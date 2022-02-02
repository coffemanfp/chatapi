package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const PORT = 8080
const STATIC_FILES_PATH = "./public"

func main() {
	fmt.Println("Starting...")
	server := newServer(PORT)

	fmt.Printf("Listening on port: %d\n", PORT)
	log.Fatal(server.ListenAndServe())
}

func newServer(port int) (srv *http.Server) {
	r := mux.NewRouter().StrictSlash(false)

	v1R := r.PathPrefix("/api/v1").Subrouter()

	// API
	{
		v1R.HandleFunc("/healthcheck/{id:[0-9]+}", handleHealthCheck)
	}

	// Templates
	{
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(STATIC_FILES_PATH)))
	}

	srv = &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("localhost:%d", PORT),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	return
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
