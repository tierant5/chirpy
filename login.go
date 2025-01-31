package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tierant5/chirpy/internal/auth"
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

	var expiresInSeconds int
	if user.ExpiresInSeconds != "" {
		expiresInSeconds, err = strconv.Atoi(user.ExpiresInSeconds)
		if err != nil {
			msg := fmt.Sprintf("could not parse int: %v", user.ExpiresInSeconds)
			respondWithError(w, 400, msg, err)
			return
		}
		if expiresInSeconds > 3600 {
			expiresInSeconds = 3600
		}
	} else {
		expiresInSeconds = 3600
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.signingToken, time.Second*time.Duration(expiresInSeconds))
	if err != nil {
		msg := "could not generate token"
		respondWithError(w, 400, msg, err)
		return
	}
	authUser := AuthUser{Token: token}
	authUser.mapDBType(&dbUser)
	respondWithJSON(w, 200, authUser)
}
