package middleware

import (
	"net/http"
	"strings"

	"github.com/faeb5/winecellar/internal/auth"
)

func Authorized(jwtSecret string) Middleware {
	return func(next http.Handler) http.Handler {
		secret := jwtSecret
		const bearer = "Bearer "
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authStr := r.Header.Get("Authorization")
			if !strings.HasPrefix(authStr, bearer) {
				http.Error(w, "No bearer found", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authStr, bearer)

			userID, err := auth.ValidateJWT(tokenStr, secret)
			if err != nil {
				http.Error(w, "Invalid access token", http.StatusUnauthorized)
				return
			}

			r.Header.Set("X-User-ID", userID)

			next.ServeHTTP(w, r)
		})
	}
}
