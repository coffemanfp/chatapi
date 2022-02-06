package auth

import "fmt"

func Login(nickname, pwHash, pwTry string) (Session, error) {
	sL := SystemLogger{
		nickname: nickname,
		pwHash:   pwHash,
		pwTry:    pwTry,
	}
	return sL.Login()
}

type Logger interface {
	Login() (Session, error)
}

type SystemLogger struct {
	nickname string
	pwHash   string
	pwTry    string
}

func (s SystemLogger) Login() (session Session, err error) {
	if !checkPasswordHash(s.pwTry, s.pwHash) {
		err = fmt.Errorf("invalid credentials: password don't match")
		return
	}

	session, err = NewSession(s.nickname)
	return
}
