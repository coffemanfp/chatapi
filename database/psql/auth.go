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

func NewAuthRepository(conn *PostgreSQLConnector) (usersRepo AuthRepository, err error) {
	db, err := conn.GetConn()
	if err != nil {
		return
	}
	usersRepo = AuthRepository{
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
			users(nickname, password, created_at)
		values
			($1, $2, $3)
		returning
			id
	`

	err = tx.QueryRow(qInsertUser, user.Nickname, user.Password, user.CreatedAt).Scan(&id)
	if err != nil {
		err = fmt.Errorf("failed to insert user %s: %s", user.Nickname, err)
		return
	}

	qInsertSession := `
		insert into
			user_session(id, user_id, logged_at, last_seen_at)
		values
			($1, $2, $3, $4)
	`

	_, err = tx.Exec(qInsertSession, session.ID, id, session.LoggedAt, session.LastSeenAt)
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
