package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
}

type ServerConfig struct {
	Mode                 string
	Name                 string
	Listener             string `yaml:"listener"`
	ServerTimeoutSeconds int    `yaml:"server_timeout_seconds"`
	Logfile              string `yaml:"logfile"`
}

type DBConfig struct {
	DBName   string `yaml:"db_name"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func MustLoadConfig(path string) *Config {
	f, err := os.ReadFile("setting/setting.yaml")
	if err != nil {
		panic(err)
	}

	cfg := new(Config)
	err = yaml.Unmarshal(f, cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
