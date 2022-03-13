package database

import (
	"fmt"

	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/users"
)

// AUTH_REPOSITORY is the key to be used when creating the repositories hashmap.
const AUTH_REPOSITORY RepositoryID = "AUTH"

// GetAuthRepository gets the AuthRepository instance inside the repositories hashmap.
// 	@param repoMap map[RepositoryID]interface{}: repositories hashmap.
// 	@return repo AuthRepository: found AuthRepository instance.
//  @return err error: missing or invalid repository instance error.
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

// AuthRepository defines the behaviors to be used by a AuthRepository implementation.
type AuthRepository interface {
	// SignUp creates the records for a new user register and its session.
	//  Returns a new id for the user created.
	//  @param user users.User: user to be created.
	//  @param session auth.Session: session to be created.
	//  @return $1 int: new generated ID.
	//  @return $2 error: failed record creation.
	SignUp(user users.User, session auth.Session) (int, error)

	// MatchCredentials locate a user by its credentials which check if match.
	//  Returns the id of the user if match.
	//	@param user users.User: user to check the credentials.
	//	@return $1 int: id of the matched user. Is 0 if it don't match.
	//	@return $2 error: failed credentials validation process.
	MatchCredentials(user users.User) (int, error)

	// UpsertSession creates or updates the user session.
	//  @param session auth.Session: session to create or update.
	//  @return $1 error: failed record creation or update.
	UpsertSession(session auth.Session) error
}
