package main

import (
	"log"
	"net/http"
)

func (state *apiConfig) handlerResetVisits(res http.ResponseWriter, req *http.Request) {
	state.fileserverHits.Store(0)
	if state.platform != "dev" {
		respondWithError(res, http.StatusForbidden, http.StatusText(http.StatusForbidden), nil)
	}
	err := state.db.DeleteUsers(req.Context())
	if err != nil {
		log.Printf("error deleting all users . %v", err)
		respondWithError(res, http.StatusInternalServerError, "Couldn't delete users", err)
	}

	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)

}
