package main

import (
	"encoding/json"
	"net/http"

	"github.com/abelgideon/chirpy/internal/auth"
	"github.com/abelgideon/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	tk, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "authorization header missing")
		return
	}

	userID, valErr := auth.ValidateJWT(tk, cfg.secret)
	if valErr != nil {
		respondWithError(w, http.StatusUnauthorized, "you are not authenticated")
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	param := parameters{}

	decodeErr := decoder.Decode(&param)
	if decodeErr != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding json")
		return
	}

	newPass, err := auth.HashPassword(param.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing new password")
		return
	}

	newUser, err := cfg.db.UpdateUserById(r.Context(), database.UpdateUserByIdParams{
		Email:          param.Email,
		HashedPassword: newPass,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating user")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	})
}
