package psql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/coffemanfp/chat/auth"
	"github.com/coffemanfp/chat/users"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(conn *PostgreSQLConnector) (authRepo AuthRepository, err error) {
	db, err := conn.GetConn()
	if err != nil {
		return
	}
	authRepo = AuthRepository{
		db: db,
	}
	return
}

func (u AuthRepository) SignUp(user users.User, session auth.Session) (id int, err error) {
	tx, err := u.db.Begin()
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

	qInsertUser := `
		insert into	
			users(nickname, email, password, picture, created_at)
		values
			(nullif($1, ''), nullif($2, ''), $3, $4, $5)
		returning
			id
	`

	err = tx.QueryRow(qInsertUser, user.Nickname, user.Email, user.Password, user.Picture, user.CreatedAt).Scan(&id)
	if err != nil {
		err = fmt.Errorf("failed to insert user %s: %s", user.Nickname, err)
		return
	}

	if len(user.SignedWith) > 0 {
		sign := user.SignedWith[0]

		qInsertExtUserAuth := `
			insert into
				external_user_auth(id, user_id, email, picture, platform, created_at)
			values
				($1, $2, $3, $4, $5, $6)
		`
		_, err = tx.Exec(qInsertExtUserAuth, sign.ID, id, sign.Email, sign.Picture, sign.Platform, sign.CreatedAt)
		if err != nil {
			err = fmt.Errorf("failed to insert external user auth %s: %s", user.Nickname, err)
			return
		}
	}

	qInsertSession := `
		insert into
			user_session(id, user_id, logged_at, last_seen_at, logged_with)
		values
			($1, $2, $3, $4, $5)
	`

	_, err = tx.Exec(qInsertSession, session.ID, id, session.LoggedAt, session.LastSeenAt, session.LoggedWith)
	if err != nil {
		err = fmt.Errorf("failed to insert session of %s: %s", user.Nickname, err)
	}
	return
}

func (u AuthRepository) FindPassword(nickname string) (pw string, err error) {
	query := `
		select password from users where nickname = $1
	`

	err = u.db.QueryRow(query, nickname).Scan(&pw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("not found: %s not exists", nickname)
			return
		}
		err = fmt.Errorf("failed to get user credentials: %s", err)
	}
	return
}
