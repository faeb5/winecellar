package main

import (
	"database/sql"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/faeb5/winecellar/internal/database"
	"github.com/faeb5/winecellar/internal/middleware"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

const userIdHeader = "X-User-ID"

type apiConfig struct {
	profile   string
	port      string
	jwtSecret string
	dbQueries *database.Queries
}

func main() {
	envFile := flag.String("e", ".env", "The environment file")
	flag.Parse()

	if err := godotenv.Load(*envFile); err != nil {
		log.Fatal(err)
	}

	apiConfig, err := createApiConfig()
	if err != nil {
		log.Fatal(err)
	}

	// API routes
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(http.StatusText(http.StatusOK))); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	apiMux.HandleFunc("GET /wines", handleGetWines(apiConfig))
	apiMux.HandleFunc("POST /wines", handleCreateWine(apiConfig))
	apiMux.HandleFunc("GET /wines/{wineID}", handleGetWine(apiConfig))
	apiMux.HandleFunc("PUT /wines/{wineID}", handleUpdateWine(apiConfig))
	apiMux.HandleFunc("DELETE /wines/{wineID}", handleDeleteWine(apiConfig))

	apiMux.HandleFunc("GET /ratings", handleGetRatings(apiConfig))
	apiMux.HandleFunc("POST /ratings", handleCreateRating(apiConfig))
	apiMux.HandleFunc("GET /ratings/{ratingID}", handleGetRating(apiConfig))
	apiMux.HandleFunc("PUT /ratings/{ratingID}", handleUpdateRating(apiConfig))
	apiMux.HandleFunc("DELETE /ratings/{ratingID}", handleDeleteRating(apiConfig))

	apiStack := middleware.CreateStack(middleware.Authorized(apiConfig.jwtSecret))

	// DEV routes
	devMux := http.NewServeMux()
	devMux.HandleFunc("POST /reset", handleReset(apiConfig))
	devStack := middleware.CreateStack(middleware.DevOnly(apiConfig.profile))

	// default routes
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", apiStack(apiMux)))
	mux.Handle("/dev/", http.StripPrefix("/dev", devStack(devMux)))
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
	profile := os.Getenv("PROFILE")
	if profile == "" {
		return apiConfig{}, errors.New("PROFILE is not set")
	}
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
		profile:   profile,
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
