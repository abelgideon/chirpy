package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	param := parameters{}

	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, 400, "error decoding json")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), param.Email)
	if err != nil {
		respondWithError(w, 500, "error creating user")
		return
	}

	respondWithJSON(w, 201, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
