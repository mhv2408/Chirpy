package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/mhv2408/Chirpy/internal/database"
)

func GetApiChirps(dbChirps []database.Chirp) []Chirp {
	resApiChirps := make([]Chirp, 0, len(dbChirps))

	for _, c := range dbChirps {
		resApiChirps = append(resApiChirps, Chirp{
			Id:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserId:    c.UserID,
		})
	}
	return resApiChirps
}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("author_id")

	if id != "" {
		author_id, err := uuid.Parse(id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "cannot parse the id into UUID format", err)
			return
		}
		dbAuthorChirps, err := cfg.db.GetChirpsByUserID(r.Context(), author_id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirps for User", err)
			return
		}
		apiChirps := GetApiChirps(dbAuthorChirps)
		respondWithJson(w, http.StatusOK, apiChirps)
		return

	}

	dbChirps, err := cfg.db.GetAllChirps(r.Context()) //dbChrips are the Chirp(database.Chirp) format
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}
	apiChirps := GetApiChirps(dbChirps)

	sorting_method := r.URL.Query().Get("sort")
	if sorting_method == "desc" { // sorting by desc order of created_at...by defualt it is ASC
		sort.Slice(apiChirps, func(i, j int) bool { return apiChirps[i].CreatedAt.After(apiChirps[j].CreatedAt) })
	}
	respondWithJson(w, http.StatusOK, apiChirps)

}
