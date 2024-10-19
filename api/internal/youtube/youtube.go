package youtube

import (
	"context"
	"errors"
	"io"

	"github.com/rs/zerolog/log"

	"github.com/pkulik0/autocc/api/internal/pb"
	"github.com/pkulik0/autocc/api/internal/store"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
)

// Youtube is an interface for the YouTube service.
//
//go:generate mockgen -destination=../mock/youtube.go -package=mock . Youtube
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
	// UploadClosedCaptions uploads closed captions for a video.
	UploadClosedCaptions(ctx context.Context, userID, videoID, language string, srt io.Reader) (string, error)
}

var _ Youtube = &youtube{}

type youtube struct {
	store store.Store
}

// New creates a new YouTube service.
func New(store store.Store) *youtube {
	log.Debug().Msg("created youtube service")
	return &youtube{
		store: store,
	}
}
