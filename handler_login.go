package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Brent-the-carpenter/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (state *apiConfig) handleLogin(res http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type returnVal struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	defer req.Body.Close()
	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&params)

	if params.Email == "" || params.Password == "" {
		respondWithError(res, http.StatusBadRequest, "Either email or password was left out of request", nil)
		return
	}

	user, err := state.db.GetUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "Failed to retrieve user", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(res, 401, "Incorrect password", err)
		return
	}
	resVal := returnVal{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(res, 200, resVal)

}
