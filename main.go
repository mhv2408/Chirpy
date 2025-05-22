package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func main() {
	const root_file_path = "."
	const port = "8080"

	//new config
	apiCfg := &apiConfig{}

	mux := http.NewServeMux() //creating a serve multiplexer :- connects request types --> handlers

	//Register the Handler
	mux.HandleFunc("GET /api/healthz", handleReadiness) //only accesssible for GET req

	file_handler := http.FileServer(http.Dir(root_file_path))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", file_handler)))

	mux.HandleFunc("GET /api/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /api/reset", apiCfg.handleReset)

	srv := &http.Server{Handler: mux, Addr: ":" + port}

	log.Printf("Serving files from %s on port: %s\n", root_file_path, port)

	log.Fatal(http.ListenAndServe(srv.Addr, srv.Handler))
}
