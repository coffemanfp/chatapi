package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coffemanfp/chat/config"
	"github.com/coffemanfp/chat/users"
	"golang.org/x/oauth2"
)

const facebookHandlerName handlerName = "facebook"

type facebookHandler struct {
	oauth oauth
	name  handlerName
}

func (f facebookHandler) read(w http.ResponseWriter, r *http.Request) (user users.User, err error) {
	resp, err := f.oauth.callback(w, r)
	if err != nil {
		return
	}
	user, err = f.parseUserInfo(resp)
	return
}

func (f facebookHandler) requestSignUp(w http.ResponseWriter, r *http.Request) (err error) {
	err = f.oauth.redirectToHandler(w, r)
	return
}

func (f facebookHandler) parseUserInfo(b []byte) (user users.User, err error) {
	userInfo := struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}{}

	err = json.Unmarshal(b, &userInfo)
	if err != nil {
		err = fmt.Errorf("failed to read user data from %s: %s", f.name, err)
		return
	}

	user.SignedWith = append(user.SignedWith, users.ExternalSigned{
		ID:        userInfo.ID,
		Email:     userInfo.Email,
		Platform:  f.name.string(),
		CreatedAt: time.Now(),
		Picture:   userInfo.Picture.Data.URL,
	})
	return
}

func newFacebookHandler(conf config.ConfigInfo) facebookHandler {
	return facebookHandler{
		oauth: oauth{
			conf: &oauth2.Config{
				ClientID:     conf.OAuth.Facebook.ClientID,
				ClientSecret: conf.OAuth.Facebook.ClientSecret,
				Endpoint:     conf.OAuth.Facebook.Endpoint,
				RedirectURL:  conf.OAuth.Facebook.RedirectURIS[0],
				Scopes:       conf.OAuth.Facebook.Scopes,
			},
			handler:     facebookHandlerName,
			validStates: make(map[string]bool),
		},
		name: facebookHandlerName,
	}
}
