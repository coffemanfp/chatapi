package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/coffemanfp/chat/account"
	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/database"
	sErrors "github.com/coffemanfp/chat/errors"
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

func (a AuthRepository) SignUp(account account.Account, session auth.Session) (id int, err error) {
	tx, err := a.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	qInsertAccount := `
		insert into	
			account(name, last_name, nickname, email, password, picture_url, created_at)
		values
			($1, nullif($2, ''), $3, $4, $5, $6, $7)
		returning
			id
	`

	err = tx.QueryRow(qInsertAccount, account.Name, account.LastName, account.Nickname, account.Email, account.Password, account.PictureURL, account.CreatedAt).Scan(&id)
	if err != nil {
		var match bool
		match, err = newPQError(err).asAlreadyExists()
		if match {
			err = sErrors.NewClientError(http.StatusConflict, err.Error())
		} else {
			err = fmt.Errorf("failed to insert account %s: %s", account.Nickname, err)
		}
	}
	return
}

func (u AuthRepository) UpsertSession(session auth.Session) (err error) {
	qInsertSession := `
		insert into
			account_session(id, account_id, logged_at, last_seen_at, logged_with, actived)
		values
			($1, $2, $3, $4, $5, $6)
		on conflict (account_id, actived) do update set
			last_seen_at=$4
		where
			account_session.account_id=$2 and account_session.actived;
	`

	_, err = u.db.Exec(qInsertSession, session.ID, session.AccountID, session.LoggedAt, session.LastSeenAt, session.LoggedWith, session.Actived)
	if err != nil {
		err = fmt.Errorf("failed to insert session of %d: %s", session.AccountID, err)
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
