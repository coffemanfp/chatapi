package auth

import (
	"fmt"
	"net/http"

	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/config"
	"github.com/coffemanfp/chat/database"
	"github.com/coffemanfp/chat/server/handlers"
	"github.com/coffemanfp/chat/users"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	config               config.ConfigInfo
	repository           database.AuthRepository
	writer               handlers.ResponseWriter
	reader               handlers.RequestReader
	signHandlers         map[string]signUpHandler
	externalSignHandlers map[string]externalSignUpHandler
}

type signUpHandler interface {
	getUser(w http.ResponseWriter, r *http.Request) (users.User, error)
}

type externalSignUpHandler interface {
	requestSignUp(w http.ResponseWriter, r *http.Request) error
}

func NewAuthHandler(repo database.AuthRepository, r handlers.RequestReader, w handlers.ResponseWriter, conf config.ConfigInfo) (u AuthHandler) {
	fbHandler := newFacebookHandler(conf)
	gHandler := newGoogleHandler(conf)
	return AuthHandler{
		reader:     r,
		writer:     w,
		repository: repo,
		config:     conf,
		signHandlers: map[string]signUpHandler{
			"system": systemSignUpHandler{
				reader: r,
				writer: w,
			},
			"google":   gHandler,
			"facebook": fbHandler,
		},
		externalSignHandlers: map[string]externalSignUpHandler{
			"google":   gHandler,
			"facebook": fbHandler,
		},
	}
}

func (a AuthHandler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling sign up...")
	vars := mux.Vars(r)
	hName := vars["handler"]
	if hName == "" {
		hName = "system"
	}

	h, ok := a.signHandlers[hName]
	if !ok {
		a.writer.JSON(w, http.StatusBadRequest, handlers.Hash{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("invalid signup handler: %s not exists", hName),
		})
	}

	fmt.Printf("Sending %s sign up\n", hName)
	user, err := h.getUser(w, r)
	if err != nil {
		return
	}

	a.writer.JSON(w, 200, user)

	user, sessionID := a.signUp(user, w, r)
	if err != nil {
		return
	}

	a.writer.JSON(w, 200, handlers.Hash{
		"session_id": sessionID,
		"user":       user,
	})
	fmt.Printf("Success %s sign up", hName)
}

func (a AuthHandler) HandleExternalSign(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling external logging...")

	vars := mux.Vars(r)
	hName := vars["handler"]
	h, ok := a.externalSignHandlers[hName]
	if !ok {
		a.writer.JSON(w, http.StatusBadRequest, handlers.Hash{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("invalid callback handler: %s not exists", hName),
		})
		return
	}
	fmt.Printf("Handling %s logging...\n", hName)

	err := h.requestSignUp(w, r)
	if err != nil {
		a.writer.JSON(w, http.StatusBadRequest, handlers.Hash{
			"code":    http.StatusInternalServerError,
			"message": "failed to request a callback request",
		})
	}
	fmt.Printf("Successfully redirected to %s", hName)
}

func (a AuthHandler) signUp(userR users.User, w http.ResponseWriter, r *http.Request) (user users.User, sessionID string) {
	userR, err := users.New(userR)
	if err != nil {
		a.writer.JSON(w, http.StatusBadRequest, handlers.Hash{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	session, err := auth.NewSession(userR.Nickname, userR.SignedWith[0].Platform)
	if err != nil {
		a.writer.JSON(w, http.StatusInternalServerError, handlers.Hash{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	id, err := a.repository.SignUp(userR, session)
	if err != nil {
		a.writer.JSON(w, http.StatusInternalServerError, handlers.Hash{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	sessionID = session.ID
	user = userR
	user.ID = id
	user.Password = ""
	return
}
