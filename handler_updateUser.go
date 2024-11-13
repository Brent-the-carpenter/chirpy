package main

import (
	"encoding/json"
	"net/http"

	"github.com/Brent-the-carpenter/chirpy/internal/auth"
	"github.com/Brent-the-carpenter/chirpy/internal/database"
)

func (state *apiConfig) handlerUpdateUser(res http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(res, http.StatusUnauthorized, "couldn't get Bearer token", err)
		return
	}

	userId, err := auth.ValidateJWT(accessToken, state.secret)
	if err != nil {
		respondWithError(res, http.StatusUnauthorized, "access token missing or malformed", err)
		return
	}

	defer req.Body.Close()
	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "couldn't decode JSON", err)
		return
	}

	newPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "couldn't hash new password", err)
		return
	}

	user, err := state.db.UpdatePasswordAndEmail(req.Context(), database.UpdatePasswordAndEmailParams{
		Email:          params.Email,
		HashedPassword: newPassword,
		ID:             userId,
	})
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "couldn't update user", err)
		return
	}

	respondWithJSON(res, http.StatusOK, response{
		User{
			ID:        userId,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			ChirpyRed: user.IsChirpyRed.Bool,
		},
	})

}
