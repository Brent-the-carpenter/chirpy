package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Brent-the-carpenter/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (state *apiConfig) handlerUpgradeUser(res http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(res, http.StatusUnauthorized, "error getting api key", err)
		return
	}
	if apiKey != state.polkaApiKey {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := parameters{}
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "couldn't unmarshal JSON", err)
		return
	}

	if params.Event != "user.upgraded" {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	_, err = state.db.UpgradeUser(req.Context(), params.Data.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	respondWithJSON(res, http.StatusNoContent, nil)
}
