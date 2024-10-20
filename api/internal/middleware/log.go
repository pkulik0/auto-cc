package middleware

import (
	"net/http"
	"time"

	"github.com/pkulik0/autocc/api/internal/auth"
	"github.com/rs/zerolog/log"
)

// Log is a middleware that logs requests.
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writer := newHttpWriter(w)
		next.ServeHTTP(writer, r)

		userID, isSuperuser, _ := auth.UserFromContext(r.Context())
		log.Info().Str("method", r.Method).Str("path", r.URL.Path).Int("status",
			writer.StatusCode()).Dur("duration", time.Since(start)).Str("remote",
			r.RemoteAddr).Str("user_id", userID).Bool("is_superuser", isSuperuser).Msg("request")
	})
}
