package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/faeb5/winecellar/internal/database"
	"github.com/google/uuid"
)

type createWineParameters struct {
	Name     string `json:"name"`
	Color    string `json:"color"`
	Producer string `json:"producer"`
	Country  string `json:"country"`
	Vintage  int    `json:"vintage"`
}

type updateWineParameters struct {
	Name     string `json:"name"`
	Color    string `json:"color"`
	Producer string `json:"producer"`
	Country  string `json:"country"`
	Vintage  int    `json:"vintage"`
}

func handleGetWines(apiConfig apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbWines, err := apiConfig.dbQueries.GetAllWines(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
			return
		}

		wines := make([]wine, len(dbWines))
		for i, dbWine := range dbWines {
			wines[i] = wine{
				ID:        dbWine.ID,
				Name:      dbWine.Name,
				Color:     dbWine.Color,
				Producer:  dbWine.Producer,
				Country:   dbWine.Country,
				Vintage:   int(dbWine.Vintage),
				CreatedAt: dbWine.CreatedAt,
				UpdatedAt: dbWine.UpdatedAt,
			}
		}

		respondWithJSON(w, http.StatusOK, wines)
	}
}

func handleUpdateWine(apiConfig apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wineID := r.PathValue("wineID")
		if _, err := apiConfig.dbQueries.GetWineByID(r.Context(), wineID); err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), err)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var params updateWineParameters
		if err := decoder.Decode(&params); err != nil {
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
			return
		}

		dbWine, err := apiConfig.dbQueries.UpdateWineByID(r.Context(), database.UpdateWineByIDParams{
			Color:    params.Color,
			Name:     params.Name,
			Producer: params.Producer,
			Country:  params.Country,
			Vintage:  int64(params.Vintage),
			ID:       wineID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
			return
		}

		respondWithJSON(w, http.StatusOK, wine{
			ID:        dbWine.ID,
			Name:      dbWine.Name,
			Color:     dbWine.Color,
			Producer:  dbWine.Producer,
			Country:   dbWine.Country,
			Vintage:   int(dbWine.Vintage),
			CreatedAt: dbWine.CreatedAt,
			UpdatedAt: dbWine.UpdatedAt,
		})
	}
}

func handleGetWine(apiConfig apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wineID := r.PathValue("wineID")
		dbWine, err := apiConfig.dbQueries.GetWineByID(r.Context(), wineID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), err)
			return
		}

		respondWithJSON(w, http.StatusOK, wine{
			ID:        dbWine.ID,
			Name:      dbWine.Name,
			Color:     dbWine.Color,
			Producer:  dbWine.Producer,
			Country:   dbWine.Country,
			Vintage:   int(dbWine.Vintage),
			CreatedAt: dbWine.CreatedAt,
			UpdatedAt: dbWine.UpdatedAt,
		})
	}
}

func handleDeleteWine(apiConfig apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wineID := r.PathValue("wineID")
		if _, err := apiConfig.dbQueries.GetWineByID(r.Context(), wineID); err != nil {
			respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound), err)
			return
		}

		if err := apiConfig.dbQueries.DeleteWine(r.Context(), wineID); err != nil {
			respondWithError(w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func handleCreateWine(apiConfig apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var params createWineParameters
		if err := decoder.Decode(&params); err != nil {
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
			return
		}

		_, err := apiConfig.dbQueries.GetWineByProducerAndNameAndVintage(
			r.Context(),
			database.GetWineByProducerAndNameAndVintageParams{
				Producer: params.Producer,
				Name:     params.Name,
				Vintage:  int64(params.Vintage),
			},
		)
		// No error means the wine already exists
		if err == nil {
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err)
			return
		}

		dbWine, err := apiConfig.dbQueries.CreateWine(r.Context(), database.CreateWineParams{
			ID:        uuid.NewString(),
			Name:      params.Name,
			Color:     params.Color,
			Producer:  params.Producer,
			Country:   params.Country,
			Vintage:   int64(params.Vintage),
			CreatedBy: r.Header.Get(userIdHeader),
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError),
				err)
			return
		}

		respondWithJSON(w, http.StatusCreated, wine{
			ID:        dbWine.ID,
			Name:      dbWine.Name,
			Color:     dbWine.Color,
			Producer:  dbWine.Producer,
			Country:   dbWine.Country,
			Vintage:   int(dbWine.Vintage),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
}
