package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
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
	queries := r.URL.Query()
	sortby := queries.Get("sort")
	author_id := queries.Get("author_id")
	if author_id != "" {
		author_id, err := uuid.Parse(author_id)
		if err != nil {
			msg := "Error parsing author_id"
			respondWithError(w, 400, msg, err)
			return
		}
		dbChirps, err := cfg.dbQueries.GetAllChirpsByAuthor(r.Context(), author_id)
		if err != nil {
			msg := "Error getting Chirps from database"
			respondWithError(w, 400, msg, err)
			return
		}
		chirps := make(Chirps, len(dbChirps))
		chirps.mapDBType(&dbChirps)
		if sortby == "desc" {
			sort.Slice(chirps, func(i int, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
		}
		respondWithJSON(w, 200, chirps)
		return
	} else {
		dbChirps, err := cfg.dbQueries.GetAllChirps(r.Context())
		if err != nil {
			msg := "Error getting Chirps from database"
			respondWithError(w, 400, msg, err)
			return
		}
		chirps := make(Chirps, len(dbChirps))
		chirps.mapDBType(&dbChirps)
		if sortby == "desc" {
			sort.Slice(chirps, func(i int, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
		}
		respondWithJSON(w, 200, chirps)
		return
	}
}

func (cfg *apiConfig) handleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		msg := "Error parsing chirpID"
		respondWithError(w, 400, msg, err)
		return
	}
	dbChirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		msg := fmt.Sprintf("Error getting chirpID: %v from the database", chirpID)
		respondWithError(w, 404, msg, err)
		return
	}
	chirp := Chirp{}
	chirp.mapDBType(&dbChirp)
	respondWithJSON(w, 200, chirp)
}

func (cfg *apiConfig) handleDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		msg := "token not found"
		respondWithError(w, 401, msg, err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.signingToken)
	if err != nil {
		msg := "invalid token"
		respondWithError(w, 401, msg, err)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		msg := "Error parsing chirpID"
		respondWithError(w, 400, msg, err)
		return
	}
	dbChirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		msg := fmt.Sprintf("Error getting chirpID: %v from the database", chirpID)
		respondWithError(w, 404, msg, err)
		return
	}

	if userID != dbChirp.UserID {
		msg := "user is not the owner"
		respondWithError(w, 403, msg, nil)
		return
	}

	err = cfg.dbQueries.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		msg := "error removing chirp from database"
		respondWithError(w, 400, msg, err)
		return
	}
	respondWithJSON(w, 204, struct{}{})
}
