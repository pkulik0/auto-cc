package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// Panic is a middleware that recovers from panics.
func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rec := recover(); rec != nil {
			log.Error().Interface("recovered", rec).Msg("panic recovered")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
