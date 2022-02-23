package users

import (
	"net"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/errors"
)

type User struct {
	ID         int              `json:"id"`
	Nickname   string           `json:"nickname"`
	Email      string           `json:"email"`
	Password   string           `json:"password,omitempty"`
	Picture    string           `json:"picture,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
	SignedWith []ExternalSigned `json:"signed_with,omitempty"`
}

type ExternalSigned struct {
	ID        string    `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Picture   string    `json:"picture,omitempty"`
	Platform  string    `json:"platform,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func New(userR User) (user User, err error) {
	if len(userR.SignedWith) == 0 {
		err = ValidateNickname(userR.Nickname)
		if err != nil {
			return
		}
		err = HashPassword(&user.Password)
		if err != nil {
			return
		}
	}
	user = userR
	user.CreatedAt = time.Now()
	return
}

func HashPassword(orig *string) (err error) {
	h, err := auth.HashPassword(*orig)
	if err != nil {
		return
	}
	orig = &h
	return
}

var nicknameRegex = regexp.MustCompile(`^[^0-9]\w+$`)

func ValidateNickname(nickname string) (err error) {
	if !nicknameRegex.MatchString(nickname) {
		err = errors.NewClientError(http.StatusBadRequest, "invalid nickname: invalid nickname format of %s", nickname)
	}
	return
}

func ValidateEmail(email string) (err error) {
	_, err = mail.ParseAddress(email)
	if err != nil {
		err = errors.NewClientError(http.StatusBadRequest, "invalid email format: %s is not valid, cause %s", email, err)
		return
	}

	parts := strings.Split(email, "@")
	_, err = net.LookupHost(parts[1])
	if err != nil {
		err = errors.NewClientError(http.StatusBadRequest, "invalid email host: %s not exists", parts[1])
	}
	return
}
