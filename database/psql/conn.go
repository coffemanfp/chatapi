package psql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type properties struct {
	user string
	pass string
	name string
	host string
	port int
}

type PostgreSQLConnector struct {
	props properties
	db    *sql.DB
}

func (p *PostgreSQLConnector) Connect() (err error) {
	db, err := sql.Open("postgres", connURL(p.props))
	if err != nil {
		return
	}

	err = db.Ping()
	if err != nil {
		err = fmt.Errorf("failed to ping database: %s", err)
		return
	}
	p.db = db
	return
}

func (p PostgreSQLConnector) GetConn() (conn *sql.DB, err error) {
	err = p.db.Ping()
	if err != nil {
		err = fmt.Errorf("failed to ping database: %s", err)
		return
	}

	conn = p.db
	return
}

func NewPostgreSQLConnector(user, pass, name, host string, port int) (conn *PostgreSQLConnector) {
	return &PostgreSQLConnector{
		props: properties{
			user: user,
			pass: pass,
			name: name,
			host: host,
			port: port,
		},
	}
}

func connURL(props properties) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", props.user, props.pass, props.host, props.port, props.name, "disable")
}
