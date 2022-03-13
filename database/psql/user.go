package psql

import (
	"database/sql"
)

// UsersRepository is the implementation of a user repository for the PostgreSQL database.
type UsersRepository struct {
	db *sql.DB
}

// NewAuthRepository initializes a new UsersRepository instance.
// 	@param conn *PostgreSQLConnector: is the PostgreSQLConnector handler.
//	@return usersRepo UsersRepository: new UsersRepository instance.
//	@return err error: database connection error.
func NewUsersRepository(conn *PostgreSQLConnector) (usersRepo UsersRepository, err error) {
	db, err := conn.GetConn()
	if err != nil {
		return
	}
	usersRepo = UsersRepository{
		db: db,
	}
	return
}
