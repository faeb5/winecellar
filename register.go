package main

import (
	"encoding/json"
	"net/http"

	"github.com/faeb5/winecellar/internal/auth"
	"github.com/faeb5/winecellar/internal/database"
	"github.com/google/uuid"
)

type registerParameters struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type registerResponse struct {
	user
}

func handleRegister(conf apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var params registerParameters
		if err := decoder.Decode(&params); err != nil {
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
			return
		}

		// No error means the e-mail address already exists
		if _, err := conf.dbQueries.GetUserByUsername(r.Context(), params.Username); err == nil {
			respondWithError(w, http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest), err)
			return
		}

		// No error means the e-mail address already exists
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
			Username:       params.Username,
			Email:          params.Email,
			HashedPassword: hashedPassword,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
			return
		}

		respondWithJSON(w, http.StatusOK, registerResponse{
			user{
				ID:        dbUser.ID,
				Username:  dbUser.Username,
				Email:     dbUser.Email,
				CreatedAt: dbUser.CreatedAt,
				UpdatedAt: dbUser.UpdatedAt,
			},
		})
	}
}
