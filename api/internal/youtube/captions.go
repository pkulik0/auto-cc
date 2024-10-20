package youtube

import (
	"context"
	"io"
	"strings"
	"time"

	yt "google.golang.org/api/youtube/v3"

	"github.com/pkulik0/autocc/api/internal/cache"
	"github.com/pkulik0/autocc/api/internal/errs"
	"github.com/pkulik0/autocc/api/internal/pb"
	"github.com/pkulik0/autocc/api/internal/quota"
	"github.com/pkulik0/autocc/api/internal/srt"
	"github.com/rs/zerolog/log"
)

const (
	captionsFormat = "srt"
)

func (y *youtube) GetClosedCaptions(ctx context.Context, userID, videoID string) ([]*pb.ClosedCaptionsEntry, error) {
	if userID == "" || videoID == "" {
		return nil, errs.InvalidInput
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

func (y *youtube) DownloadClosedCaptions(ctx context.Context, userID, ccID string) (*srt.Srt, error) {
	if userID == "" || ccID == "" {
		return nil, errs.InvalidInput
	}

	service, err := y.getInstance(ctx, userID, quota.YoutubeCaptionsDownload)
	if err != nil {
		return nil, err
	}

	resp, err := service.Captions.Download(ccID).Tfmt(captionsFormat).Download()
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	srt, err := srt.Parse(string(body))
	if err != nil {
		return nil, err
	}
	return srt, nil
}

func (y *youtube) UploadClosedCaptions(ctx context.Context, userID, videoID, language string, srt *srt.Srt) (string, error) {
	if userID == "" || videoID == "" || language == "" || srt == nil {
		return "", errs.InvalidInput
	}

	key := cache.CreateKey(userID, videoID, language, srt.String())
	if value, err := y.cache.Get(ctx, key); err == nil {
		return value, nil
	}

	service, err := y.getInstance(ctx, userID, quota.YoutubeCaptionsUpload)
	if err != nil {
		return "", err
	}

	log.Trace().Str("video_id", videoID).Str("language", language).Msg("uploading closed captions")
	resp, err := service.Captions.Insert([]string{"snippet"}, &yt.Caption{Snippet: &yt.CaptionSnippet{
		Language:        language,
		VideoId:         videoID,
		Name:            "",
		ForceSendFields: []string{"Name"}, // If not set `omitempty` kicks in and the required field is not sent.
	}}).Media(strings.NewReader(srt.String())).Do()
	if err != nil {
		return "", err
	}

	go func() {
		log.Trace().Str("key", key).Str("id", resp.Id).Msg("setting cache")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := y.cache.Set(ctx, key, resp.Id, time.Hour*24); err != nil {
			log.Error().Err(err).Str("key", key).Msg("failed to set cache")
		}
	}()
	return resp.Id, nil
}
