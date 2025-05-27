package main

import (
	"net/http"

	"github.com/mhv2408/Chirpy/internal/auth"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	//find the refresh token
	refresh_token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to extract refresh token from header", err)
		return
	}

	// revoke the refersh_token

	err = cfg.db.RevokeToken(r.Context(), refresh_token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to extract refresh token from header", err)
		return
	}

	respondWithJson(w, http.StatusNoContent, nil)
}
