package main

import (
	"net/http"

	"github.com/tierant5/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		msg := "refresh token not found"
		respondWithError(w, 400, msg, err)
		return
	}

	err = cfg.dbQueries.SetRefreshTokenRevoked(r.Context(), refreshToken)
	if err != nil {
		msg := "could not revoke refresh token"
		respondWithError(w, 400, msg, err)
		return
	}

	respondWithJSON(w, 204, struct{}{})
}
