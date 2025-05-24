package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mhv2408/Chirpy/internal/database"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func validateChirp(w http.ResponseWriter, message string) {

	if len(message) > 140 { // if the char limit crosses 140...send an 400 response
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

}

func (cfg *apiConfig) handleChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}         //struct for the json body expected to get...
	err := decoder.Decode(&params) //decoding the json into jsonBody struct
	if err != nil {
		log.Fatal("unable to decode the received json: ", err)
	}
	validateChirp(w, params.Body)                 // validation chirp length
	params.Body = removeProfaneWords(params.Body) //removing profane words

	// creating a chirp

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserId,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJson(w, http.StatusCreated, Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})

}
