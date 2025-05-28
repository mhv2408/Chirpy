package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mhv2408/Chirpy/internal/auth"
	"github.com/mhv2408/Chirpy/internal/database"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password  string        `json:"password"`
		Email     string        `json:"email"`
		ExpiresIn time.Duration `json:"expires_in"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Fatal("unable to decode the received json: ", err)
	}
	/*
		if params.ExpiresIn == 0 { // set the default expiration value to 1 hour
			params.ExpiresIn = time.Duration(3600) * time.Second // one hour in seconds
		}*/
	params.ExpiresIn = time.Duration(3600) * time.Second // defualt access token time

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

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to make a refresh token", err)
		return
	}

	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to create a refresh token", err)
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refresh_token,
	})

}
