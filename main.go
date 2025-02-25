package main

import (
	"log"
	"net/http"
	"os"

	"github.com/faeb5/winecellar/internal/middleware"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not set")
	}

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
		Addr:    ":" + port,
		Handler: defaultStack(mux),
	}
	log.Println("Listening on port", port)
	log.Fatal(server.ListenAndServe())
}
