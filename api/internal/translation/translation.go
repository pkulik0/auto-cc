package translation

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/pkulik0/autocc/api/internal/store"
)

// Language represents a natural language.
type Language string

func (l Language) String() string {
	return string(l)
}

// Translator is the interface that wraps translation methods.
//
//go:generate mockgen -destination=../mock/translation.go -package=mock . Translator
type Translator interface {
	// GetTargetLanguages returns a list of supported languages.
	GetLanguages(ctx context.Context) ([]Language, error)
	// Translate translates the text from the source language to the target language.
	Translate(ctx context.Context, text []string, source, target Language) ([]string, error)
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

func (d *translator) GetLanguages(ctx context.Context) ([]Language, error) {
	apiClient, err := newDeeplApiClient(ctx, d.store, 0)
	if err != nil {
		return nil, err
	}

	return apiClient.getLanguages(ctx)
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

func (d *translator) Translate(ctx context.Context, text []string, source, target Language) ([]string, error) {
	if len(text) == 0 || source == "" || target == "" {
		return nil, ErrInvalidInput
	}

	cost := countTextLen(text)

	apiClient, err := newDeeplApiClient(ctx, d.store, cost)
	if err != nil {
		return nil, err
	}

	return apiClient.translate(ctx, text, source, target)
}
