package main

import (
	"fmt"
	"log"
	"net/http"
)

func (state *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		state.fileserverHits.Add(1)
		next.ServeHTTP(res, req)

	})
}

func (state *apiConfig) handlerNumOfVisits(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	_, err := res.Write([]byte(fmt.Sprintf(`
		<html>
		  <body>
		    <h1>Welcome, Chirpy Admin</h1>
		    <p>Chirpy has been visited %d times!</p>
		  </body>
		</html>`, state.fileserverHits.Load())))
	if err != nil {
		log.Printf("Error responding to request at %v: %v", req.URL.Path, err)
		return
	}
}
