package main

import (
	"database/sql"
	"net/http"

	"github.com/Brent-the-carpenter/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (state *apiConfig) handlerDeleteChirp(res http.ResponseWriter, req *http.Request) {
	chirpId := req.PathValue("chirpID")
	if chirpId == "" {
		respondWithError(res, http.StatusBadRequest, "chirp id in URL path is missing", nil)
		return
	}

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(res, http.StatusUnauthorized, "couldn't get access token from http.header", err)
		return
	}
	userID, err := auth.ValidateJWT(accessToken, state.secret)
	if err != nil {
		respondWithError(res, http.StatusUnauthorized, "malformed or missing accesstoken", err)
		return
	}

	chirpUUID, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "couldn't parse chirpId to uuid", err)
		return
	}
	chirp, err := state.db.GetChirp(req.Context(), chirpUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(res, http.StatusNotFound, "chirp not found", err)
		} else {
			respondWithError(res, http.StatusInternalServerError, "couldn't get chirp record", err)
		}
		return
	}

	if chirp.UserID != userID {
		respondWithError(res, http.StatusForbidden, "not your chirp", nil)
		return
	}

	err = state.db.DeleteChirp(req.Context(), chirp.ID)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, "couldn't delete chirp", err)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}
