package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Brent-the-carpenter/chirpy/internal/auth"
	"github.com/Brent-the-carpenter/chirpy/internal/database"
)

func (state *apiConfig) handleLogin(res http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type loginResponse struct {
		User
	}

	defer req.Body.Close()
	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(res, http.StatusBadRequest, "Either email or password was left out of request", nil)
		return
	}

	user, err := state.db.GetUser(req.Context(), params.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(res, http.StatusNotFound, "User not found", err)
		} else {
			respondWithError(res, http.StatusInternalServerError, "Failed to retrieve user", err)
		}
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(res, http.StatusUnauthorized, "Incorrect password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, state.secret, state.AccessTokenExpiry)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "Error making JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "couldn't make refresh token", err)
		return
	}

	_, err = state.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
		RevokedAt: sql.NullTime{},
	})
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "Couldn't save refresh token to db", err)
		return
	}
	response := loginResponse{
		User: User{
			ID:           user.ID,
			Email:        user.Email,
			Token:        token,
			RefreshToken: refreshToken,
			ChirpyRed:    user.IsChirpyRed.Bool,
		},
	}

	respondWithJSON(res, http.StatusOK, response)

}
