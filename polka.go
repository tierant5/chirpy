package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/tierant5/chirpy/internal/auth"
	"github.com/tierant5/chirpy/internal/database"
)

type PolkaWebhook struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "no api key found", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, 401, "invalid api key", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var polkaWebhook PolkaWebhook
	err = decoder.Decode(&polkaWebhook)
	if err != nil {
		msg := "error decoding webhook body"
		respondWithError(w, 400, msg, err)
		return
	}

	if polkaWebhook.Event != "user.upgraded" {
		respondWithJSON(w, 204, struct{}{})
		return
	}

	err = cfg.dbQueries.UpdateUserChirpyRed(r.Context(), database.UpdateUserChirpyRedParams{
		ID:          polkaWebhook.Data.UserID,
		IsChirpyRed: true,
	})
	if err != nil {
		respondWithError(w, 404, "user not found", err)
		return
	}

	respondWithJSON(w, 204, struct{}{})
}
