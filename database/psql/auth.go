package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coffemanfp/chat/auth"
	sErrors "github.com/coffemanfp/chat/errors"
	"github.com/coffemanfp/chat/users"
	"github.com/lib/pq"
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
	if pErr, ok := err.(*pq.Error); ok {
		switch pErr.Code {
		case foreign_key_violation:
			field := pErr.Detail[strings.Index(pErr.Detail, "(")+1 : strings.Index(pErr.Detail, ")")]
			var v interface{}
			switch field {
			case "nickname":
				v = user.Nickname
			case "email":
				v = user.Email
			}
			err = sErrors.NewClientError(http.StatusConflict, "already exists %s: %v", field, v)
		default:
			err = sErrors.NewClientError(http.StatusInternalServerError, "failed to insert user %s: %s", user.Nickname, err)
		}
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
		if pErr, ok := err.(*pq.Error); ok {
			switch pErr.Code {
			case foreign_key_violation:
				err = sErrors.NewClientError(http.StatusConflict, "already exists id: %s", sign.ID)
			default:
				err = sErrors.NewClientError(http.StatusInternalServerError, "failed to insert external user auth %s: %s", user.Nickname, err)
			}
			return
		}
	}
	return
}

func (u AuthRepository) UpsertSession(session auth.Session) (err error) {
	qInsertSession := `
		insert into
			user_session(id, tmp_id, user_id, logged_at, last_seen_at, logged_with, actived)
		values
			($1, $2, $3, $4, $5, $6, $7)
		on conflict (user_id, actived) do update set
			last_seen_at=$5, tmp_id = null
		where
			user_session.user_id=$3 and user_session.actived;
	`

	_, err = u.db.Exec(qInsertSession, session.ID, session.TmpID, session.UserID, session.LoggedAt, session.LastSeenAt, session.LoggedWith, session.Actived)
	if err != nil {
		err = fmt.Errorf("failed to insert session of %d: %s", session.UserID, err)
	}
	return
}

func (u AuthRepository) MatchCredentials(user users.User) (id int, err error) {
	query := `
		select id from users where (nickname = $1 and password = $3) or (email = $2 and password = $3)
	`

	err = u.db.QueryRow(query, user.Nickname, user.Email, user.Password).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			return
		}
		err = fmt.Errorf("failed to get user credentials: %s", err)
	}
	return
}
