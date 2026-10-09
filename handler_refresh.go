package main

import (
	"net/http"
	"time"

	"github.com/abelgideon/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	tk, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "authorization header is missing")
		return
	}

	user, err := cfg.db.GetUserByRefreshToken(r.Context(), tk)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "you are not authenticated")
		return
	}

	newToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create new jwt")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": newToken,
	})
}
