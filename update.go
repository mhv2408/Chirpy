package main

import (
	"encoding/json"
	"net/http"

	"github.com/mhv2408/Chirpy/internal/auth"
	"github.com/mhv2408/Chirpy/internal/database"
)

func (cfg *apiConfig) handleUpdate(w http.ResponseWriter, r *http.Request) {

	access_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "cannot find Authorization token in header", err)
		return
	}
	user_id, err := auth.ValidateJWT(access_token, cfg.jwt_secret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to validate the JWT token", err)
		return
	}
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to parse the json to extract new email and password", err)
		return
	}

	new_hashed_password, err := auth.HashPassword(params.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to hash the new_password", err)
		return
	}
	user, err := cfg.db.UpdateUserByID(r.Context(), database.UpdateUserByIDParams{
		ID:             user_id,
		Email:          params.Email,
		HashedPassword: new_hashed_password,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to update the user credentials", err)
		return
	}

	respondWithJson(w, http.StatusOK, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})

}
