package auth

import (
	"net/http"

	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/database"
	"github.com/coffemanfp/chat/server/handlers"
	"github.com/coffemanfp/chat/users"
)

type AuthHandler struct {
	Repository database.AuthRepository
	Writer     handlers.ResponseWriter
	Reader     handlers.RequestReader
}

func NewAuthHandler(db database.AuthRepository, r handlers.RequestReader, w handlers.ResponseWriter) (u AuthHandler) {
	return AuthHandler{
		Reader:     r,
		Writer:     w,
		Repository: db,
	}
}

func (a AuthHandler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	var user users.User
	err := a.Reader.JSON(r, &user)
	if err != nil {
		a.Writer.JSON(w, http.StatusBadRequest, handlers.Hash{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	user, err = users.New(user)
	if err != nil {
		a.Writer.JSON(w, http.StatusBadRequest, handlers.Hash{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	session, err := auth.NewSession(user.Nickname)
	if err != nil {
		a.Writer.JSON(w, http.StatusInternalServerError, handlers.Hash{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	id, err := a.Repository.SignUp(user, session)
	if err != nil {
		a.Writer.JSON(w, http.StatusInternalServerError, handlers.Hash{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	user.ID = id
	user.Password = ""

	a.Writer.JSON(w, http.StatusCreated, user)
}
