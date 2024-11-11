package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(res http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errResposne struct {
		Error string `json:"error"`
	}

	errRes := errResposne{Error: msg}

	respondWithJSON(res, code, errRes)
}

func respondWithJSON(res http.ResponseWriter, code int, payload interface{}) {
	res.Header().Set("Content-Type", "application/json")
	payloadJSON, err := json.Marshal(&payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		res.WriteHeader(500)
		return
	}
	res.WriteHeader(code)
	res.Write(payloadJSON)

}
