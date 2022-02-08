package psql

import (
	"database/sql"
)

type UsersRepository struct {
	db *sql.DB
}

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
