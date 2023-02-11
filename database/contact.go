package database

import (
	"fmt"

	"github.com/coffemanfp/chat/contact"
)

// CONTACT_REPOSITORY is the key to be used when creating the repositories hashmap.
const CONTACT_REPOSITORY RepositoryID = "CONTACT"

// GetContactRepository gets the ContactRepository instance inside the repositories hashmap.
//
//	@param repoMap map[RepositoryID]interface{}: repositories hashmap.
//	@return repo ContactRepository: found ContactRepository instance.
//	@return err error: missing or invalid repository instance error.
func GetContactRepository(repoMap map[RepositoryID]interface{}) (repo ContactRepository, err error) {
	repoI, ok := repoMap[CONTACT_REPOSITORY]
	if !ok {
		err = fmt.Errorf("missing repository: %s not found in repository map", CONTACT_REPOSITORY)
		return
	}
	repo, ok = repoI.(ContactRepository)
	if !ok {
		err = fmt.Errorf("invalid repository value: %s has a invalid %s repository handler", CONTACT_REPOSITORY, CONTACT_REPOSITORY)
	}
	return
}

// ContactRepository defines the behaviors to be used by a ContactRepository implementation.
type ContactRepository interface {
	GetByRange(id, limit, offset int) ([]contact.Contact, error)
}
