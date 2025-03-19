package main

import "net/http"

func handleReset(conf apiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := conf.dbQueries.DeleteAllUsers(r.Context()); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to delete users", err)
		}
		if err := conf.dbQueries.DeleteAllWines(r.Context()); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to delete wines", err)
		}
		if err := conf.dbQueries.DeleteAllRatings(r.Context()); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to delete ratings", err)
		}
		if err := conf.dbQueries.DeleteAllRefreshTokens(r.Context()); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to delete refresh tokens", err)
		}
		w.WriteHeader(http.StatusOK)
	}
}
