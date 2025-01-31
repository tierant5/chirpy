package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tierant5/chirpy/internal/auth"
	"github.com/tierant5/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type Chirps []Chirp

func (c Chirps) mapDBType(d *[]database.Chirp) {
	chirp := Chirp{}
	for i, dbChirp := range *d {
		chirp.mapDBType(&dbChirp)
		c[i] = chirp
	}
}

func (c *Chirp) mapDBType(d *database.Chirp) {
	c.ID = d.ID
	c.CreatedAt = d.CreatedAt
	c.UpdatedAt = d.UpdatedAt
	c.Body = d.Body
	c.UserID = d.UserID
}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		msg := "token not found"
		respondWithError(w, 400, msg, err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.signingToken)
	if err != nil {
		msg := "invalid token"
		respondWithError(w, 401, msg, err)
		return
	}

	var chirp Chirp
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&chirp)
	if err != nil {
		msg := "Error decoding body"
		respondWithError(w, 400, msg, err)
		return
	}
	cleanedBody, ok := chirp.Validate()
	if !ok {
		respondWithError(w, 400, cleanedBody, nil)
		return
	}
	dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: userID,
		Body:   cleanedBody,
	})
	if err != nil {
		msg := "Error creating Chirp in database"
		respondWithError(w, 400, msg, err)
		return
	}
	chirp.mapDBType(&dbChirp)
	respondWithJSON(w, 201, chirp)
}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		msg := "Error getting Chirps from database"
		respondWithError(w, 400, msg, err)
	}
	chirps := make(Chirps, len(dbChirps))
	chirps.mapDBType(&dbChirps)
	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) handleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		msg := "Error parsing chirpID"
		respondWithError(w, 400, msg, err)
	}
	dbChirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		msg := fmt.Sprintf("Error getting chirpID: %v from the database", chirpID)
		respondWithError(w, 404, msg, err)
	}
	chirp := Chirp{}
	chirp.mapDBType(&dbChirp)
	respondWithJSON(w, 200, chirp)
}
