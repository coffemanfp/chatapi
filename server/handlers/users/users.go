package users

import (
	"net/http"

	"github.com/coffemanfp/chat/database"
	"github.com/coffemanfp/chat/models"
	"github.com/coffemanfp/chat/server/handlers"
)

type UsersHandler struct {
	Repository database.UsersRepository
	Writer     handlers.ResponseWriter
	Reader     handlers.RequestReader
}

func (u UsersHandler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := u.Reader.JSON(r, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Writer.JSON(w, handlers.Hash{
			"code":    404,
			"message": err.Error(),
		})
		return
	}

	user, err = models.NewUser(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Writer.JSON(w, handlers.Hash{
			"code":    404,
			"message": err.Error(),
		})
		return
	}

	id, err := u.Repository.SignUp(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Writer.JSON(w, handlers.Hash{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	user.ID = id
	user.Password = ""

	err = u.Writer.JSON(w, user)
	if err != nil {
		return
	}
}

func NewUsersHandler(db database.UsersRepository, r handlers.RequestReader, w handlers.ResponseWriter) (u UsersHandler) {
	return UsersHandler{
		Reader:     r,
		Writer:     w,
		Repository: db,
	}
}
