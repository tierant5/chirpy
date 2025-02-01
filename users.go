package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tierant5/chirpy/internal/auth"
	"github.com/tierant5/chirpy/internal/database"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type AuthUser struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (u *User) mapDBType(d *database.User) {
	u.ID = d.ID
	u.CreatedAt = d.CreatedAt
	u.UpdatedAt = d.UpdatedAt
	u.Email = d.Email
	u.IsChirpyRed = d.IsChirpyRed
}

func (u *AuthUser) mapDBType(d *database.User) {
	u.ID = d.ID
	u.CreatedAt = d.CreatedAt
	u.UpdatedAt = d.UpdatedAt
	u.Email = d.Email
	u.IsChirpyRed = d.IsChirpyRed
}

func (cfg *apiConfig) handleUsersPost(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		msg := "Could not decode request body"
		respondWithError(w, 400, msg, err)
		return
	}
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		msg := "Error creating hashedPassword"
		respondWithError(w, 400, msg, err)
		return
	}
	dbUser, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          user.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		msg := "Error creating user in database"
		respondWithError(w, 400, msg, err)
		return
	}
	user.mapDBType(&dbUser)
	respondWithJSON(w, 201, user)
}

func (cfg *apiConfig) handleUsersPut(w http.ResponseWriter, r *http.Request) {
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

	var user User
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&user)
	if err != nil {
		msg := "Could not decode request body"
		respondWithError(w, 400, msg, err)
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		msg := "Error creating hashedPassword"
		respondWithError(w, 400, msg, err)
		return
	}

	dbUser, err := cfg.dbQueries.UpdateUserEmailPassword(r.Context(), database.UpdateUserEmailPasswordParams{
		ID:             userID,
		Email:          user.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		msg := "Error updating user in database"
		respondWithError(w, 400, msg, err)
		return
	}
	user.mapDBType(&dbUser)
	respondWithJSON(w, 200, user)
}
