package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/faeb5/winecellar/internal/auth"
)

type Middleware func(http.Handler) http.Handler

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func CreateStack(mws ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			mw := mws[i]
			next = mw(next)
		}
		return next
	}
}

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

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(sw, r)
		log.Println(sw.statusCode, r.Method, r.URL.Path, time.Until(start))
	})
}
