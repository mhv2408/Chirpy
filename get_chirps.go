package main

import "net/http"

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {

	dbChirps, err := cfg.db.GetAllChirps(r.Context()) //dbChrips are the Chirp(database.Chirp) format
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}
	apiChirps := make([]Chirp, 0, len(dbChirps)) // slice to hold apiChirps (the format we want)
	for _, c := range dbChirps {
		apiChirps = append(apiChirps, Chirp{
			Id:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserId:    c.UserID,
		})
	}

	respondWithJson(w, http.StatusOK, apiChirps)

}
