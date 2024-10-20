package translation

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/pkulik0/autocc/api/internal/cache"
	"github.com/pkulik0/autocc/api/internal/errs"
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
	cache cache.Cache
}

var _ Translator = &translator{}

// New creates a new translation service.
func New(store store.Store, cache cache.Cache) *translator {
	log.Debug().Msg("created translation service")
	return &translator{
		store: store,
		cache: cache,
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
	for i, l := range languages {
		languages[i] = strings.ToLower(l)
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

func (t *translator) Translate(ctx context.Context, text []string, sourceLanguage, targetLanguage string) ([]string, error) {
	if len(text) == 0 || sourceLanguage == "" || targetLanguage == "" {
		return nil, errs.InvalidInput
	}

	// Check if the translation is already in the cache.
	key := cache.CreateKey(append(text, sourceLanguage, targetLanguage)...)
	log.Trace().Str("key", key).Strs("text", text).Str("source_language", sourceLanguage).Str("target_language", targetLanguage).Msg("checking cache")
	if value, err := t.cache.GetList(ctx, key); err == nil {
		log.Trace().Strs("text", text).Strs("translated_text", value).Msg("cache hit")
		return value, nil
	}

	apiClient, err := newDeeplApiClient(ctx, t.store, countTextLen(text))
	if err != nil {
		log.Error().Err(err).Msg("failed to create DeepL API client")
		return nil, err
	}

	log.Trace().Str("source_language", sourceLanguage).Str("target_language", targetLanguage).Strs("text", text).Msg("translating text")
	translatedText, err := apiClient.translate(ctx, text, sourceLanguage, targetLanguage)
	if err != nil {
		return nil, err
	}

	go func() {
		log.Trace().Str("key", key).Str("source_language", sourceLanguage).Str("target_language", targetLanguage).Strs("translated_text", translatedText).Msg("setting cache")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := t.cache.SetList(ctx, key, translatedText, time.Hour*24)
		if err != nil {
			log.Error().Err(err).Str("key", key).Msg("failed to set cache")
		}
	}()
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
