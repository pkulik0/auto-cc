package youtube

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	yt "google.golang.org/api/youtube/v3"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/pkulik0/autocc/api/internal/pb"
	"github.com/pkulik0/autocc/api/internal/quota"
	"github.com/pkulik0/autocc/api/internal/store"
)

// Youtube is an interface for the YouTube service.
//
//go:generate mockgen -destination=../mocks/youtube.go -package=mocks . Youtube
type Youtube interface {
	// GetVideos returns a list of videos uploaded by the authenticated user.
	GetVideos(ctx context.Context, userID, nextPageToken string) ([]*pb.Video, string, error)

	// GetMetadata returns metadata for a video.
	GetMetadata(ctx context.Context, userID, videoID string) (*pb.Metadata, error)
	// UpdateMetadata updates metadata for a video for each language.
	UpdateMetadata(ctx context.Context, userID, videoID string, metadata map[string]*pb.Metadata) error

	// GetClosedCaptions returns a list of closed captions for a video.
	GetClosedCaptions(ctx context.Context, userID, videoID string) ([]*pb.ClosedCaptionsEntry, error)
	// DownloadClosedCaptions downloads closed captions for a video.
	DownloadClosedCaptions(ctx context.Context, userID, ccID string) (string, error)
}

var _ Youtube = &youtube{}

type youtube struct {
	store store.Store
}

// New creates a new YouTube service.
func New(store store.Store) *youtube {
	return &youtube{
		store: store,
	}
}

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

const (
	videosMaxResults = 50
	captionsFormat   = "srt"
)

func (y *youtube) GetVideos(ctx context.Context, userID, nextPageToken string) ([]*pb.Video, string, error) {
	if userID == "" {
		return nil, "", ErrInvalidInput
	}

	service, err := y.getInstance(ctx, userID, quota.YoutubeSearchList)
	if err != nil {
		return nil, "", err
	}

	call := service.Search.List([]string{"snippet"}).ForMine(true).MaxResults(videosMaxResults).Type("video")
	if nextPageToken != "" {
		call.PageToken(nextPageToken)
	}

	resp, err := call.Do()
	if err != nil {
		return nil, "", err
	}

	var videos []*pb.Video
	for _, item := range resp.Items {
		if item.Id.Kind != "youtube#video" {
			continue
		}

		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			continue
		}

		videos = append(videos, &pb.Video{
			Id:           item.Id.VideoId,
			Title:        item.Snippet.Title,
			ThumbnailUrl: item.Snippet.Thumbnails.Default.Url,
			Description:  item.Snippet.Description,
			PublishedAt:  &timestamppb.Timestamp{Seconds: publishedAt.Unix()},
		})
	}

	return videos, resp.NextPageToken, nil
}

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

func (y *youtube) GetMetadata(ctx context.Context, userID, videoID string) (*pb.Metadata, error) {
	if userID == "" || videoID == "" {
		return nil, ErrInvalidInput
	}

	service, err := y.getInstance(ctx, userID, quota.YoutubeVideosList)
	if err != nil {
		return nil, err
	}

	resp, err := service.Videos.List([]string{"snippet"}).Id(videoID).Do()
	if err != nil {
		return nil, err
	}
	if len(resp.Items) == 0 {
		return nil, ErrNotFound
	}

	metadata := resp.Items[0].Snippet
	return &pb.Metadata{
		Title:       metadata.Title,
		Description: metadata.Description,
		Language:    metadata.DefaultLanguage,
	}, nil
}

func (y *youtube) UpdateMetadata(ctx context.Context, userID, videoID string, metadata map[string]*pb.Metadata) error {
	if userID == "" || videoID == "" || len(metadata) == 0 {
		return ErrInvalidInput
	}

	service, err := y.getInstance(ctx, userID, quota.YoutubeVideosUpdate)
	if err != nil {
		return err
	}

	localizations := make(map[string]yt.VideoLocalization)
	for lang, meta := range metadata {
		localizations[lang] = yt.VideoLocalization{
			Title:       meta.Title,
			Description: meta.Description,
		}
	}

	_, err = service.Videos.Update([]string{"localizations"}, &yt.Video{
		Id:            videoID,
		Localizations: localizations,
	}).Do()
	// TODO: Handle not found
	return err
}

func (y *youtube) GetClosedCaptions(ctx context.Context, userID, videoID string) ([]*pb.ClosedCaptionsEntry, error) {
	if userID == "" || videoID == "" {
		return nil, ErrInvalidInput
	}

	service, err := y.getInstance(ctx, userID, quota.YoutubeCaptionsList)
	if err != nil {
		return nil, err
	}

	resp, err := service.Captions.List([]string{"id", "snippet"}, videoID).Do()
	if err != nil {
		return nil, err
	}

	var captions []*pb.ClosedCaptionsEntry
	for _, item := range resp.Items {
		captions = append(captions, &pb.ClosedCaptionsEntry{
			Id:       item.Id,
			Language: item.Snippet.Language,
		})
	}

	return captions, nil
}

func (y *youtube) DownloadClosedCaptions(ctx context.Context, userID, ccID string) (string, error) {
	if userID == "" || ccID == "" {
		return "", ErrInvalidInput
	}

	service, err := y.getInstance(ctx, userID, quota.YoutubeCaptionsDownload)
	if err != nil {
		return "", err
	}

	resp, err := service.Captions.Download(ccID).Tfmt(captionsFormat).Download()
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
