package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mhv2408/Chirpy/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password  string        `json:"password"`
		Email     string        `json:"email"`
		ExpiresIn time.Duration `json:"expires_in_seconds,omitempty"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Fatal("unable to decode the received json: ", err)
	}
	if params.ExpiresIn == 0 { // set the default expiration value to 1 hour
		params.ExpiresIn = time.Duration(3600) * time.Second // one hour in seconds
	}
	user, err := cfg.db.UserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve user by email", err)
	}
	if err := auth.CheckPasswordHash(user.HashedPassword, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwt_secret, params.ExpiresIn)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to create a token", err)
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})

}
