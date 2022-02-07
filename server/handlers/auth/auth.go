package auth

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/config"
	"github.com/coffemanfp/chat/database"
	"github.com/coffemanfp/chat/server/handlers"
	"github.com/coffemanfp/chat/users"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	config     config.ConfigInfo
	repository database.AuthRepository
	writer     handlers.ResponseWriter
	reader     handlers.RequestReader
}

func NewAuthHandler(repo database.AuthRepository, r handlers.RequestReader, w handlers.ResponseWriter, config config.ConfigInfo) (u AuthHandler) {
	return AuthHandler{
		reader:     r,
		writer:     w,
		repository: repo,
		config:     config,
	}
}

func (a AuthHandler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.Llongfile)
	defer log.Println(nil)
	vars := mux.Vars(r)
	handler := vars["handler"]

	switch handler {
	case "system":
		a.systemLoggin(w, r)
	case "google":
		a.oauthLogin(w, r, &oauth2.Config{
			ClientID:     a.config.OAuth.Google.ClientID,
			ClientSecret: a.config.OAuth.Google.ClientSecret,
			Endpoint:     a.config.OAuth.Google.Endpoint,
			RedirectURL:  a.config.OAuth.Google.RedirectURIS[0],
			Scopes:       a.config.OAuth.Google.Scopes,
		}, a.config.OAuth.Google.State)
	default:
		a.writer.JSON(w, http.StatusBadRequest, handlers.Hash{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("invalid signup login: %s not exists", handler),
		})
	}
}

func (a AuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	handler := vars["handler"]
	fmt.Println("////// Metodo:", r.Method)

	switch handler {
	case "google":
		a.oauthCallback(w, r, &oauth2.Config{
			ClientID:     a.config.OAuth.Google.ClientID,
			ClientSecret: a.config.OAuth.Google.ClientSecret,
			Endpoint:     a.config.OAuth.Google.Endpoint,
			RedirectURL:  a.config.OAuth.Google.RedirectURIS[0],
			Scopes:       a.config.OAuth.Google.Scopes,
		}, a.config.OAuth.Google.State, handler)
	default:
		a.writer.JSON(w, http.StatusBadRequest, handlers.Hash{
			"code":    http.StatusBadRequest,
			"message": fmt.Sprintf("invalid signup login: %s not exists", handler),
		})
	}
}

func (a AuthHandler) systemLoggin(w http.ResponseWriter, r *http.Request) {
	var user users.User
	err := a.reader.JSON(r, &user)
	if err != nil {
		a.writer.JSON(w, http.StatusBadRequest, handlers.Hash{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	user, err = users.New(user)
	if err != nil {
		a.writer.JSON(w, http.StatusBadRequest, handlers.Hash{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	session, err := auth.NewSession(user.Nickname)
	if err != nil {
		a.writer.JSON(w, http.StatusInternalServerError, handlers.Hash{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	id, err := a.repository.SignUp(user, session)
	if err != nil {
		a.writer.JSON(w, http.StatusInternalServerError, handlers.Hash{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	user.ID = id
	user.Password = ""

	a.writer.JSON(w, http.StatusCreated, user)
}

func (a AuthHandler) oauthLogin(w http.ResponseWriter, r *http.Request, oauthConf *oauth2.Config, oauthState string) {
	log.SetFlags(log.Llongfile)
	defer log.Println(nil)

	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		a.writer.JSON(w, http.StatusInternalServerError, handlers.Hash{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	fmt.Println(URL.String())

	parameters := url.Values{}

	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthState)

	URL.RawQuery = parameters.Encode()
	url := URL.String()
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (a AuthHandler) oauthCallback(w http.ResponseWriter, r *http.Request, oauthConf *oauth2.Config, oauthState, handler string) {
	fmt.Printf("Callback %s...", handler)

	state := r.FormValue("state")
	fmt.Println(state)
	if state != oauthState {
		a.writer.JSON(w, http.StatusInternalServerError, handlers.Hash{
			"code":    500,
			"message": "failed oauth callback error",
		})
		fmt.Printf("invalid oauth state: expected %s, got %s\n", oauthState, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	fmt.Println("oauth code:", code)

	if code == "" {
		w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission.."))
		}
		return
	}
	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("oauth exchange failed with: %s\n", err)
		a.writer.JSON(w, http.StatusInternalServerError, handlers.Hash{
			"code":    500,
			"message": "failed oauth callback error",
		})
		return
	}
	fmt.Println("TOKEN>> AccessToken>> " + token.AccessToken)
	fmt.Println("TOKEN>> Expiration Time>> " + token.Expiry.String())
	fmt.Println("TOKEN>> RefreshToken>> " + token.RefreshToken)

	resp, err := http.Get(genURLToRequestUserInfo(token.AccessToken, handler))
	if err != nil {
		fmt.Printf("failed to request user info: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read user info: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Println("parseResponseBody: " + string(response) + "\n")

	w.Write([]byte("Hello, I'm protected\n"))
	w.Write([]byte(string(response)))
}

func genURLToRequestUserInfo(accessToken, handler string) (s string) {
	switch handler {
	case "google":
		s = fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", url.QueryEscape(accessToken))
	}
	return
}
