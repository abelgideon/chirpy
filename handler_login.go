package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/abelgideon/chirpy/internal/auth"
	"github.com/abelgideon/chirpy/internal/database"
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

	tk, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(time.Hour))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error generating token")
		return
	}

	refToken := auth.MakeRefreshToken()
	genRefTk, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error generating refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     tk,
		RefToken:  genRefTk.Token,
	})
}
