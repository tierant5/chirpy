package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tierant5/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (c *Chirp) mapDBType(d *database.Chirp) {
	c.ID = d.ID
	c.CreatedAt = d.CreatedAt
	c.UpdatedAt = d.UpdatedAt
	c.Body = d.Body
	c.UserID = d.UserID
}

func (cfg *apiConfig) handleChirps(w http.ResponseWriter, r *http.Request) {
	var chirp Chirp
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
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
		UserID: chirp.UserID,
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
