package database

import (
	"fmt"

	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/users"
)

const AUTH_REPOSITORY RepositoryID = "AUTH"

func GetAuthRepository(repoMap map[RepositoryID]interface{}) (repo AuthRepository, err error) {
	repoI, ok := repoMap[AUTH_REPOSITORY]
	if !ok {
		err = fmt.Errorf("missing repository: %s not found in repository map", AUTH_REPOSITORY)
		return
	}
	repo, ok = repoI.(AuthRepository)
	if !ok {
		err = fmt.Errorf("invalid repository value: %s has a invalid %s repository handler", AUTH_REPOSITORY, AUTH_REPOSITORY)
	}
	return
}

type AuthRepository interface {
	SignUp(user users.User, session auth.Session) (int, error)
	FindPassword(nickname string) (string, error)
}
