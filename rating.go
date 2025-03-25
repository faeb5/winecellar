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

type updateRatingParameters struct {
	Rating string `json:"rating"`
}

func handleDeleteRating(apiConfig apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(userIdHeader)
		if userID == "" {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized),
				fmt.Errorf("Missing header %s in http request", userIdHeader))
			return
		}

		ratingID := r.PathValue("ratingID")
		if _, err := apiConfig.dbQueries.GetRatingByID(r.Context(), ratingID); err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), err)
			return
		}

		if err := apiConfig.dbQueries.DeleteRatingByID(r.Context(), ratingID); err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func handleGetRatings(apiConfig apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(userIdHeader)
		if userID == "" {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized),
				fmt.Errorf("Missing header %s in http request", userIdHeader))
			return
		}

		dbRatings, err := apiConfig.dbQueries.GetAllRatings(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
		}

		ratings := make([]rating, len(dbRatings))
		for i, dbRating := range dbRatings {
			ratings[i] = rating{
				ID:        dbRating.ID,
				WineID:    dbRating.WineID,
				UserID:    dbRating.UserID,
				Rating:    dbRating.Rating,
				CreatedAt: dbRating.CreatedAt,
				UpdatedAt: dbRating.UpdatedAt,
			}
		}

		respondWithJSON(w, http.StatusOK, ratings)
	}
}

func handleUpdateRating(apiConfig apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(userIdHeader)
		if userID == "" {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized),
				fmt.Errorf("Missing header %s in http request", userIdHeader))
			return
		}

		ratingID := r.PathValue("ratingID")
		if _, err := apiConfig.dbQueries.GetRatingByID(r.Context(), ratingID); err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), err)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var params updateRatingParameters
		if err := decoder.Decode(&params); err != nil {
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
			return
		}

		dbRating, err := apiConfig.dbQueries.UpdateRatingByID(r.Context(), database.UpdateRatingByIDParams{
			ID:     ratingID,
			Rating: params.Rating,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
			return
		}

		respondWithJSON(w, http.StatusOK, rating{
			ID:        dbRating.ID,
			WineID:    dbRating.WineID,
			UserID:    dbRating.UserID,
			Rating:    dbRating.Rating,
			CreatedAt: dbRating.CreatedAt,
			UpdatedAt: dbRating.UpdatedAt,
		})
	}
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
