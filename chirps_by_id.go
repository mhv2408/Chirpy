package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetChirpsById(w http.ResponseWriter, r *http.Request) {
	chirp_id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot convert chirp id to uuid", err)
		return
	}
	chirp, err := cfg.db.GetChirpById(r.Context(), chirp_id)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirp by id", err)
		return
	}
	respondWithJson(w, http.StatusOK, Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})

}
