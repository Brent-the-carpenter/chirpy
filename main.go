package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/Brent-the-carpenter/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db                 *database.Queries
	platform           string
	secret             string
	polkaApiKey        string
	RefreshTokenExpiry time.Duration
	AccessTokenExpiry  time.Duration
	fileserverHits     atomic.Int32
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

	polkaApikey := os.Getenv("POLKA_API_KEY")
	if polkaApikey == "" {
		log.Fatal("POLKA_API_KEY must be set.")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("Platform must be set.")
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRETE must be set")
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database: %w", err)
	}
	// Get the pointer to the database to access queries made by sqlc
	dbQueries := database.New(dbConn)
	apiCfg := apiConfig{
		fileserverHits:     atomic.Int32{},
		db:                 dbQueries,
		platform:           platform,
		secret:             secret,
		polkaApiKey:        polkaApikey,
		RefreshTokenExpiry: time.Hour * 24 * 60,
		AccessTokenExpiry:  time.Hour,
	}
	// Create server router
	mux := http.NewServeMux()
	// Register handlers for specific routes
	fsHandler := apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadyCheck)

	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeRefreshToken)

	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUpgradeUser)

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
