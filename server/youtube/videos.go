package youtube

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/api/youtube/v3"
	"time"
)

type Video struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	Description  string `json:"description"`
	PublishedAt  int64  `json:"publishedAt"`
}

func (s *Service) videosHandler(ctx *fiber.Ctx) error {
	if err := s.checkQuotaAndRotateIdentity(ctx, quotaCostVideos); err != nil {
		log.Errorf("Quota check failed: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to check quota.")
	}

	youtubeClient := ctx.Locals("youtubeClient").(*youtube.Service)
	req := youtubeClient.Search.List([]string{"snippet"}).ForMine(true).MaxResults(50).Type("video")
	if pageToken := ctx.Query("token"); pageToken != "" {
		req = req.PageToken(pageToken)
	}

	res, err := req.Do()
	if err != nil {
		log.Errorf("YT video search failure: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to get yt response.")
	}

	_, err = s.addToQuota(quotaCostVideos)
	if err != nil {
		log.Errorf("Failed to increment quota: %s", err)
	}

	var videos []Video
	for _, raw := range res.Items {
		publicationTime, err := time.Parse(time.RFC3339, raw.Snippet.PublishedAt)
		if err != nil {
			continue
		}

		video := Video{
			Id:           raw.Id.VideoId,
			Title:        raw.Snippet.Title,
			ThumbnailUrl: raw.Snippet.Thumbnails.Maxres.Url,
			Description:  raw.Snippet.Description,
			PublishedAt:  publicationTime.UnixMilli(),
		}
		videos = append(videos, video)
	}

	return ctx.JSON(fiber.Map{
		"videos":        videos,
		"nextPageToken": res.NextPageToken,
	})
}

type VideoMetadata struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Language    string `json:"language"`
}

func (s *Service) videoMetadataHandler(ctx *fiber.Ctx) error {
	videoId := ctx.Params("videoId")
	if videoId == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing video id")
	}

	if err := s.checkQuotaAndRotateIdentity(ctx, quotaCostVideoInfo); err != nil {
		log.Errorf("Quota check failed: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to check quota.")
	}

	youtubeClient := ctx.Locals("youtubeClient").(*youtube.Service)
	res, err := youtubeClient.Videos.List([]string{"snippet"}).Id(videoId).Do()
	if err != nil {
		log.Errorf("Failed to get video metadata: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to fetch yt response")
	}

	_, err = s.addToQuota(quotaCostVideoInfo)
	if err != nil {
		log.Errorf("Failed to increment quota: %s", err)
	}

	if len(res.Items) < 1 {
		return ctx.Status(fiber.StatusNotFound).SendString("Video not found")
	}
	videoSnippet := res.Items[0].Snippet

	return ctx.JSON(VideoMetadata{
		Title:       videoSnippet.Title,
		Description: videoSnippet.Description,
		Language:    videoSnippet.DefaultLanguage,
	})
}

func convertLanguageCode(translatorCode string) string {
	switch translatorCode {
	case "NB":
		return "NO"
	default:
		return translatorCode
	}
}

func (s *Service) videoUpdateMetadataHandler(ctx *fiber.Ctx) error {
	videoId := ctx.Params("videoId")
	if videoId == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing video id")
	}

	var metadataArray []VideoMetadata
	if err := ctx.BodyParser(&metadataArray); err != nil {
		log.Errorf("Failed to parse video update body: %s", err)
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	localizations := make(map[string]youtube.VideoLocalization)
	for _, metadata := range metadataArray {
		localizations[convertLanguageCode(metadata.Language)] = youtube.VideoLocalization{
			Title:       metadata.Title,
			Description: metadata.Description,
		}
	}

	if err := s.checkQuotaAndRotateIdentity(ctx, quotaCostUpdateVideo); err != nil {
		log.Errorf("Quota check failed: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to check quota.")
	}

	video := &youtube.Video{
		Id:            videoId,
		Localizations: localizations,
	}

	youtubeClient := ctx.Locals("youtubeClient").(*youtube.Service)
	res, err := youtubeClient.Videos.Update([]string{"localizations"}, video).Do()
	if err != nil {
		log.Errorf("Failed to update video metadata: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to fetch yt response")
	}

	_, err = s.addToQuota(quotaCostUpdateVideo)
	if err != nil {
		log.Errorf("Failed to increment quota: %s", err)
	}

	return ctx.JSON(res)
}
