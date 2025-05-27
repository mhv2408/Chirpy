package main

import (
	"net/http"
	"time"

	"github.com/mhv2408/Chirpy/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	// get the refresh token
	refresh_token, err := auth.GetBearerToken(r.Header)

	type response struct {
		Token string `json:"token"`
	}

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to extract refresh token from header", err)
		return
	}
	db_ref_token, err := cfg.db.GetUserFromRefreshToken(r.Context(), refresh_token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to get user_id from refresh_token DB", err)
		return
	}

	if db_ref_token.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token is revoked", err)
		return
	}

	access_token, err := auth.MakeJWT(db_ref_token.UserID, cfg.jwt_secret, time.Duration(3600)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to create a access-token", err)
	}

	respondWithJson(w, http.StatusOK, response{
		Token: access_token,
	})

}
