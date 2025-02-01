package main

import (
	"net/http"
	"time"

	"github.com/tierant5/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		msg := "refresh token not found"
		respondWithError(w, 400, msg, err)
		return
	}
	dbRefreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		msg := "refresh token not found"
		respondWithError(w, 401, msg, err)
		return
	}
	is_expired := dbRefreshToken.ExpiresAt.UTC().Before(time.Now().UTC())
	if is_expired {
		msg := "refresh token expired"
		respondWithError(w, 401, msg, err)
		return
	}
	if dbRefreshToken.RevokedAt.Valid {
		msg := "refresh token revoked"
		respondWithError(w, 401, msg, err)
		return
	}

	const expiresInSeconds = 3600
	token, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.signingToken, time.Second*time.Duration(expiresInSeconds))
	if err != nil {
		msg := "could not generate token"
		respondWithError(w, 400, msg, err)
		return
	}
	type response struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, 200, response{
		Token: token,
	})
}
