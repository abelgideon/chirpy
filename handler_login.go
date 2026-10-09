package main

import (
	"encoding/json"
	"net/http"

	"github.com/abelgideon/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	param := parameters{}

	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding json")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), param.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found")
	}

	isValid, err := auth.CheckPasswordHash(param.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error comparing passwords")
		return
	}

	if !isValid {
		respondWithError(w, http.StatusUnauthorized, "you are not authorized")
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
