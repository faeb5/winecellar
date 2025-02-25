package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/faeb5/winecellar/internal/auth"
	"github.com/faeb5/winecellar/internal/database"
	"github.com/faeb5/winecellar/internal/middleware"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

const (
	accessTokenExpiresIn  = time.Hour
	refreshTokenExpiresIn = 60 * 24 * time.Hour
	userIdHeader          = "X-User-ID"
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

func handleLogin(conf apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		type response struct {
			user
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		}

		decoder := json.NewDecoder(r.Body)
		var params parameters
		if err := decoder.Decode(&params); err != nil {
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
			return
		}

		dbUser, err := conf.dbQueries.GetUserByEmail(r.Context(), params.Email)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err)
			return
		}

		if err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword); err != nil {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err)
			return
		}

		accessToken, err := auth.MakeJWT(dbUser.ID, conf.jwtSecret, accessTokenExpiresIn)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
			return
		}

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
			return
		}

		if _, err := conf.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    dbUser.ID,
			ExpiresAt: time.Now().Add(refreshTokenExpiresIn),
		}); err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
			return
		}

		respondWithJSON(w, http.StatusOK, response{
			user: user{
				ID:        dbUser.ID,
				Email:     dbUser.Email,
				CreatedAt: dbUser.CreatedAt,
				UpdatedAt: dbUser.UpdatedAt,
			},
			Token:        accessToken,
			RefreshToken: refreshToken,
		})
	}
}

func handleRegister(conf apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		type response struct {
			user
		}

		decoder := json.NewDecoder(r.Body)
		var params parameters
		if err := decoder.Decode(&params); err != nil {
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
			return
		}

		// No error means the user already exists
		if _, err := conf.dbQueries.GetUserByEmail(r.Context(), params.Email); err == nil {
			respondWithError(w, http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest), err)
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
			return
		}

		dbUser, err := conf.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
			ID:             uuid.NewString(),
			Email:          params.Email,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
			return
		}

		respondWithJSON(w, http.StatusOK, response{
			user{
				ID:        dbUser.ID,
				Email:     dbUser.Email,
				CreatedAt: dbUser.CreatedAt,
				UpdatedAt: dbUser.UpdatedAt,
			},
		})
	}
}
