package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/faeb5/winecellar/internal/auth"
	"github.com/faeb5/winecellar/internal/database"
)

const (
	accessTokenExpiresIn  = time.Hour
	refreshTokenExpiresIn = 60 * 24 * time.Hour
)

type loginParameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginResponse struct {
	user
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func handleLogin(conf apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var params loginParameters
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

		respondWithJSON(w, http.StatusOK, loginResponse{
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
