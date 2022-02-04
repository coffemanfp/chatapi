package database

type RepositoryID string

type Database struct {
	Conn         DatabaseConnector
	Repositories map[RepositoryID]interface{}
}

type DatabaseConnector interface {
	Connect() error
}
