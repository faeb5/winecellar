package main

import (
	"log"
	"net/http"

	"github.com/faeb5/winecellar/internal/middleware"
)

func main() {
	port := "8080"
	jwtSecret := "mySecret"
	platform := "dev"

	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.DevOnly(platform),
		middleware.Authorized(jwtSecret),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	log.Println("Listening on port", port)
	server := http.Server{
		Addr:    ":" + port,
		Handler: stack(mux),
	}
	log.Fatal(server.ListenAndServe())
}
