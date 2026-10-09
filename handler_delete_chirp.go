package main

import (
	"net/http"

	"github.com/abelgideon/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	tk, bearerErr := auth.GetBearerToken(r.Header)
	if bearerErr != nil {
		respondWithError(w, http.StatusUnauthorized, "authorization header missing")
		return
	}

	userID, valErr := auth.ValidateJWT(tk, cfg.secret)
	if valErr != nil {
		respondWithError(w, http.StatusUnauthorized, "you are not authenticated")
		return
	}

	pathParam := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(pathParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id")
		return
	}

	chirp, foundErr := cfg.db.GetChirp(r.Context(), chirpID)
	if foundErr != nil {
		respondWithError(w, http.StatusNotFound, "could not find chirp")
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "chirp does not belong to you")
		return
	}

	deleteErr := cfg.db.DeleteChirp(r.Context(), chirpID)
	if deleteErr != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to delete chirp")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
