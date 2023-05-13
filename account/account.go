package account

import (
	"net/http"
	"regexp"

	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/errors"
)

// Account is the representation of the common account data
type Account struct {
	ID       int    `json:"id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// New initializes a new account based on the basic data provided from the account passed as param.
//
//	@param accountR Account: Basic data of the account to build.
//	@return account Account: Account builded
//	@return err error: error in the validation of the based account.
func New(accountR Account) (account Account, err error) {
	err = ValidateNickname(accountR.Nickname)
	if err != nil {
		return
	}
	account = accountR
	account.Password, err = auth.HashPassword(account.Password)
	return
}

var nicknameRegex = regexp.MustCompile(`^[a-z0-9_-]{3,32}$`)

// ValidateNickname validate the nickname with a regular expression.
//
//	@param nickname string: nickname to validate.
//	 @return err error: don't match the regex with the string provided.
func ValidateNickname(nickname string) (err error) {
	if !nicknameRegex.MatchString(nickname) {
		err = errors.NewClientError(http.StatusBadRequest, "invalid nickname: invalid nickname format of %s", nickname)
	}
	return
}

var emailRegex = regexp.MustCompile(`^([A-Za-zÑñÁáÉéÍíÓóÚú]+['-]{0,1}[A-Za-zÑñÁáÉéÍíÓóÚú]+)(\s+([A-Za-zÑñÁáÉéÍíÓóÚú]+['-]{0,1}[A-Za-zÑñÁáÉéÍíÓóÚú]+))*$`)

// ValidateEmail validate the email with a regular expression.
//
//	@param email string: email to validate.
//	 @return err error: don't match the regex with the string provided.
func ValidateEmail(email string) (err error) {
	if !emailRegex.MatchString(email) {
		err = errors.NewClientError(http.StatusBadRequest, "invalid email: invalid email format of %s", email)
	}
	return
}
