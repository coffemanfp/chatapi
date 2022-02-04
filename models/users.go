package models

import (
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(userR User) (user User, err error) {
	err = validateNickname(userR.Nickname)
	if err != nil {
		return
	}
	ph, err := hashPassword(userR.Password)
	if err != nil {
		return
	}
	user = User{
		Nickname:  userR.Nickname,
		Password:  ph,
		CreatedAt: time.Now(),
	}
	return
}

var nicknameRegex = regexp.MustCompile(`^[^0-9]\w+$`)

func validateNickname(nickname string) (err error) {
	if !nicknameRegex.MatchString(nickname) {
		err = fmt.Errorf("invalid nickname: invalid nickname format of %s", nickname)
	}
	return
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		err = fmt.Errorf("failed to generate password: %s", err)
	}
	return string(bytes), err
}
