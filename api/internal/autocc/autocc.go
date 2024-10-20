package autocc

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkulik0/autocc/api/internal/errs"
	"github.com/pkulik0/autocc/api/internal/translation"
	"github.com/pkulik0/autocc/api/internal/youtube"
	"github.com/rs/zerolog/log"
)

// AutoCC is the interface that wraps the main workflow of AutoCC.
type AutoCC interface {
	// Process processes the video and uploads translated closed captions and metadata.
	Process(ctx context.Context, userID, videoID string) error
}

var _ AutoCC = &autoCC{}

type autoCC struct {
	translator translation.Translator
	youtube    youtube.Youtube
}

// New creates a new AutoCC service.
func New(translator translation.Translator, youtube youtube.Youtube) *autoCC {
	return &autoCC{
		translator: translator,
		youtube:    youtube,
	}
}

func (a *autoCC) Process(ctx context.Context, userID, videoID string) error {
	if userID == "" || videoID == "" {
		return errs.InvalidInput
	}

	languages, err := a.translator.GetLanguages(ctx)
	if err != nil {
		return err
	}

	metadata, err := a.youtube.GetMetadata(ctx, userID, videoID)
	if err != nil {
		return err
	}

	allCC, err := a.youtube.GetCC(ctx, userID, videoID)
	if err != nil {
		return err
	}

	var srcCC *youtube.CC
	for _, cc := range allCC {
		if cc.Language == metadata.Language {
			srcCC = cc
			break
		}
	}
	if srcCC == nil {
		return errs.SourceClosedCaptionsNotFound
	}

	srt, err := a.youtube.DownloadCC(ctx, userID, srcCC.Id)
	if err != nil {
		return err
	}

	errChan := make(chan error)
	metadataChan := make(chan *youtube.Metadata, len(languages))
	waitGroupMetadata := sync.WaitGroup{}
	waitGroupCC := sync.WaitGroup{}

	for i, targetLang := range languages {
		srcLang := translation.CodeGoogleToTranslation(metadata.Language)
		if targetLang == srcLang {
			continue
		}

		waitGroupCC.Add(1)
		go func() {
			defer func() {
				log.Debug().Str("src_lang", srcLang).Str("target_lang", targetLang).Int("i", i).Int("len", len(languages)-1).Msg("finished translating cc")
				waitGroupCC.Done()
			}()

			text, err := a.translator.Translate(ctx, srt.Text(), srcLang, targetLang)
			if err != nil {
				log.Error().Err(err).Str("src_lang", srcLang).Str("target_lang", targetLang).Msg("failed to translate text")
				errChan <- err
				return
			}

			translatedSrt := srt.Clone()
			err = translatedSrt.ReplaceText(text)
			if err != nil {
				log.Error().Err(err).Msg("failed to replace text")
				errChan <- err
				return
			}

			_, err = a.youtube.UploadCC(ctx, userID, videoID, translation.CodeTranslationToGoogle(targetLang), translatedSrt)
			if err != nil {
				log.Error().Err(err).Str("src_lang", srcLang).Str("target_lang", targetLang).Msg("failed to upload cc")
				errChan <- err
			}
		}()

		waitGroupMetadata.Add(1)
		go func() {
			defer func() {
				log.Debug().Str("src_lang", srcLang).Str("target_lang", targetLang).Int("i", i).Int("len", len(languages)-1).Msg("finished translating metadata")
				waitGroupMetadata.Done()
			}()
			text, err := a.translator.Translate(ctx, []string{metadata.Title, metadata.Description}, srcLang, targetLang)
			if err != nil {
				log.Error().Err(err).Str("src_lang", srcLang).Str("target_lang", targetLang).Msg("failed to translate metadata")
				errChan <- err
				return
			}

			if len(text) != 2 {
				log.Error().Strs("text", text).Msg("invalid translation response")
				errChan <- fmt.Errorf("invalid translation response")
				return
			}

			title := text[0]
			description := text[1]

			metadataChan <- &youtube.Metadata{
				Title:       title,
				Description: description,
				Language:    translation.CodeTranslationToGoogle(targetLang),
			}
		}()
	}

	waitGroupMetadata.Wait()
	close(metadataChan)

	if len(errChan) > 0 {
		err := <-errChan
		return err
	}

	metadataMap := make(map[string]*youtube.Metadata)
	for m := range metadataChan {
		metadataMap[m.Language] = m
	}

	err = a.youtube.UpdateMetadata(ctx, userID, videoID, metadataMap)
	if err != nil {
		return err
	}

	waitGroupCC.Wait()
	if len(errChan) > 0 {
		err := <-errChan
		return err
	}
	close(errChan)

	return nil
}
