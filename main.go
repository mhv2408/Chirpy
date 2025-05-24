package main

import (
	"database/sql"
	"fmt"
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
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	const root_file_path = "."
	const port = "8080"
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")               // getting the database url
	dbConn, err := sql.Open("postgres", dbURL) //opening the connection
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
		return
	}

	//new config
	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(dbConn), //new database connection
		platform:       os.Getenv("PLATFORM"),
	}
	mux := http.NewServeMux() //creating a serve multiplexer :- connects request types --> handlers

	//Register the Handler
	mux.HandleFunc("GET /api/healthz", handleReadiness) //only accesssible for GET req

	file_handler := http.FileServer(http.Dir(root_file_path))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", file_handler)))

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)

	mux.HandleFunc("POST /api/users", apiCfg.handleUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.handleChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleGetChirpsById)

	srv := &http.Server{Handler: mux, Addr: ":" + port}

	log.Printf("Serving files from %s on port: %s\n", root_file_path, port)

	log.Fatal(http.ListenAndServe(srv.Addr, srv.Handler))
}
