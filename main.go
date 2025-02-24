package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/faeb5/winecellar/internal/database"
	"github.com/faeb5/winecellar/internal/middleware"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

type apiConfig struct {
	port      string
	jwtSecret string
	dbQueries *database.Queries
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	apiConfig, err := createApiConfig()
	if err != nil {
		log.Fatal(err)
	}

	// API routes
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	apiStack := middleware.CreateStack(middleware.Authorized(apiConfig.jwtSecret))

	// default routes
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", apiStack(apiMux)))

	mux.HandleFunc("POST /login", handleLogin(apiConfig))
	mux.HandleFunc("POST /register", handleRegister(apiConfig))

	defaultStack := middleware.CreateStack(middleware.Logging)

	server := http.Server{
		Addr:    ":" + apiConfig.port,
		Handler: defaultStack(mux),
	}
	log.Println("Listening on port", apiConfig.port)
	log.Fatal(server.ListenAndServe())
}

func createApiConfig() (apiConfig, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return apiConfig{}, errors.New("PORT is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return apiConfig{}, errors.New("JWT_SECRET is not set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return apiConfig{}, errors.New("DB_URL is not set")
	}

	dbQueries, err := openDatabase(dbURL)
	if err != nil {
		return apiConfig{}, err
	}

	config := apiConfig{
		port:      port,
		jwtSecret: jwtSecret,
		dbQueries: dbQueries,
	}

	return config, nil
}

func openDatabase(dbURL string) (*database.Queries, error) {
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return database.New(db), nil
}
