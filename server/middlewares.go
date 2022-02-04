package server

import (
	"net/http"
	"os"

	muxhandlers "github.com/gorilla/handlers"
)

func logginMiddleware(next http.Handler) http.Handler {
	return muxhandlers.LoggingHandler(os.Stdout, next)
}
