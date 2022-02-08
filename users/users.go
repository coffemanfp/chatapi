package users

import (
	"fmt"
	"net"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/coffemanfp/chat/auth"
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
	var ph string
	if len(userR.SignedWith) == 0 {
		err = ValidateNickname(userR.Nickname)
		if err != nil {
			return
		}
		ph, err = auth.HashPassword(userR.Password)
		if err != nil {
			return
		}
	}
	user = userR
	user.Password = ph
	user.CreatedAt = time.Now()
	return
}

var nicknameRegex = regexp.MustCompile(`^[^0-9]\w+$`)

func ValidateNickname(nickname string) (err error) {
	if !nicknameRegex.MatchString(nickname) {
		err = fmt.Errorf("invalid nickname: invalid nickname format of %s", nickname)
	}
	return
}

func ValidateEmail(email string) (err error) {
	_, err = mail.ParseAddress(email)
	if err != nil {
		err = fmt.Errorf("invalid email format: %s is not valid, cause %s", email, err)
		return
	}

	parts := strings.Split(email, "@")
	_, err = net.LookupHost(parts[1])
	if err != nil {
		err = fmt.Errorf("invalid email host: %s not exists", parts[1])
	}
	return
}
