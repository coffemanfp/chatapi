package psql

import (
	"database/sql"
	"fmt"

	"github.com/coffemanfp/chat/contact"
	"github.com/coffemanfp/chat/database"
)

type ContactRepository struct {
	db *sql.DB
}

func NewContactRepository(conn *PostgreSQLConnector) (repo database.ContactRepository, err error) {
	db, err := conn.getConn()
	if err != nil {
		return
	}
	repo = ContactRepository{
		db: db,
	}
	return
}

func (c ContactRepository) GetByRange(id, limit, offset int) (contacts []contact.Contact, err error) {
	contacts = make([]contact.Contact, 0)
	qSelectContacts := `
		SELECT
			c.id, c.name, c.last_name, c.created_at,
			a.id, a.name, a.last_name, a.nickname, a.picture_url
		FROM
			contact c
		INNER JOIN
			account a ON c.to_account_id = a.id
		WHERE
			c.from_account_id = $1
		LIMIT
			$2
		OFFSET
			$3
	`

	rows, err := c.db.Query(qSelectContacts, id, limit, offset)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var c contact.Contact
		err = rows.Scan(&c.ID, &c.Name, &c.LastName, &c.CreatedAt,
			&c.Account.ID, &c.Account.Name, &c.Account.LastName, &c.Account.Nickname, &c.Account.PictureURL,
		)
		if err != nil {
			contacts = make([]contact.Contact, 0)
			err = fmt.Errorf("failed to get contacts of account %d", id)
			return
		}
		contacts = append(contacts, c)
	}
	err = rows.Err()
	if err != nil {
		contacts = make([]contact.Contact, 0)
	}
	return
}
