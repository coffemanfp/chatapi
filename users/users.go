// Package users handles all the user logic, like creation and validation.
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

// User is the representation of the common user data
type User struct {
	ID         int              `json:"id"`
	Nickname   string           `json:"nickname"`
	Email      string           `json:"email"`
	Password   string           `json:"password,omitempty"`
	Picture    string           `json:"picture,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
	SignedWith []ExternalSigned `json:"signed_with,omitempty"`
}

// ExternalSigned represents the data required for external sign in services models.
type ExternalSigned struct {
	ID        string    `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Picture   string    `json:"picture,omitempty"`
	Platform  string    `json:"platform,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// New initializes a new user based on the basic data provided from the user passed as param.
// 	@param userR User: Basic data of the user to build.
// 	@return user User: User builded
// 	@return err error: error in the validation of the based user.
func New(userR User) (user User, err error) {
	// If the user is not registered with an external platform, validate the nickname and password.
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

// HashPassword encrypt the password provided with a bcrypt algorithm.
// 	@param orig *string: Is the password to encrypt.
//	  The result of the bcrypt is assigned to this same param.
//  @return err error: bcrypt encriptation error.
func HashPassword(orig *string) (err error) {
	h, err := auth.HashPassword(*orig)
	if err != nil {
		return
	}
	orig = &h
	return
}

var nicknameRegex = regexp.MustCompile(`^[^0-9]\w+$`)

// ValidateNickname validate the nickname with a regular expression.
// 	@param nickname string: nickname to validate.
//  @return err error: don't match the regex with the string provided.
func ValidateNickname(nickname string) (err error) {
	if !nicknameRegex.MatchString(nickname) {
		err = errors.NewClientError(http.StatusBadRequest, "invalid nickname: invalid nickname format of %s", nickname)
	}
	return
}

// ValidateEmail validate the email with a standart library and
//	check the host.
// @param email string: email to validate.
// @return err error: invalid format of the email or the host.
func ValidateEmail(email string) (err error) {
	// Check email format
	_, err = mail.ParseAddress(email)
	if err != nil {
		err = errors.NewClientError(http.StatusBadRequest, "invalid email format: %s is not valid, cause %s", email, err)
		return
	}

	// Check the host
	host := strings.Split(email, "@")[1]
	_, err = net.LookupHost(host)
	if err != nil {
		err = errors.NewClientError(http.StatusBadRequest, "invalid email host: %s not exists", host)
	}
	return
}
