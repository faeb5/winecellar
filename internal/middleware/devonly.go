package middleware

import (
	"net/http"
)

func DevOnly(currentPlatform string) Middleware {
	return func(next http.Handler) http.Handler {
		platform := currentPlatform
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if platform != "dev" {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
