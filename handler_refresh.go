package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Brent-the-carpenter/chirpy/internal/auth"
)

func (state *apiConfig) handlerRefresh(res http.ResponseWriter, req *http.Request) {
	type refreshResponse struct {
		Token string `json:"token"`
	}
	header := req.Header
	bearerToken, err := auth.GetBearerToken(header)
	if err != nil {
		respondWithError(res,
			http.StatusBadRequest,
			"Couldn't get refresh token from request headers",
			err)
		return
	}

	refreshToken, err := state.db.GetRefreshToken(req.Context(), bearerToken)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(res, http.StatusUnauthorized, "refresh token doesn't exsist", err)
		} else {
			respondWithError(res, http.StatusUnauthorized, http.StatusText(http.StatusInternalServerError), err)
		}
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(res, http.StatusUnauthorized, fmt.Sprintf("token expired at %v", refreshToken.ExpiresAt), nil)
		return
	}
	if refreshToken.RevokedAt.Valid {
		respondWithError(res, http.StatusUnauthorized, "refresher token has been revoked", nil)
		return
	}
	authToken, err := auth.MakeJWT(refreshToken.UserID, state.secret, state.AccessTokenExpiry)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "couldn't create new JWT token", err)
		return
	}
	respondWithJSON(res, http.StatusOK, refreshResponse{Token: authToken})

}

func (state *apiConfig) handlerRevokeRefreshToken(res http.ResponseWriter, req *http.Request) {

	header := req.Header
	refreshToken, err := auth.GetBearerToken(header)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "couldn't get refresh token from auth header", err)
		return
	}
	_, err = state.db.RevokeRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "couldn't revoke token", err)
		return
	}
	res.WriteHeader(http.StatusNoContent)
}
