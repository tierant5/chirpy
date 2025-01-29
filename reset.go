package main

import "net/http"

func (cfg *apiConfig) handleReset(w http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	if cfg.platform != "dev" {
		w.WriteHeader(403)
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.dbQueries.DeleteAllUsers(req.Context())
	if err != nil {
		msg := "Error deleting users from database"
		respondWithError(w, 400, msg, err)
		return
	}
	err = cfg.dbQueries.DeleteAllChirps(req.Context())
	if err != nil {
		msg := "Error deleting Chirps from database"
		respondWithError(w, 400, msg, err)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("data reset"))
}
