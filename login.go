package main

import (
	"encoding/json"
	"net/http"

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
	authUser := AuthUser{}
	authUser.mapDBType(&dbUser)
	respondWithJSON(w, 200, authUser)
}
