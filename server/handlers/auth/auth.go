package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/config"
	"github.com/coffemanfp/chat/database"
	sErrors "github.com/coffemanfp/chat/errors"
	"github.com/coffemanfp/chat/server/handlers"
	"github.com/coffemanfp/chat/users"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	config               config.ConfigInfo
	repository           database.AuthRepository
	writer               handlers.ResponseWriter
	reader               handlers.RequestReader
	userReaders          map[string]userReader
	externalSignHandlers map[string]externalSignUpHandler
}

type userReader interface {
	readUser(w http.ResponseWriter, r *http.Request) (users.User, error)
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
		userReaders: map[string]userReader{
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

func (a AuthHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling auth...")
	vars := mux.Vars(r)
	action := vars["action"]

	if action == "external-sign" {
		a.HandleExternalSign(w, r)
		return
	}

	hName := vars["handler"]
	if hName == "" {
		hName = "system"
	}

	h, err := a.getUserReader(hName)
	if err != nil {
		a.handleError(w, err)
		return
	}

	log.Printf("Sending %s sign up", hName)
	user, err := h.readUser(w, r)
	if err != nil {
		a.handleError(w, err)
		return
	}

	var tmpID string

	switch action {
	case "signup":
		tmpID, err = a.handleSignUp(user, w, r)
	case "login":
		tmpID, err = a.handleLogin(user, w, r)
	}
	if err != nil {
		a.handleError(w, err)
		return
	}

	rURL, _ := url.Parse("http://localhost:3000/chat")
	qValues := rURL.Query()
	qValues.Set("tmp_id", tmpID)
	rURL.RawQuery = qValues.Encode()

	if hName != "system" {
		http.Redirect(w, r, rURL.String(), http.StatusTemporaryRedirect)
	}

	log.Printf("Success %s %s", hName, action)
}

func (a AuthHandler) handleSignUp(user users.User, w http.ResponseWriter, r *http.Request) (tmpID string, err error) {
	user, session, err := a.signUp(user)
	if err != nil {
		return
	}

	tmpID = session.TmpID
	return
}

func (a AuthHandler) handleLogin(user users.User, w http.ResponseWriter, r *http.Request) (tmpID string, err error) {
	session, err := a.login(user)
	if err != nil {
		return
	}

	tmpID = session.TmpID
	return
}

func (a AuthHandler) HandleExternalSign(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling external logging...")

	vars := mux.Vars(r)
	hName := vars["handler"]

	h, err := a.getExternalSignHandler(hName)
	if err != nil {
		a.handleError(w, err)
		return
	}

	log.Printf("Handling %s logging...", hName)

	err = h.requestSignUp(w, r)
	if err != nil {
		a.handleError(w, err)
		return
	}

	log.Printf("Successfully redirected to %s", hName)
}

func (a AuthHandler) signUp(userR users.User) (user users.User, session auth.Session, err error) {
	log.Printf("Saving sign up of %s %s", userR.Nickname, userR.Email)

	userR, err = users.New(userR)
	if err != nil {
		return
	}

	log.Println("New generated user")

	platform := "system"
	if len(userR.SignedWith) > 0 {
		platform = userR.SignedWith[0].Platform
	}

	id, err := a.repository.SignUp(userR, session)
	if err != nil {
		return
	}

	userR.ID = id

	session, err = auth.NewSession(userR.ID, platform)
	if err != nil {
		return
	}

	log.Println("New generated session")

	err = a.repository.UpsertSession(session)
	if err != nil {
		return
	}

	log.Println("Successfully registered in database")

	user = userR
	user.ID = id
	user.Password = ""
	return
}

func (a AuthHandler) login(userR users.User) (session auth.Session, err error) {
	log.Printf("Creating login session of %s %s", userR.Nickname, userR.Email)

	err = users.HashPassword(&userR.Password)
	if err != nil {
		return
	}

	id, err := a.repository.MatchCredentials(userR)
	if err != nil {
		return
	}

	if id == 0 {
		err = sErrors.NewClientError(http.StatusUnauthorized, "credentials don't match: invalid credentials of user %s %s", userR.Nickname, userR.Email)
		return
	}

	platform := "system"
	if len(userR.SignedWith) > 0 {
		platform = userR.SignedWith[0].Platform
	}

	session, err = auth.NewSession(userR.ID, platform)
	if err != nil {
		return
	}

	err = a.repository.UpsertSession(session)
	return
}

func (a AuthHandler) getExternalSignHandler(name string) (h externalSignUpHandler, err error) {
	h, ok := a.externalSignHandlers[name]
	if !ok {
		err = sErrors.NewClientError(http.StatusBadRequest, "invalid callback handler: %s not exists", name)
	}
	return
}

func (a AuthHandler) getUserReader(name string) (r userReader, err error) {
	r, ok := a.userReaders[name]
	if !ok {
		err = sErrors.NewClientError(http.StatusBadRequest, "invalid signup handler: %s not exists", name)
	}
	return
}

func (a AuthHandler) handleError(w http.ResponseWriter, err error) {
	hErr, ok := err.(sErrors.ClientError)
	if !ok {
		log.Println(err)
		a.writer.JSON(w, http.StatusInternalServerError, handlers.Hash{
			"message": sErrors.SERVER_ERROR_MESSAGE,
		})
		return
	}
	a.writer.JSON(w, hErr.HTTPCode(), handlers.Hash{
		"message": hErr.Error(),
	})
}

type CheckAuthHandler struct {
	next   http.Handler
	writer handlers.ResponseWriter
}

func (c CheckAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			// not authorized
			w.Header().Set("Location", "/")
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			log.Printf("failed to load authentication cookie: %s", err)
			c.writer.JSON(w, http.StatusInternalServerError, handlers.Hash{
				"message": sErrors.SERVER_ERROR_MESSAGE,
			})
		}
		return
	}

	c.next.ServeHTTP(w, r)
}

func NewCheckAuthHandler(w handlers.ResponseWriter) func(http.Handler) http.Handler {
	return func(n http.Handler) http.Handler {
		return &CheckAuthHandler{
			next:   n,
			writer: w,
		}
	}
}

type CheckNoAuthHandler struct {
	next   http.Handler
	writer handlers.ResponseWriter
}

func (c CheckNoAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err != nil && err != http.ErrNoCookie {
		log.Printf("failed to load authentication cookie: %s", err)
		c.writer.JSON(w, http.StatusInternalServerError, handlers.Hash{
			"message": sErrors.SERVER_ERROR_MESSAGE,
		})
		return
	}

	w.Header().Set("Location", "/chat")
	w.WriteHeader(http.StatusTemporaryRedirect)

	c.next.ServeHTTP(w, r)
}

func NewNoCheckAuthHandler(w handlers.ResponseWriter) func(http.Handler) http.Handler {
	return func(n http.Handler) http.Handler {
		return &CheckNoAuthHandler{
			next:   n,
			writer: w,
		}
	}
}
