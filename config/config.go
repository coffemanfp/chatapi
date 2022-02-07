package config

import "golang.org/x/oauth2"

type Config interface {
	Get() ConfigInfo
}

type ConfigInfo struct {
	Server               server               `yaml:"server"`
	OAuth                oauth                `yaml:"oauth"`
	PostgreSQLProperties postgreSQLProperties `yaml:"psql"`
}

type server struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type oauth struct {
	Google   oauthProperties `yaml:"google"`
	Facebook oauthProperties `yaml:"facebook"`
}

type oauthProperties struct {
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	RedirectURIS []string `yaml:"redirect_uris"`
	Scopes       []string `yaml:"scopes"`
	State        string   `yaml:"state"`
	Endpoint     oauth2.Endpoint
}

type postgreSQLProperties struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}
