package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tierant5/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (u *User) mapDBType(d *database.User) {
	u.ID = d.ID
	u.CreatedAt = d.CreatedAt
	u.UpdatedAt = d.UpdatedAt
	u.Email = d.Email
}

func (cfg *apiConfig) handleUsers(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		msg := "Could not decode request body"
		respondWithError(w, 400, msg, err)
		return
	}
	dbUser, err := cfg.dbQueries.CreateUser(r.Context(), user.Email)
	if err != nil {
		msg := "Error creating user in database"
		respondWithError(w, 400, msg, err)
		return
	}
	user.mapDBType(&dbUser)
	respondWithJSON(w, 201, user)
}
