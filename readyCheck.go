package main

import (
	"log"
	"net/http"
)

func handlerReadyCheck(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	_, err := res.Write([]byte("OK"))
	if err != nil {
		log.Printf("Error responding to request: %s", err)
	}
}
