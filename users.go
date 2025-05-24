package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mhv2408/Chirpy/internal/auth"
	"github.com/mhv2408/Chirpy/internal/database"
)

func (cfg *apiConfig) handleUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	// create a decoder
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	password_hash, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Fatal("Error generating password hash: ", err)
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: password_hash,
	})
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
