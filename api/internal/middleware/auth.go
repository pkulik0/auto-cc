package middleware

import (
	"net/http"
	"strings"

	"github.com/pkulik0/autocc/api/internal/auth"
)

// Auth is a middleware that authenticates users.
func Auth(a auth.Auth, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" || !strings.HasPrefix(accessToken, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken = strings.TrimPrefix(accessToken, "Bearer ")
		userID, isSuperuser, err := a.Authenticate(r.Context(), accessToken)
		switch err {
		case nil:
		case auth.ErrInvalidToken:
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r.WithContext(auth.ContextWithUser(r.Context(), userID, isSuperuser)))
	})
}
