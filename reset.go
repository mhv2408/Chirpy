package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	if err := cfg.db.DeleteUsers(r.Context()); err != nil {
		log.Fatalf("cannot delete the users in reset endpoint: %s", err)
	}
}
