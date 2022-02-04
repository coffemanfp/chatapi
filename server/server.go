package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/coffemanfp/chat/database"
	"github.com/coffemanfp/chat/server/handlers"
	"github.com/coffemanfp/chat/server/handlers/users"
	muxhandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	srv *http.Server
}

func (s Server) Run() (err error) {
	return s.srv.ListenAndServe()
}

func NewServer(db database.Database, host string, port int) (server *Server, err error) {
	r := mux.NewRouter().StrictSlash(false)
	v1R := r.PathPrefix("/api/v1").Subrouter()

	setUpMiddlewares(r)
	setUpAPIHandlers(r)
	setUpUsersHandlers(v1R, db)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%d", host, port),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	server = &Server{
		srv: srv,
	}
	return
}

func setUpAPIHandlers(r *mux.Router) {
	r.HandleFunc("/healthcheck/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}).Methods("GET")
}

func setUpMiddlewares(r *mux.Router) {
	r.Use(logginMiddleware)
	r.Use(muxhandlers.RecoveryHandler())
}

func setUpUsersHandlers(r *mux.Router, db database.Database) {
	repo, err := database.GetUsersRepository(db.Repositories)
	if err != nil {
		return
	}

	uh := users.NewUsersHandler(
		repo,
		handlers.NewRequestReaderImpl(),
		handlers.NewResponseWriterImpl(),
	)

	r.HandleFunc("/users", uh.HandleSignUp).Methods("POST")
}
