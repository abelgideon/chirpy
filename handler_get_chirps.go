package main

import (
	"net/http"

	"github.com/abelgideon/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userIDQuery := query.Get("author_id")

	var chirps []database.Chirp
	var err error

	if userIDQuery != "" {
		userID, err := uuid.Parse(userIDQuery)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error parsing author id")
			return
		}

		chirps, err = cfg.db.GetChirpsByAuthor(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error fetching author's chirps")
			return
		}
	} else {
		chirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error fetching chirps")
			return
		}
	}

	chirpsList := []Chirp{}

	for _, chirp := range chirps {
		chirpsList = append(chirpsList, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, 200, chirpsList)
}
