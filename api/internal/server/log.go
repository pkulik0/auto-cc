package server

import (
	"net/http"
	"time"

	"github.com/pkulik0/autocc/api/internal/auth"
	"github.com/rs/zerolog/log"
)

type httpWriter struct {
	w          http.ResponseWriter
	statusCode int
}

func (h *httpWriter) Header() http.Header {
	return h.w.Header()
}

func (h *httpWriter) Write(data []byte) (int, error) {
	return h.w.Write(data)
}

func (h *httpWriter) WriteHeader(statusCode int) {
	h.statusCode = statusCode
	h.w.WriteHeader(statusCode)
}

func (h *httpWriter) StatusCode() int {
	if h.statusCode == 0 {
		return http.StatusOK
	}
	return h.statusCode
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writer := &httpWriter{w: w}
		next.ServeHTTP(writer, r)

		userID, isSuperuser, _ := auth.UserFromContext(r.Context())
		log.Info().Str("method", r.Method).Str("path", r.URL.Path).Int("status",
			writer.StatusCode()).Dur("duration", time.Since(start)).Str("remote",
			r.RemoteAddr).Str("user_id", userID).Bool("is_superuser", isSuperuser).Msg("request")
	})
}
