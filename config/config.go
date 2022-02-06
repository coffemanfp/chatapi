package config

type Config interface {
	Get() ConfigInfo
}

type ConfigInfo struct {
	Server               server               `yaml:"server"`
	PostgreSQLProperties postgreSQLProperties `yaml:"psql"`
}

type server struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type postgreSQLProperties struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}
