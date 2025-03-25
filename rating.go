package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/faeb5/winecellar/internal/database"
	"github.com/google/uuid"
)

type createRatingParameters struct {
	WineID string `json:"wine_id"`
	Rating string `json:"rating"`
}

func handleCreateRating(apiConfig apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(userIdHeader)
		if userID == "" {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized),
				fmt.Errorf("Missing header %s in http request", userIdHeader))
			return
		}

		dbUser, err := apiConfig.dbQueries.GetUserByID(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var params createRatingParameters
		if err := decoder.Decode(&params); err != nil {
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
			return
		}

		dbWine, err := apiConfig.dbQueries.GetWineByID(r.Context(), params.WineID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
			return
		}

		dbRating, err := apiConfig.dbQueries.CreateRating(r.Context(), database.CreateRatingParams{
			ID:     uuid.NewString(),
			WineID: dbWine.ID,
			UserID: dbUser.ID,
			Rating: params.Rating,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
			return
		}

		respondWithJSON(w, http.StatusCreated, rating{
			ID:        dbRating.ID,
			WineID:    dbRating.WineID,
			UserID:    dbRating.UserID,
			Rating:    dbRating.Rating,
			CreatedAt: dbRating.CreatedAt,
			UpdatedAt: dbRating.UpdatedAt,
		})
	}
}
