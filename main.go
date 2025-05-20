package main

import (
	"log"
	"net/http"
)

func main() {
	const root_file_path = "."
	const port = "8080"
	mux := http.NewServeMux() //creating a serve multiplexer :- connects request types --> handlers

	//Register the Handler
	mux.HandleFunc("/healthz", handleReadiness)

	file_handler := http.FileServer(http.Dir(root_file_path))
	mux.Handle("/app/", http.StripPrefix("/app/", file_handler))

	srv := &http.Server{Handler: mux, Addr: ":" + port}

	log.Printf("Serving files from %s on port: %s\n", root_file_path, port)

	log.Fatal(http.ListenAndServe(srv.Addr, srv.Handler))
}
func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
