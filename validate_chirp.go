package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}         //struct for the json body expected to get...
	err := decoder.Decode(&params) //decoding the json into jsonBody struct
	if err != nil {
		log.Fatal("unable to decode the received json: ", err)
	}
	if len(params.Body) > 140 { // if the char limit crosses 140...send an 400 response
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}
	params.Body = removeProfaneWords(params.Body)
	respondWithJson(w, 200, returnVals{
		CleanedBody: params.Body,
	})

}
