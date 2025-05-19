package main

import (
	"log"
	"net/http"
)

func main() {
	const root_file_path = "."
	const port = "8080"
	mux := http.NewServeMux() //creating a serve multiplexer :- connects request types --> handlers

	file_handler := http.FileServer(http.Dir(root_file_path))
	mux.Handle("/", file_handler)

	srv := &http.Server{Handler: mux, Addr: ":" + port}

	log.Printf("Serving files from %s on port: %s\n", root_file_path, port)

	log.Fatal(http.ListenAndServe(srv.Addr, srv.Handler))
}
