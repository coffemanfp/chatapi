package psql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/coffemanfp/chat/account"
	"github.com/coffemanfp/chat/database"
)

// AuthRepository is the implementation of a authentication repository for the PostgreSQL database.
type AuthRepository struct {
	db *sql.DB
}

// NewAuthRepository initializes a new auth repository instance.
//
//	@param conn *PostgreSQLConnector: is the PostgreSQLConnector handler.
//	@return authRepo database.AuthRepository: is the final interface to keep
//	 the AuthRepository implementation.
//	@return err error: database connection error.
func NewAuthRepository(conn *PostgreSQLConnector) (repo database.AuthRepository, err error) {
	db, err := conn.getConn()
	if err != nil {
		return
	}
	repo = AuthRepository{
		db: db,
	}
	return
}

func (u AuthRepository) MatchCredentials(account account.Account) (id int, err error) {
	query := `
		select id from account where (nickname = $1 and password = $3) or (email = $2 and password = $3)
	`

	err = u.db.QueryRow(query, account.Nickname, account.Email, account.Password).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			return
		}
		err = fmt.Errorf("failed to get account credentials: %s", err)
	}
	return
}
