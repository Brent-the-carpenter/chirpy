package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Brent-the-carpenter/chirpy/internal/database"
	"github.com/google/uuid"
)

func (state *apiConfig) createChirp(res http.ResponseWriter, req *http.Request) {
	type params struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type returnVals struct {
		ID         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		User_id    uuid.UUID `json:"user_id"`
	}

	parameters := params{}
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&parameters)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(res, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(parameters.Body) > maxChirpLength {
		respondWithError(res, 400, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}
	cleaned := cleanInput(parameters.Body, badWords)
	chirp, err := state.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: parameters.UserID,
	})
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "Unable to save chirp to database", err)
	}

	respondWithJSON(res, http.StatusCreated, returnVals{
		ID:         chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.CreatedAt,
		Body:       chirp.Body,
		User_id:    chirp.UserID,
	})
}

const maxChirpLength = 140

func cleanInput(s string, badwords map[string]struct{}) string {
	splitString := strings.Split(s, " ")
	for index, word := range splitString {
		lowerCaseWord := strings.ToLower(word)
		if _, exist := badwords[lowerCaseWord]; exist {
			splitString[index] = "****"
		}
	}

	return strings.Join(splitString, " ")
}

func (state *apiConfig) handlerGetAllChirps(res http.ResponseWriter, req *http.Request) {
	type chirp struct {
		ID         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		User_id    uuid.UUID `json:"user_id"`
	}

	type returnVals []chirp

	chirpRecords, err := state.db.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(res,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err)
	}
	chirps := make(returnVals, len(chirpRecords))
	for index, c := range chirpRecords {
		chirps[index] = chirp{
			ID:         c.ID,
			Created_at: c.CreatedAt,
			Updated_at: c.UpdatedAt,
			Body:       c.Body,
			User_id:    c.UserID,
		}
	}
	respondWithJSON(res, 200, chirps)

}

func (state *apiConfig) handlerGetChirp(res http.ResponseWriter, req *http.Request) {
	type chirp struct {
		ID         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		User_id    uuid.UUID `json:"user_id"`
	}

	urlParam := req.PathValue("chirpID")
	if urlParam == "" {
		respondWithError(res, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil)
	}
	chirpID, err := uuid.Parse(urlParam)
	if err != nil {
		respondWithError(res, 500, "Invalid chirp id", err)
	}
	c, err := state.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(res, 404, "chirp not found", nil)
	}

	respondWithJSON(res, 200, chirp{
		ID:         c.ID,
		Created_at: c.CreatedAt,
		Updated_at: c.UpdatedAt,
		Body:       c.Body,
		User_id:    c.UserID,
	})
}
