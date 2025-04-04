package main

import (
	"encoding/json"
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

func handleGetRating(apiConfig apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ratingID := r.PathValue("ratingID")
		dbRating, err := apiConfig.dbQueries.GetRatingByID(r.Context(), ratingID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), err)
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

func handleDeleteRating(apiConfig apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbUser, err := apiConfig.dbQueries.GetUserByID(r.Context(), r.Header.Get(userIdHeader))
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err)
			return
		}

		dbRating, err := apiConfig.dbQueries.GetRatingByID(r.Context(), r.PathValue("ratingID"))
		if err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), err)
			return
		}

		if dbRating.UserID != dbUser.ID {
			respondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden), err)
			return
		}

		if err := apiConfig.dbQueries.DeleteRatingByID(r.Context(), dbRating.ID); err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func handleGetRatings(apiConfig apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		dbUser, err := apiConfig.dbQueries.GetUserByID(r.Context(), r.Header.Get(userIdHeader))
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), err)
			return
		}

		dbRating, err := apiConfig.dbQueries.GetRatingByID(r.Context(), r.PathValue("ratingID"))
		if err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), err)
			return
		}

		if dbRating.UserID != dbUser.ID {
			respondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden), err)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var params updateRatingParameters
		if err := decoder.Decode(&params); err != nil {
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
			return
		}

		newDBRating, err := apiConfig.dbQueries.UpdateRatingByID(r.Context(), database.UpdateRatingByIDParams{
			ID:     dbRating.ID,
			Rating: params.Rating,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
			return
		}

		respondWithJSON(w, http.StatusOK, rating{
			ID:        newDBRating.ID,
			WineID:    newDBRating.WineID,
			UserID:    newDBRating.UserID,
			Rating:    newDBRating.Rating,
			CreatedAt: newDBRating.CreatedAt,
			UpdatedAt: newDBRating.UpdatedAt,
		})
	}
}

func handleCreateRating(apiConfig apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbUser, err := apiConfig.dbQueries.GetUserByID(r.Context(), r.Header.Get(userIdHeader))
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
