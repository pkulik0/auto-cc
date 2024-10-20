package youtube

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/pkulik0/autocc/api/internal/cache"
	"github.com/pkulik0/autocc/api/internal/pb"
	"github.com/pkulik0/autocc/api/internal/srt"
	"github.com/pkulik0/autocc/api/internal/store"
)

// Youtube is an interface for the YouTube service.
//
//go:generate mockgen -destination=../mock/youtube.go -package=mock . Youtube
type Youtube interface {
	// GetVideos returns a list of videos uploaded by the authenticated user.
	GetVideos(ctx context.Context, userID, nextPageToken string) ([]*pb.Video, string, error)

	// GetMetadata returns metadata for a video.
	GetMetadata(ctx context.Context, userID, videoID string) (*Metadata, error)
	// UpdateMetadata updates metadata for a video for each language.
	UpdateMetadata(ctx context.Context, userID, videoID string, metadata map[string]*Metadata) error

	// GetCC returns a list of closed captions for a video.
	GetCC(ctx context.Context, userID, videoID string) ([]*CC, error)
	// DownloadCC downloads closed captions for a video.
	DownloadCC(ctx context.Context, userID, ccID string) (*srt.Srt, error)
	// UploadCC uploads closed captions for a video.
	UploadCC(ctx context.Context, userID, videoID, language string, srt *srt.Srt) (string, error)
}

var _ Youtube = &youtube{}

type youtube struct {
	store store.Store
	cache cache.Cache
}

// New creates a new YouTube service.
func New(store store.Store, cache cache.Cache) *youtube {
	log.Debug().Msg("created youtube service")
	return &youtube{
		store: store,
		cache: cache,
	}
}
