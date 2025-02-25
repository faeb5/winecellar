package main

import (
	"log"
	"net/http"

	"github.com/faeb5/winecellar/internal/middleware"
)

func main() {
	log.Println("Listening on port 8080")

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	apiStack := middleware.CreateStack(middleware.Authorized("my-little-secret"))

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", apiStack(apiMux)))
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login successful"))
	})

	defaultStack := middleware.CreateStack(middleware.Logging)

	server := http.Server{
		Addr:    ":8080",
		Handler: defaultStack(mux),
	}
	log.Fatal(server.ListenAndServe())
}
