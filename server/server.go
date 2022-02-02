package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	srv *http.Server
}

func (s Server) Run() (err error) {
	return s.srv.ListenAndServe()
}

func NewServer(port int, staticFilePath string) (server *Server) {
	r := mux.NewRouter().StrictSlash(false)

	v1R := r.PathPrefix("/api/v1").Subrouter()

	// API
	{
		v1R.HandleFunc("/healthcheck/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "OK")
		}).Methods("GET")

	}

	// Templates
	{
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(staticFilePath)))
	}

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("localhost:%d", port),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	server = &Server{
		srv: srv,
	}
	return server
}
