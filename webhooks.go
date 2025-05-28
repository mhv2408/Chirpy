package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/mhv2408/Chirpy/internal/auth"
)

func (cfg *apiConfig) handleWebHooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	// validate the API key
	api_key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key", err)
		return
	}
	if api_key != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "API key is invalid", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot parse the json body", err)
		return
	}

	if params.Event != "user.upgraded" { //do nothing
		respondWithJson(w, http.StatusNoContent, nil)
		return
	}

	// upgrade the user subscription
	err = cfg.db.UpgradeUserToChirpyRed(r.Context(), uuid.MustParse(params.Data.UserID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	respondWithJson(w, http.StatusNoContent, nil)
}
