package auth

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type oauth struct {
	conf        *oauth2.Config
	handler     string
	validStates map[string]bool
}

func (o oauth) redirectToHandler(w http.ResponseWriter, r *http.Request) (err error) {
	authURL, err := url.Parse(o.conf.Endpoint.AuthURL)
	if err != nil {
		err = fmt.Errorf("failed to parse auth url: %s", err)
		return
	}

	state := uuid.NewString()
	o.validStates[state] = true

	parameters := url.Values{}

	parameters.Add("client_id", o.conf.ClientID)
	parameters.Add("scope", strings.Join(o.conf.Scopes, " "))
	parameters.Add("redirect_uri", o.conf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", state)

	authURL.RawQuery = parameters.Encode()
	http.Redirect(w, r, authURL.String(), http.StatusTemporaryRedirect)
	return
}

func (o oauth) callback(w http.ResponseWriter, r *http.Request) (response []byte, err error) {
	state := r.FormValue("state")
	_, ok := o.validStates[state]
	if !ok {
		err = fmt.Errorf("invalid callback state: %s is not a valid response callback state", state)
		return
	}
	delete(o.validStates, state)

	code := r.FormValue("code")

	if code == "" {
		w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission.."))
		}
		return
	}
	token, err := o.conf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("oauth in %s exchange failed with: %s\n", o.handler, err)
		err = errors.New("failed oauth callback error")
		return
	}

	resp, err := http.Get(genURLToRequestUserInfo(token.AccessToken, o.handler))
	if err != nil {
		err = fmt.Errorf("failed to request user info from %s: %s", o.handler, err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("failed to read user info from %s: %s", o.handler, err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	return
}

func genURLToRequestUserInfo(accessToken, handler string) (s string) {
	switch handler {
	case "google":
		s = fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", url.QueryEscape(accessToken))
	case "facebook":
		s = fmt.Sprintf("https://graph.facebook.com/me?access_token=%s&fields=email,picture,name", url.QueryEscape(accessToken))
	}
	return
}
