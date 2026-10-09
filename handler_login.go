package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/abelgideon/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		ExpiresAt *int   `json:"expires_in_seconds"`
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
		return
	}

	isValid, err := auth.CheckPasswordHash(param.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error comparing passwords")
		return
	}

	if !isValid {
		respondWithError(w, http.StatusUnauthorized, "you are not authorized")
		return
	}

	var expiresAt int

	if param.ExpiresAt == nil || *param.ExpiresAt*int(time.Second) > int(time.Hour) {
		expiresAt = int(time.Hour)
	} else {
		expiresAt = *param.ExpiresAt * int(time.Second)
	}

	tk, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(expiresAt))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error generating token")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     tk,
	})
}
