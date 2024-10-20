package youtube

import (
	"context"

	yt "google.golang.org/api/youtube/v3"

	"github.com/pkulik0/autocc/api/internal/errs"
	"github.com/pkulik0/autocc/api/internal/pb"
	"github.com/pkulik0/autocc/api/internal/quota"
)

func (y *youtube) GetMetadata(ctx context.Context, userID, videoID string) (*pb.Metadata, error) {
	if userID == "" || videoID == "" {
		return nil, errs.InvalidInput
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
		return nil, errs.NotFound
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
		return errs.InvalidInput
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
