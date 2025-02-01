package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tierant5/chirpy/internal/auth"
	"github.com/tierant5/chirpy/internal/database"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		msg := "Could not decode request body"
		respondWithError(w, 400, msg, err)
		return
	}
	dbUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		respondUnauthorized(w, err)
		return
	}
	err = auth.CheckPasswordHash(user.Password, dbUser.HashedPassword)
	if err != nil {
		respondUnauthorized(w, err)
		return
	}

	const refreshTokenExpiresHours = 60 * 24 * time.Hour
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		msg := "could not generate refresh token"
		respondWithError(w, 400, msg, err)
		return
	}
	dbRefreshToken, err := cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().UTC().Add(refreshTokenExpiresHours),
	})
	if err != nil {
		msg := "could not generate refresh token in database"
		respondWithError(w, 400, msg, err)
		return
	}

	const expiresInSeconds = 3600
	token, err := auth.MakeJWT(dbUser.ID, cfg.signingToken, time.Second*time.Duration(expiresInSeconds))
	if err != nil {
		msg := "could not generate token"
		respondWithError(w, 400, msg, err)
		return
	}
	authUser := AuthUser{Token: token, RefreshToken: dbRefreshToken.Token}
	authUser.mapDBType(&dbUser)
	respondWithJSON(w, 200, authUser)
}
