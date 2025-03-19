package middleware

import (
	"net/http"
)

func DevOnly(currentProfile string) Middleware {
	return func(next http.Handler) http.Handler {
		profile := currentProfile
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if profile != "dev" {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
