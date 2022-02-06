package config

import (
	"fmt"
	"os"
	"strconv"
)

type EnvManagerConfig struct {
	config ConfigInfo
}

func (f EnvManagerConfig) Get() ConfigInfo {
	return f.config
}

func NewEnvManagerConfig() (conf ConfigInfo, err error) {
	conf, err = newConfigWithEnvVars()
	return
}

func newConfigWithEnvVars() (conf ConfigInfo, err error) {
	srvPort, err := getEnvInt("PORT")
	if err != nil {
		return
	}

	dbPort, err := getEnvInt("DB_PORT")
	if err != nil {
		return
	}

	conf = ConfigInfo{
		Server: server{
			Port: srvPort,
			Host: os.Getenv("SRV_HOST"),
		},
		PostgreSQLProperties: postgreSQLProperties{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			Name:     os.Getenv("DB_NAME"),
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
		},
	}
	return
}

func getEnvInt(n string) (i int, err error) {
	i, err = strconv.Atoi(os.Getenv(n))
	if err != nil {
		err = fmt.Errorf("failed to load env var int %s: %s", n, err)
	}
	return
}
