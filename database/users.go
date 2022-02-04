package database

import (
	"fmt"

	"github.com/coffemanfp/chat/models"
)

const USERS_REPOSITORY RepositoryID = "USERS"

func GetUsersRepository(repoMap map[RepositoryID]interface{}) (repo UsersRepository, err error) {
	repoI, ok := repoMap[USERS_REPOSITORY]
	if !ok {
		err = fmt.Errorf("missing repository: %s not found in repository map", USERS_REPOSITORY)
		return
	}
	repo, ok = repoI.(UsersRepository)
	if !ok {
		err = fmt.Errorf("invalid repository value: %s has a invalid %s repository handler", USERS_REPOSITORY, USERS_REPOSITORY)
	}
	return
}

type UsersRepository interface {
	SignUp(user models.User) (int, error)
}
