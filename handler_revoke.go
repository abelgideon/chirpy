package main

import (
	"net/http"

	"github.com/abelgideon/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	tk, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "authorization header is missing")
		return
	}

	revokeErr := cfg.db.SetRevokedAtAndUpdatedAt(r.Context(), tk)
	if revokeErr != nil {
		respondWithError(w, http.StatusInternalServerError, "could not revoke refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
