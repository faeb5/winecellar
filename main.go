package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"
	log.Println("Listening on port", port)
	server := http.Server{
		Addr: ":" + port,
	}

	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	log.Fatal(server.ListenAndServe())
}
