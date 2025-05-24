package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handleUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	// create a decoder
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Fatalf("unable to create the user in DB: %s", err)
	}

	respondWithJson(w, 201, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
