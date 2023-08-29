package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const identitiesKey = "identities"

type Identity struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func (i *Identity) Hash() string {
	hash := sha256.Sum256([]byte(i.ClientSecret))
	return hex.EncodeToString(hash[:])
}

func (i *Identity) getOAuth2Config() *oauth2.Config {
	redirectUri := getEnv("GOOGLE_REDIRECT_URI")

	return &oauth2.Config{
		ClientID:     i.ClientId,
		ClientSecret: i.ClientSecret,
		RedirectURL:  redirectUri,
		Scopes:       []string{youtube.YoutubeForceSslScope},
		Endpoint:     google.Endpoint,
	}
}

type IdentityInfo struct {
	IdentityHash string `json:"identityHash"`
	UsedQuota    uint64 `json:"usedQuota"`
	IsSelected   bool   `json:"isSelected"`
}

func (s *Service) getIdentitiesHandler(ctx *fiber.Ctx) error {
	currentIdentity, err := s.Identity()
	if err != nil {
		log.Errorf("Failed to get identity: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to get current identity.")
	}
	currentIdentityHash := currentIdentity.Hash()

	identityHashes := []string{}
	err = s.forEachIdentity(func(identity *Identity) error {
		identityHashes = append(identityHashes, identity.Hash())
		return nil
	})
	if err != nil {
		log.Errorf("Failed to iterate over identities: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve identities")
	}

	identityInfos := []IdentityInfo{}
	for _, identityHash := range identityHashes {
		usedQuota, err := s.getQuota(identityHash)
		if err != nil {
			log.Errorf("Failed to get quota: %s", err)
			return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to check quotas.")
		}
		identityInfos = append(identityInfos, IdentityInfo{
			IdentityHash: identityHash,
			UsedQuota:    usedQuota,
			IsSelected:   identityHash == currentIdentityHash,
		})
	}

	return ctx.JSON(identityInfos)
}

func (s *Service) addIdentityHandler(ctx *fiber.Ctx) error {
	var identity Identity
	if err := ctx.BodyParser(&identity); err != nil {
		log.Errorf("Failed to parse identity from body: %s", err)
		return err
	}

	identityJson, err := json.Marshal(identity)
	if err != nil {
		log.Errorf("Failed to marshal received identity: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Received identity couldn't be processed.")
	}

	if err := s.rdb.RPush(identitiesKey, identityJson).Err(); err != nil {
		log.Error("Failed to add identity to db: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to save identity.")
	}

	return ctx.SendString("OK")
}

func (s *Service) nextIdentity() (*Identity, error) {
	identityEntry, err := s.rdb.LPop(identitiesKey).Result()
	if err != nil {
		return nil, err
	}

	if err := s.rdb.RPush(identitiesKey, identityEntry).Err(); err != nil {
		return nil, err
	}

	var identity Identity
	if err := json.Unmarshal([]byte(identityEntry), &identity); err != nil {
		return nil, err
	}

	s.identity = &identity
	return &identity, nil
}

func (s *Service) forEachIdentity(identityProcessor func(identity *Identity) error) error {
	identitiesJson, err := s.rdb.LRange(identitiesKey, 0, -1).Result()
	if err != nil {
		return err
	}

	for _, identityJson := range identitiesJson {
		var identity Identity
		if err := json.Unmarshal([]byte(identityJson), &identity); err != nil {
			return err
		}

		if err := identityProcessor(&identity); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) callbackHandler(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	state := ctx.Query("state")
	if code == "" || state == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing oauth2 code/state.")
	}

	var identityFromState *Identity = nil
	err := s.forEachIdentity(func(identity *Identity) error {
		if identity.Hash() == state {
			identityFromState = identity
		}
		return nil
	})
	if err != nil {
		return err
	}
	if identityFromState == nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid state value.")
	}

	config := identityFromState.getOAuth2Config()
	token, err := exchangeToken(ctx.Context(), config, code, state)
	if err != nil {
		log.Errorf("Failed to exchange token: %s", err)
		return ctx.Status(fiber.StatusBadRequest).SendString("Failed to exchange token.")
	}

	if err := saveToken(s.rdb, token, identityFromState); err != nil {
		log.Errorf("Failed to save token: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to save token.")
	}

	return ctx.SendString("OK")
}

type AuthUrl struct {
	Url          string `json:"url"`
	IdentityHash string `json:"identityHash"`
}

func (s *Service) authUrlsHandler(ctx *fiber.Ctx) error {
	authUrls := []AuthUrl{}

	err := s.forEachIdentity(func(identity *Identity) error {
		_, err := s.rdb.Get(tokenCacheKey + identity.Hash()).Result()
		if err == nil {
			return nil
		}
		if err != redis.Nil {
			log.Errorf("Token retrieval error: %s", err)
			return ctx.Status(fiber.StatusInternalServerError).SendString("Error encountered while retrieving identities.")
		}

		url := identity.getOAuth2Config().AuthCodeURL(identity.Hash(), oauth2.AccessTypeOffline)
		authUrls = append(authUrls, AuthUrl{
			Url:          url,
			IdentityHash: identity.Hash(),
		})

		return nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(authUrls)
}

func (s *Service) generateClientFromCurrentIdentity(ctx *fiber.Ctx) error {
	identity, err := s.Identity()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Identity error.")
	}

	token, err := getTokenFromCache(s.rdb, identity)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString("No access.")
	}

	config := identity.getOAuth2Config()
	tokenSource := config.TokenSource(ctx.Context(), token)
	tokenWrapper := NewTokenWrapper(s.rdb, tokenSource, identity)

	oauth2Client := oauth2.NewClient(ctx.Context(), tokenWrapper)
	youtubeClient, err := youtube.NewService(ctx.Context(), option.WithHTTPClient(oauth2Client))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to create youtube client.")
	}

	ctx.Locals("youtubeClient", youtubeClient)
	return nil
}
