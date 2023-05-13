package database

import (
	"github.com/coffemanfp/chat/account"
)

// AUTH_REPOSITORY is the key to be used when creating the repositories hashmap.
const AUTH_REPOSITORY RepositoryID = "AUTH"

// GetAuthRepository gets the AuthRepository instance inside the repositories hashmap.
//
//		@param repoMap map[RepositoryID]interface{}: repositories hashmap.
//		@return repo AuthRepository: found AuthRepository instance.
//	 @return err error: missing or invalid repository instance error.
func GetAuthRepository(repoMap map[RepositoryID]interface{}) (repo AuthRepository, err error) {
	return GetRepository(repoMap, AUTH_REPOSITORY)
}

// AuthRepository defines the behaviors to be used by a AuthRepository implementation.
type AuthRepository interface {

	// MatchCredentials locate a account by its credentials which check if match.
	//  Returns the id of the account if match.
	//	@param account account.Account: account to check the credentials.
	//	@return $1 int: id of the matched account. Is 0 if it don't match.
	//	@return $2 error: failed credentials validation process.
	MatchCredentials(account account.Account) (int, error)
}
