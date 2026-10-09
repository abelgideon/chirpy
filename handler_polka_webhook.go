package main

import (
	"encoding/json"
	"net/http"

	"github.com/abelgideon/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, apiErr := auth.GetAPIKey(r.Header)
	if apiErr != nil {
		respondWithError(w, http.StatusUnauthorized, "api key missing")
		return
	}

	if apiKey != cfg.polkakey {
		respondWithError(w, http.StatusUnauthorized, "you are not authorized")
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	param := parameters{}

	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error decoding json")
		return
	}

	if param.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userId, err := uuid.Parse(param.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing id")
		return
	}

	foundErr := cfg.db.GetUserById(r.Context(), userId)
	if foundErr != nil {
		respondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	upgradeErr := cfg.db.UpgradeUserById(r.Context(), userId)
	if upgradeErr != nil {
		respondWithError(w, http.StatusInternalServerError, "error upgrading user")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
