package youtube

import (
	"context"
	"io"

	yt "google.golang.org/api/youtube/v3"

	"github.com/pkulik0/autocc/api/internal/pb"
	"github.com/pkulik0/autocc/api/internal/quota"
)

const (
	captionsFormat = "srt"
)

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

func (y *youtube) UploadClosedCaptions(ctx context.Context, userID, videoID, language string, srt io.Reader) (string, error) {
	if userID == "" || videoID == "" || language == "" || srt == nil {
		return "", ErrInvalidInput
	}

	service, err := y.getInstance(ctx, userID, quota.YoutubeCaptionsUpload)
	if err != nil {
		return "", err
	}

	resp, err := service.Captions.Insert([]string{"snippet"}, &yt.Caption{Snippet: &yt.CaptionSnippet{
		Language: language,
		Name:     language,
		VideoId:  videoID,
	}}).Media(srt).Do()
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}
