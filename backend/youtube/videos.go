package main

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
	youtubeClient := ctx.Locals("youtubeClient").(*youtube.Service)

	if err := s.checkQuotaAndRotateIdentity(quotaCostVideos); err != nil {
		log.Errorf("Quota check failed: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to check quota.")
	}

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
