package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/coffemanfp/chat/account"
	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/config"
	"github.com/coffemanfp/chat/database"
	sErrors "github.com/coffemanfp/chat/errors"
	"github.com/coffemanfp/chat/server/handlers"
	"github.com/gorilla/mux"
)

type handlerName string

// AuthHandler represents a handler for sign actions like sign up and login.
// Handles external sign platforms and own-server sign service.
type AuthHandler struct {
	config     config.ConfigInfo
	repository database.AuthRepository
	writer     handlers.ResponseWriter
	reader     handlers.RequestReader

	// accountReaders keeps the services to be used for read the account info which is trying to sign.
	accountReaders map[handlerName]accountReader
}

// accountReader represents a service which reads the account info.
type accountReader interface {
	// read reads the account info and return it in a new instance.
	//  @param w http.ResponseWriter: response writer of the call.
	//  @param r *http.Request: request instance of the call.
	//	@return $1 account.Account: new account information instance.
	//	@return $2 error: connection or reading error.
	read(w http.ResponseWriter, r *http.Request) (account.Account, error)
}

// NewAuthHandler initializes a new AuthHandler instance.
//
//	@param repo database.AuthRepository: AuthRepository interface for the authentication handling.
//	@param r handlers.RequestReader: RequestReader interface for reading request operations.
//	@param w handlers.ResponseWriter: ResponseWriter interface for writing request response operations.
//	@param conf config.ConfigInfo: keeps the config information of the service.
//	@return u AuthHandler: new AuthHandler instance.
func NewAuthHandler(repo database.AuthRepository, r handlers.RequestReader, w handlers.ResponseWriter, conf config.ConfigInfo) (u AuthHandler) {
	return AuthHandler{
		reader:     r,
		writer:     w,
		repository: repo,
		config:     conf,
		accountReaders: map[handlerName]accountReader{
			systemHandlerName: systemAccountReader{
				reader: r,
				writer: w,
			},
		},
	}
}

// HandleAuth implements the account authentication actions.
func (a AuthHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling auth...")
	vars := mux.Vars(r)
	action := vars["action"]

	accountReader, err := a.getAccountReader(handlerName("system"))
	if err != nil {
		a.handleError(w, err)
		return
	}

	log.Printf("Sending system sign up")
	account, err := accountReader.read(w, r)
	if err != nil {
		a.handleError(w, err)
		return
	}

	var id, code int

	switch action {
	case "signup":
		id, err = a.handleSignUp(account, w, r)
		if err != nil {
			a.handleError(w, err)
			return
		}
		code = 201
	case "login":
		id, err = a.handleLogin(account, w, r)
		if err != nil {
			a.handleError(w, err)
			return
		}
		code = 200
	}

	token, err := auth.GenerateJWT(a.config.Server.SecretKey, id)
	if err != nil {
		a.handleError(w, err)
		return
	}

	a.writer.JSON(w, code, token)
	log.Println("Success", action)
}

// handleSignUp performs a sign up process for the account requested.
//
//	@param account account.Account: account to sign up.
//	@return id int: account authenticated id.
func (a AuthHandler) handleSignUp(account account.Account, w http.ResponseWriter, r *http.Request) (id int, err error) {
	id, _, err = a.signUp(account)
	return
}

// handleLogin performs a login process for the account requested.
//
//	@param account account.Account: account to login.
//	@return id int: account authenticated id.
func (a AuthHandler) handleLogin(account account.Account, w http.ResponseWriter, r *http.Request) (id int, err error) {
	id, _, err = a.login(account)
	return
}

// signUp perfoms the account sign up process.
//
//	 @param accountR account.Account: account to sign up.
//		@return id int: account authenticated id.
//		@return session auth.Session: new session of the account.
//		@return err error: sign up, validation or connection error
func (a AuthHandler) signUp(accountR account.Account) (id int, session auth.Session, err error) {
	log.Printf("Saving sign up of %s %s", accountR.Nickname, accountR.Email)

	accountR, err = account.New(accountR)
	if err != nil {
		return
	}

	log.Println("New generated account")

	id, err = a.repository.SignUp(accountR, session)
	if err != nil {
		return
	}

	accountR.ID = id

	session, err = auth.NewSession(accountR.ID, "system")
	if err != nil {
		return
	}

	log.Println("New generated session")

	err = a.repository.UpsertSession(session)
	if err != nil {
		return
	}

	log.Println("Successfully registered in database")
	return
}

// login performs the account login process.
//
//	 @param accountR account.Account: account to login.
//		@return id int: account authenticated id.
//		@return session auth.Session: new session of the account.
//		@return err error: login, validation or connection error
func (a AuthHandler) login(accountR account.Account) (id int, session auth.Session, err error) {
	log.Printf("Creating login session of %s %s", accountR.Nickname, accountR.Email)

	accountR.Password, err = auth.HashPassword(accountR.Password)
	if err != nil {
		return
	}

	id, err = a.repository.MatchCredentials(accountR)
	if err != nil {
		return
	}

	if id == 0 {
		err = sErrors.NewClientError(http.StatusUnauthorized, "credentials don't match: invalid credentials of account %s %s", accountR.Nickname, accountR.Email)
		return
	}

	session, err = auth.NewSession(id, "system")
	if err != nil {
		return
	}

	err = a.repository.UpsertSession(session)
	return
}

func (a AuthHandler) getAccountReader(name handlerName) (r accountReader, err error) {
	r, ok := a.accountReaders[name]
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

// CheckAuthHandler handler to check if the account is authenticated for auth-required routes.
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

// NewCheckAuthHandler initializes a new CheckAuthHandler middleware.
func NewCheckAuthHandler(w handlers.ResponseWriter) func(http.Handler) http.Handler {
	return func(n http.Handler) http.Handler {
		return &CheckAuthHandler{
			next:   n,
			writer: w,
		}
	}
}

// CheckNoAuthHandler handler to check if the account is not authenticated for no-required auth routes.
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

// NewCheckNoAuthHandler initializes a new CheckNoAuthHandler middleware.
func NewCheckNoAuthHandler(w handlers.ResponseWriter) func(http.Handler) http.Handler {
	return func(n http.Handler) http.Handler {
		return &CheckNoAuthHandler{
			next:   n,
			writer: w,
		}
	}
}
