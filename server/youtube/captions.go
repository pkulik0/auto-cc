package youtube

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/api/youtube/v3"
	"io"
)

type ClosedCaptions struct {
	Id       string `json:"id"`
	Language string `json:"language"`
}

func (s *Service) ccListHandler(ctx *fiber.Ctx) error {
	videoId := ctx.Params("videoId")
	if videoId == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing video id.")
	}

	if err := s.checkQuotaAndRotateIdentity(ctx, quotaCostListCaptions); err != nil {
		log.Errorf("Quota check failed: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to check quota.")
	}

	youtubeClient := ctx.Locals("youtubeClient").(*youtube.Service)
	res, err := youtubeClient.Captions.List([]string{"id", "snippet"}, videoId).Do()
	if err != nil {
		log.Errorf("YT CC list request failure: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to get yt response.")
	}

	_, err = s.addToQuota(quotaCostListCaptions)
	if err != nil {
		log.Errorf("Failed to increment quota: %s", err)
	}

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

	if err := s.checkQuotaAndRotateIdentity(ctx, quotaCostInsertCaptions); err != nil {
		log.Errorf("Quota check failed: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to check quota.")
	}

	youtubeClient := ctx.Locals("youtubeClient").(*youtube.Service)
	res, err := youtubeClient.Captions.Download(ccId).Tfmt("srt").Download()
	if err != nil {
		log.Errorf("YT CC download failure: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to get yt response.")
	}

	_, err = s.addToQuota(quotaCostDownloadCaptions)
	if err != nil {
		log.Errorf("Failed to increment quota: %s", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("YT CC download body parsing failure: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to read parse response.")
	}
	_ = res.Body.Close()

	return ctx.Send(body)
}

func (s *Service) ccUpload(ctx *fiber.Ctx) error {
	videoId := ctx.Params("videoId")
	language := ctx.Query("language")
	if language == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing language parameter")
	}

	if err := s.checkQuotaAndRotateIdentity(ctx, quotaCostInsertCaptions); err != nil {
		log.Errorf("Quota check failed: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to check quota.")
	}

	ccInfo := &youtube.Caption{
		Snippet: &youtube.CaptionSnippet{
			Language: language,
			VideoId:  videoId,
			Name:     language,
		},
	}
	body := ctx.Body()
	if len(body) < 1 {
		return ctx.Status(fiber.StatusBadRequest).SendString("Request body does not contain SRT subtitles")
	}
	srt := bytes.NewReader(body)

	youtubeClient := ctx.Locals("youtubeClient").(*youtube.Service)
	_, err := youtubeClient.Captions.Insert([]string{"snippet"}, ccInfo).Media(srt).Do()
	if err != nil {
		log.Errorf("YT CC insert failure: %s", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to get yt response.")
	}

	_, err = s.addToQuota(quotaCostInsertCaptions)
	if err != nil {
		log.Errorf("Failed to increment quota: %s", err)
	}

	return ctx.SendString("OK")
}
