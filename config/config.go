package config

import (
	"fmt"
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Server               server               `yaml:"server"`
	PostgreSQLProperties postgreSQLProperties `yaml:"psql"`
}

type server struct {
	Port int
	Host string
}

type postgreSQLProperties struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

func NewConfig(env, configDir string) (conf Config, err error) {
	p := fmt.Sprintf("%s.yaml", env)

	fileBytes, err := ioutil.ReadFile(path.Join(configDir, p))
	if err != nil {
		err = fmt.Errorf("not found: config filepath %s not found", p)
		return
	}

	err = yaml.Unmarshal(fileBytes, &conf)
	if err != nil {
		err = fmt.Errorf("invalid config: failed to get config info. Bad structure?")
	}
	return
}
