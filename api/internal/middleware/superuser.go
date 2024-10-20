package middleware

import (
	"net/http"

	"github.com/pkulik0/autocc/api/internal/auth"
	"github.com/pkulik0/autocc/api/internal/helpers"
)

// Superuser is a middleware that checks if the user is a superuser.
func Superuser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, isSuperuser, ok := auth.UserFromContext(r.Context())
		if !ok {
			helpers.ErrLog(w, nil, "failed to get user from context", http.StatusInternalServerError)
			return
		}
		if !isSuperuser {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
