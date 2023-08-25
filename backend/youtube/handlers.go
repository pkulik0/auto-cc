package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
	"io"
	"time"
)

func (s *Service) callbackHandler(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	state := ctx.Query("state")
	if code == "" || state == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing oauth2 code/state.")
	}

	token, err := exchangeToken(ctx.Context(), s.oauth2Config, code, state)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Failed to exchange token.")
	}

	if err := saveToken(s.rdb, token); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to save token.")
	}

	return ctx.SendString("OK")
}

func (s *Service) authHandler(ctx *fiber.Ctx) error {
	url := s.oauth2Config.AuthCodeURL("stateTODO", oauth2.AccessTypeOffline)
	return ctx.Redirect(url, fiber.StatusFound)
}

type Video struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	Description  string `json:"description"`
	PublishedAt  int64  `json:"publishedAt"`
}

func (s *Service) videosHandler(ctx *fiber.Ctx) error {
	youtubeClient, ok := ctx.Locals("youtubeClient").(*youtube.Service)
	if !ok {
		log.Error("Failed to get oauth2Client from locals.")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal error.")
	}

	res, err := youtubeClient.Search.List([]string{"snippet"}).ForMine(true).MaxResults(50).Type("video").Do()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to get yt response.")
	}
	_ = s.addToQuota(100)

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
			PublishedAt:  publicationTime.Unix(),
		}
		videos = append(videos, video)
	}

	return ctx.JSON(videos)
}

type ClosedCaptions struct {
	Id       string `json:"id"`
	Language string `json:"language"`
}

func (s *Service) ccListHandler(ctx *fiber.Ctx) error {
	youtubeClient, ok := ctx.Locals("youtubeClient").(*youtube.Service)
	if !ok {
		log.Error("Failed to get oauth2Client from locals.")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal error.")
	}

	videoId := ctx.Params("videoId")
	if videoId == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing video id.")
	}

	res, err := youtubeClient.Captions.List([]string{"id", "snippet"}, videoId).Do()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to get yt response.")
	}
	_ = s.addToQuota(50)

	var ccs []ClosedCaptions
	for _, raw := range res.Items {
		cc := ClosedCaptions{
			Id:       raw.Id,
			Language: raw.Snippet.Language,
		}
		ccs = append(ccs, cc)
	}

	return ctx.JSON(ccs)
}

func (s *Service) ccDownloadHandler(ctx *fiber.Ctx) error {
	ccId := ctx.Params("ccId")
	if ccId == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing ccId.")
	}

	youtubeClient, ok := ctx.Locals("youtubeClient").(*youtube.Service)
	if !ok {
		log.Error("Failed to get oauth2Client from locals.")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal error.")
	}

	res, err := youtubeClient.Captions.Download(ccId).Tfmt("srt").Download()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to get yt response.")
	}
	_ = s.addToQuota(200)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to read parse response.")
	}
	_ = res.Body.Close()

	return ctx.Send(body)
}
