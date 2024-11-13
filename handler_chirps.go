package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Brent-the-carpenter/chirpy/internal/auth"
	"github.com/Brent-the-carpenter/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirp struct {
	ID         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Body       string    `json:"body"`
	User_id    uuid.UUID `json:"user_id"`
}

func (state *apiConfig) createChirp(res http.ResponseWriter, req *http.Request) {
	type params struct {
		Body string `json:"body"`
		//	UserID uuid.UUID `json:"user_id"`
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

	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(res, http.StatusBadRequest, "Couldn't get bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, state.secret)
	if err != nil {
		respondWithError(res, http.StatusUnauthorized, "Invalid JWT", err)
		return
	}

	badWords := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}
	cleaned := cleanInput(parameters.Body, badWords)
	newChirp, err := state.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "Unable to save chirp to database", err)
		return
	}

	respondWithJSON(res, http.StatusCreated, chirp{
		ID:         newChirp.ID,
		Created_at: newChirp.CreatedAt,
		Updated_at: newChirp.CreatedAt,
		Body:       newChirp.Body,
		User_id:    newChirp.UserID,
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
	authorParam := req.URL.Query().Get("author_id")
	sortingParam := req.URL.Query().Get("sort")

	sort := "asc"
	if sortingParam == "desc" {
		sort = "desc"
	}

	var userID uuid.UUID
	if authorParam != "" {

		var err error
		userID, err = uuid.Parse(authorParam)
		if err != nil {
			respondWithError(res, http.StatusBadRequest, "author_id malformed", err)
			return
		}
	}

	queryMap := map[string]func(ctx context.Context) ([]database.Chirp, error){
		"all_asc":   func(ctx context.Context) ([]database.Chirp, error) { return state.db.GetAllChirpsAsc(ctx) },
		"all_desc":  func(ctx context.Context) ([]database.Chirp, error) { return state.db.GetAllChirpsDesc(ctx) },
		"user_asc":  func(ctx context.Context) ([]database.Chirp, error) { return state.db.GetChirpsForUserAsc(ctx, userID) },
		"user_desc": func(ctx context.Context) ([]database.Chirp, error) { return state.db.GetChirpsForUserDesc(ctx, userID) },
	}

	queryKey := "all_" + sort
	if authorParam != "" {
		queryKey = "user_" + sort
	}

	chirpRecords, err := queryMap[queryKey](req.Context())
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "Database query failed", err)
		return
	}

	chirps := make([]chirp, len(chirpRecords))
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
	urlParam := req.PathValue("chirpID")
	if urlParam == "" {
		respondWithError(res, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil)
		return
	}
	chirpID, err := uuid.Parse(urlParam)
	if err != nil {
		respondWithError(res, 500, "Invalid chirp id", err)
		return
	}
	c, err := state.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(res, 404, "chirp not found", nil)
		return
	}

	respondWithJSON(res, 200, chirp{
		ID:         c.ID,
		Created_at: c.CreatedAt,
		Updated_at: c.UpdatedAt,
		Body:       c.Body,
		User_id:    c.UserID,
	})
}
