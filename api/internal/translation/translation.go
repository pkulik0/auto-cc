package translation

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/pkulik0/autocc/api/internal/store"
)

// Translator is the interface that wraps translation methods.
//
//go:generate mockgen -destination=../mock/translation.go -package=mock . Translator
type Translator interface {
	// GetTargetLanguages returns a list of supported languages.
	GetLanguages(ctx context.Context) ([]string, error)
	// Translate translates the text from the source language to the target language.
	Translate(ctx context.Context, text []string, sourceLanguage, targetLanguage string) ([]string, error)
	// GetUsageDeepL returns the usage of the DeepL API.
	GetUsageDeepL(ctx context.Context, apiKey string) (uint, error)
}

type translator struct {
	store store.Store
}

var _ Translator = &translator{}

// New creates a new translation service.
func New(store store.Store) *translator {
	log.Debug().Msg("created translation service")
	return &translator{
		store: store,
	}
}

func (t *translator) GetLanguages(ctx context.Context) ([]string, error) {
	apiClient, err := newDeeplApiClient(ctx, t.store, 0)
	if err != nil {
		return nil, err
	}

	languages, err := apiClient.getLanguages(ctx)
	if err != nil {
		return nil, err
	}

	log.Debug().Int("count", len(languages)).Strs("languages", languages).Msg("fetched languages")
	return languages, nil
}

func countTextLen(text []string) uint {
	var count uint
	for _, t := range text {
		count += uint(len(t))
	}
	return count
}

var (
	ErrInvalidInput = errors.New("invalid input")
)

func (t *translator) Translate(ctx context.Context, text []string, sourceLanguage, targetLanguage string) ([]string, error) {
	if len(text) == 0 || sourceLanguage == "" || targetLanguage == "" {
		return nil, ErrInvalidInput
	}

	cost := countTextLen(text)

	apiClient, err := newDeeplApiClient(ctx, t.store, cost)
	if err != nil {
		return nil, err
	}

	translatedText, err := apiClient.translate(ctx, text, sourceLanguage, targetLanguage)
	if err != nil {
		return nil, err
	}

	log.Debug().Strs("text", text).Strs("translated_text", translatedText).Msg("translated text")
	return translatedText, nil
}

func (t *translator) GetUsageDeepL(ctx context.Context, apiKey string) (uint, error) {
	client := &http.Client{
		Transport: newDeeplTransport(http.DefaultTransport, apiKey),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"usage", nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var data struct {
		CharacterCount uint `json:"character_count"`
		CharacterLimit uint `json:"character_limit"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	return data.CharacterCount, nil
}
