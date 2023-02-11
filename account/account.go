package account

import (
	"net/http"
	"regexp"
	"time"
	"unicode"

	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/errors"
)

// Account is the representation of the common account data
type Account struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Nickname   string `json:"nickname,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	PictureURL string `json:"picture_url,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
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
	if err != nil {
		return
	}
	account.CreatedAt = &time.Time{}
	*account.CreatedAt = time.Now()
	return
}

var nicknameRegex = regexp.MustCompile(`^[a-z0-9_-]{3,16}$`)

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

var nameRegex = regexp.MustCompile(`^([A-Za-zÑñÁáÉéÍíÓóÚú]+['-]{0,1}[A-Za-zÑñÁáÉéÍíÓóÚú]+)(\s+([A-Za-zÑñÁáÉéÍíÓóÚú]+['-]{0,1}[A-Za-zÑñÁáÉéÍíÓóÚú]+))*$`)

// ValidateName validate the name with a regular expression.
//
//	@param name string: name to validate.
//	 @return err error: don't match the regex with the string provided.
func ValidateName(name string) (err error) {
	if !nameRegex.MatchString(name) {
		err = errors.NewClientError(http.StatusBadRequest, "invalid name: invalid name format of %s", name)
	}
	return
}

var lastNameRegex = regexp.MustCompile(`^([A-Za-zÑñÁáÉéÍíÓóÚú]+['-]{0,1}[A-Za-zÑñÁáÉéÍíÓóÚú]+)(\s+([A-Za-zÑñÁáÉéÍíÓóÚú]+['-]{0,1}[A-Za-zÑñÁáÉéÍíÓóÚú]+))*$`)

// ValidateLastName validate the lastName with a regular expression.
//
//	@param lastName string: lastName to validate.
//	 @return err error: don't match the regex with the string provided.
func ValidateLastName(lastName string) (err error) {
	if !lastNameRegex.MatchString(lastName) {
		err = errors.NewClientError(http.StatusBadRequest, "invalid lastName: invalid lastName format of %s", lastName)
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

// ValidatePassword validate the password.
//
//	@param password string: password to validate.
//	 @return err error: don't match the required charactes with the string provided.
func ValidatePassword(password string) (err error) {
	letters := 0
	var number, upper, special, sevenOrMore bool
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		}
	}
	sevenOrMore = letters >= 7
	if !(number && upper && special && sevenOrMore) {
		err = errors.NewClientError(http.StatusBadRequest, "invalid password: invalid password format of %s", password)
	}
	return
}

// // ValidateEmail validate the email with a standart library and
// //
// //	check the host.
// //
// // @param email string: email to validate.
// // @return err error: invalid format of the email or the host.
// func ValidateEmail(email string) (err error) {
// 	// Check email format
// 	_, err = mail.ParseAddress(email)
// 	if err != nil {
// 		err = errors.NewClientError(http.StatusBadRequest, "invalid email format: %s is not valid, cause %s", email, err)
// 		return
// 	}

// 	// Check the host
// 	host := strings.Split(email, "@")[1]
// 	_, err = net.LookupHost(host)
// 	if err != nil {
// 		err = errors.NewClientError(http.StatusBadRequest, "invalid email host: %s not exists", host)
// 	}
// 	return
// }
