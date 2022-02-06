package users

import (
	"fmt"
	"regexp"
	"time"

	"github.com/coffemanfp/chat/auth"
)

type User struct {
	ID        int       `json:"id"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func New(userR User) (user User, err error) {
	err = ValidateNickname(userR.Nickname)
	if err != nil {
		return
	}
	ph, err := auth.HashPassword(userR.Password)
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

func ValidateNickname(nickname string) (err error) {
	if !nicknameRegex.MatchString(nickname) {
		err = fmt.Errorf("invalid nickname: invalid nickname format of %s", nickname)
	}
	return
}
