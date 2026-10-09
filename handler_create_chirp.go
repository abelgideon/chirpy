package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/abelgideon/chirpy/internal/auth"
	"github.com/abelgideon/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	tk, bearerErr := auth.GetBearerToken(r.Header)
	if bearerErr != nil {
		respondWithError(w, http.StatusUnauthorized, "you are not authenticated")
		return
	}

	userID, valErr := auth.ValidateJWT(tk, cfg.secret)
	if valErr != nil {
		respondWithError(w, http.StatusUnauthorized, "you are not authenticated")
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	wordList := strings.Split(params.Body, " ")
	cleaned := []string{}

	for _, word := range wordList {
		lowerWord := strings.ToLower(word)
		if lowerWord == "kerfuffle" || lowerWord == "sharbert" || lowerWord == "fornax" {
			cleaned = append(cleaned, "****")
		} else {
			cleaned = append(cleaned, word)
		}
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   strings.Join(cleaned, " "),
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating chirp")
		return
	}

	respondWithJSON(w, 201, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
