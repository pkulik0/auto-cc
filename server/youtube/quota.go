package youtube

import (
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"strconv"
	"time"
)

const (
	quotaCacheKey             = "quota_"
	quotaLimit                = 10000
	quotaCostVideoInfo        = 1
	quotaCostUpdateVideo      = 50
	quotaCostVideos           = 100
	quotaCostListCaptions     = 50
	quotaCostDownloadCaptions = 200
	quotaCostInsertCaptions   = 400
)

func getQuotaResetTime() time.Time {
	nowLocal := time.Now()
	nowPST := nowLocal.In(time.FixedZone("PST", -9*60*60))

	passedMidnight := time.Date(nowPST.Year(), nowPST.Month(), nowPST.Day(), 0, 0, 0, 0, nowPST.Location())
	nextMidnight := passedMidnight.AddDate(0, 0, 1)

	return nextMidnight
}

func (s *Service) addToQuota(usedPoints int64) (int64, error) {
	identity, err := s.Identity()
	if err != nil {
		return 0, nil
	}
	identityHash := identity.Hash()

	currentlyUsed, err := s.rdb.IncrBy(quotaCacheKey+identityHash, usedPoints).Result()
	if err != nil {
		return 0, err
	}

	if err := s.rdb.ExpireAt(quotaCacheKey+identityHash, getQuotaResetTime()).Err(); err != nil {
		log.Errorf("Failed to set expiry on key \"%s\": %s", quotaCacheKey, err)
	}

	return currentlyUsed, nil
}

func (s *Service) getQuota(identityHash string) (uint64, error) {
	cachedValue, err := s.rdb.Get(quotaCacheKey + identityHash).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	parsedValue, err := strconv.ParseUint(cachedValue, 10, 0)
	if err != nil {
		return 0, err
	}

	return parsedValue, nil
}

func (s *Service) checkQuotaAndRotateIdentity(ctx *fiber.Ctx, neededPoints uint64) error {
	identity, err := s.Identity()
	if err != nil {
		return err
	}

	var seenHashes []string
	for {
		identityHash := identity.Hash()

		for _, seenHash := range seenHashes {
			if seenHash != identityHash {
				continue
			}
			panic("No more identities under quota limit")
		}
		seenHashes = append(seenHashes, identityHash)

		quotaUsage, err := s.getQuota(identityHash)
		if err != nil {
			return err
		}
		if quotaUsage+neededPoints <= quotaLimit {
			return s.generateClientFromCurrentIdentity(ctx)
		}

		identity, err = s.nextIdentity()
		if err != nil {
			return err
		}
	}
}
