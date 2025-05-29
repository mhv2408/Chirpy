package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mhv2408/Chirpy/internal/database"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwt_secret     string
	polka_key      string
}

func main() {
	const root_file_path = "."
	const port = "8080"
	godotenv.Load(".env")

	dbURL := os.Getenv("DB_URL") // getting the database url
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	dbConn, err := sql.Open("postgres", dbURL) //opening the connection
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
		return
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("ADMIN_KEY environment variable is not set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable is not set")
	}

	//new config
	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(dbConn), //new database connection
		platform:       platform,
		jwt_secret:     jwtSecret,
		polka_key:      polkaKey,
	}
	mux := http.NewServeMux() //creating a serve multiplexer :- connects request types --> handlers

	//Register the Handler
	mux.HandleFunc("GET /api/healthz", handleReadiness) //only accesssible for GET req

	file_handler := http.FileServer(http.Dir(root_file_path))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", file_handler)))

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)

	mux.HandleFunc("POST /api/users", apiCfg.handleUsers)
	mux.HandleFunc("PUT /api/users", apiCfg.handleUpdate)

	mux.HandleFunc("POST /api/chirps", apiCfg.handleChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleGetChirpsById)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handleDeleteChirp)

	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	// refresh end point
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefresh)

	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevoke)

	// webhooks endpoint
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handleWebHooks)

	srv := &http.Server{Handler: mux, Addr: ":" + port}

	log.Printf("Serving files from %s on port: %s\n", root_file_path, port)

	log.Fatal(http.ListenAndServe(srv.Addr, srv.Handler))
}
