package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/pkulik0/autocc/api/internal/auth"
	"github.com/pkulik0/autocc/api/internal/cache"
	"github.com/pkulik0/autocc/api/internal/helpers"
)

var (
	cacheDuration = 5 * time.Minute
)

func getCacheKey(r *http.Request) (string, error) {
	hash := sha256.New()

	_, err := hash.Write([]byte(r.Method))
	if err != nil {
		return "", err
	}

	_, err = hash.Write([]byte(r.URL.String()))
	if err != nil {
		return "", err
	}

	user, _, ok := auth.UserFromContext(r.Context())
	if ok {
		_, err = hash.Write([]byte(user))
		if err != nil {
			return "", err
		}
	}

	if r.Method != http.MethodGet {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		_, err = hash.Write(body)
		if err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func nextWithCache(next http.Handler, w *httpWriter, r *http.Request, c cache.Cache, key string) {
	w.Header().Set("Cache-Control", "private, max-age="+strconv.Itoa(int(cacheDuration.Seconds())))
	next.ServeHTTP(w, r)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := c.Set(ctx, key, w.data.String(), cacheDuration); err != nil {
			log.Error().Err(err).Str("key", key).Msg("failed to set cache")
		} else {
			log.Trace().Str("key", key).Dur("duration", cacheDuration).Msg("set cache")
		}
	}()
}

// Cache is a middleware that caches responses.
func Cache(c cache.Cache, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		w, ok := writer.(*httpWriter)
		if !ok {
			log.Warn().Msg("cache middleware: writer is not httpWriter")
			w = newHttpWriter(w)
		}

		key, err := getCacheKey(r)
		if err != nil {
			helpers.ErrLog(w, err, "failed to get cache key", http.StatusInternalServerError)
			return
		}

		if r.Header.Get("Cache-Control") == "no-cache" {
			nextWithCache(next, w, r, c, key)
			return
		}

		value, err := c.Get(r.Context(), key)
		if err != nil {
			log.Trace().Str("cache-key", key).Msg("cache miss")
			w.Header().Set("X-Cache", "MISS")
			w.Header().Set("X-Cache-Key", key)

			nextWithCache(next, w, r, c, key)
		} else {
			log.Trace().Str("cache-key", key).Msg("cache hit")
			w.Header().Set("X-Cache", "HIT")
			w.Header().Set("X-Cache-Key", key)

			helpers.WriteOrLog(w, []byte(value))
		}
	})
}
