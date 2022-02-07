package main

import (
	"fmt"
	"log"

	"github.com/coffemanfp/chat/config"
	"github.com/coffemanfp/chat/database"
	"github.com/coffemanfp/chat/database/psql"
	"github.com/coffemanfp/chat/server"
)

func main() {
	fmt.Println("Starting...")

	conf, err := config.NewEnvManagerConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := setUpDatabase(conf)
	if err != nil {
		log.Fatal(err)
	}

	server, err := server.NewServer(conf, db, conf.Server.Host, conf.Server.Port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening on port: %d\n", conf.Server.Port)
	log.Fatal(server.Run())
}

func setUpDatabase(conf config.ConfigInfo) (db database.Database, err error) {
	db.Conn = psql.NewPostgreSQLConnector(
		conf.PostgreSQLProperties.User,
		conf.PostgreSQLProperties.Password,
		conf.PostgreSQLProperties.Name,
		conf.PostgreSQLProperties.Host,
		conf.PostgreSQLProperties.Port,
	)

	err = db.Conn.Connect()
	if err != nil {
		log.Fatal(err)
	}

	authRepo, err := psql.NewAuthRepository(db.Conn.(*psql.PostgreSQLConnector))
	if err != nil {
		return
	}

	db.Repositories = map[database.RepositoryID]interface{}{
		database.AUTH_REPOSITORY: authRepo,
	}
	return
}
