package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Brent-the-carpenter/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	const (
		port         = "8080"
		filepathRoot = "."
	)

	//Load enviroment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load enviroment variable. %w", err)
	}

	// Connect to database
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("Platform must be set.")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database: %w", err)
	}
	// Get the pointer to the database to access queries made by sqlc
	dbQueries := database.New(dbConn)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}
	// Create server router
	mux := http.NewServeMux()
	// Register handlers for specific routes
	fsHandler := apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadyCheck)

	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerNumOfVisits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetVisits)
	// Make server with router set up previously and port address to host on
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	// Start server wrap in log fatal so if server exits with error it gets logged.
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
