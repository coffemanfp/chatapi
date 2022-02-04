package psql

import (
	"database/sql"
	"fmt"

	"github.com/coffemanfp/chat/models"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUserRepository(conn *PostgreSQLConnector) (usersRepo UsersRepository, err error) {
	db, err := conn.GetConn()
	if err != nil {
		return
	}
	usersRepo = UsersRepository{
		db: db,
	}
	return
}

func (u UsersRepository) SignUp(user models.User) (id int, err error) {
	query := `
		insert into	
			users(nickname, password, created_at)
		values
			($1, $2, $3)
		returning
			id
	`
	err = u.db.QueryRow(query, user.Nickname, user.Password, user.CreatedAt).Scan(&id)
	if err != nil {
		err = fmt.Errorf("failed to insert user %s: %s", user.Nickname, err)
	}
	return
}
