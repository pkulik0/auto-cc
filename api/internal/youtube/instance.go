package youtube

import (
	"context"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	yt "google.golang.org/api/youtube/v3"
)

// getInstance returns an authenticated Youtube service with enough quota to make the request.
func (y *youtube) getInstance(ctx context.Context, userID string, neededQuota uint) (service *yt.Service, err error) {
	session, revert, err := y.store.GetSessionGoogleByAvailableCost(ctx, userID, neededQuota)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err == nil {
			return
		}
		revertErr := revert()
		if revertErr != nil {
			log.Err(revertErr).Str("user_id", userID).Uint("session_id", session.ID).Uint("cost", neededQuota).Msg("failed to revert session cost")
		} else {
			log.Debug().Str("user_id", userID).Uint("session_id", session.ID).Uint("cost", neededQuota).Msg("reverted session cost")
		}
	}()

	src, err := session.GetTokenSource(ctx, func(t *oauth2.Token) {
		session.AccessToken = t.AccessToken
		session.RefreshToken = t.RefreshToken
		session.Expiry = t.Expiry
		err := y.store.UpdateSessionGoogle(ctx, session)
		if err != nil {
			log.Err(err).Str("user_id", userID).Uint("session_id", session.ID).Msg("failed to update session token")
		} else {
			log.Debug().Str("user_id", userID).Uint("session_id", session.ID).Msg("updated session token")
		}
	})
	if err != nil {
		return nil, err
	}
	client := oauth2.NewClient(ctx, src)

	service, err = yt.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return service, nil
}
