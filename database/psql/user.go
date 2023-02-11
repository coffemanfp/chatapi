package psql

import (
	"database/sql"
)

// AccountRepository is the implementation of a account repository for the PostgreSQL database.
type AccountRepository struct {
	db *sql.DB
}

// NewAuthRepository initializes a new AccountRepository instance.
//
//	@param conn *PostgreSQLConnector: is the PostgreSQLConnector handler.
//	@return accountRepo AccountRepository: new AccountRepository instance.
//	@return err error: database connection error.
func NewAccountRepository(conn *PostgreSQLConnector) (accountRepo AccountRepository, err error) {
	db, err := conn.getConn()
	if err != nil {
		return
	}
	accountRepo = AccountRepository{
		db: db,
	}
	return
}
