package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/mhv2408/Chirpy/internal/auth"
)

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	access_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find the JWT", err)
		return
	}
	user_id, err := auth.ValidateJWT(access_token, cfg.jwt_secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate the JWT", err)
		return
	}
	chirp_id, err := uuid.Parse(r.PathValue("chirpID"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	chirp_details, err := cfg.db.GetChirpById(r.Context(), chirp_id)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not get the chirp", err)
		return
	}
	if chirp_details.UserID != user_id {
		respondWithError(w, http.StatusForbidden, "User does not have permission to perform the operation", err)
		return
	}
	err = cfg.db.DeleteChirpByID(r.Context(), chirp_details.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete the chirp", err)
		return
	}
	respondWithJson(w, http.StatusNoContent, nil)
}
