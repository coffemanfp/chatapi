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

const googleHandlerName handlerName = "google"

type googleHandler struct {
	oauth oauth
	name  handlerName
}

func (g googleHandler) read(w http.ResponseWriter, r *http.Request) (user users.User, err error) {
	resp, err := g.oauth.callback(w, r)
	if err != nil {
		return
	}
	user, err = g.parseUserInfo(resp)
	return
}

func (g googleHandler) requestSignUp(w http.ResponseWriter, r *http.Request) (err error) {
	err = g.oauth.redirectToHandler(w, r)
	return
}

func (g googleHandler) parseUserInfo(b []byte) (user users.User, err error) {
	userInfo := struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}{}

	err = json.Unmarshal(b, &userInfo)
	if err != nil {
		err = fmt.Errorf("failed to read user data from %s: %s", g.name, err)
		return
	}

	user = users.User{
		SignedWith: []users.ExternalSigned{
			{
				ID:        userInfo.ID,
				Email:     userInfo.Email,
				Platform:  g.name.string(),
				CreatedAt: time.Now(),
				Picture:   userInfo.Picture,
			},
		},
	}
	return
}

func newGoogleHandler(conf config.ConfigInfo) googleHandler {
	return googleHandler{
		oauth: oauth{
			conf: &oauth2.Config{
				ClientID:     conf.OAuth.Google.ClientID,
				ClientSecret: conf.OAuth.Google.ClientSecret,
				Endpoint:     conf.OAuth.Google.Endpoint,
				RedirectURL:  conf.OAuth.Google.RedirectURIS[0],
				Scopes:       conf.OAuth.Google.Scopes,
			},
			handler:     googleHandlerName,
			validStates: make(map[string]bool),
		},
		name: googleHandlerName,
	}

}
