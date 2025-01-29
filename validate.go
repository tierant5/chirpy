package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Valid bool `json:"valid"`
	}
	var params parameters

	contentType, ok := r.Header["Content-Type"]
	if !ok {
		msg := "Couldn't determine Content-Type"
		respondWithError(w, 400, msg, nil)
		return
	}
	if contentType[0] != "application/json" {
		msg := fmt.Sprintf("request must be application/json type, got: %v", contentType[0])
		respondWithError(w, 400, msg, nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		msg := "Malformed request"
		respondWithError(w, 400, msg, err)
		return
	}

	const chirp_len = 140
	if len(params.Body) > chirp_len {
		msg := fmt.Sprintf("Chirp is longer than %v, got %v", chirp_len, len(params.Body))
		respondWithError(w, 400, msg, nil)
		return
	}

	respondWithJSON(w, 200, returnVals{
		Valid: true,
	})
}
