package database

import (
	"fmt"
)

// RepositoryID is the key to use for the repositories hashmap.
type RepositoryID string

// Database is the Database manager for connections and repository instancies.
type Database struct {
	Conn         DatabaseConnector
	Repositories map[RepositoryID]interface{}
}

// DatabaseConnector defines a database connector handler.
type DatabaseConnector interface {

	// Connect creates new connection of the database implementation.
	//  @return $1 error: database connection error
	Connect() error
}

func GetRepository(repoMap map[RepositoryID]interface{}, id RepositoryID) (repo AuthRepository, err error) {
	repoI, ok := repoMap[id]
	if !ok {
		err = fmt.Errorf("missing repository: %s not found in repository map", id)
		return
	}
	repo, ok = repoI.(AuthRepository)
	if !ok {
		err = fmt.Errorf("invalid repository value: %s has a invalid %s repository handler", id, id)
	}
	return
}
