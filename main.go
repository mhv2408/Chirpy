package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	_ "github.com/lib/pq"
	"github.com/mhv2408/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	const root_file_path = "."
	const port = "8080"
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
	}

	mux := http.NewServeMux() //creating a serve multiplexer :- connects request types --> handlers

	//Register the Handler
	mux.HandleFunc("GET /api/healthz", handleReadiness) //only accesssible for GET req

	file_handler := http.FileServer(http.Dir(root_file_path))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", file_handler)))

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)

	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	srv := &http.Server{Handler: mux, Addr: ":" + port}

	log.Printf("Serving files from %s on port: %s\n", root_file_path, port)

	log.Fatal(http.ListenAndServe(srv.Addr, srv.Handler))
}
