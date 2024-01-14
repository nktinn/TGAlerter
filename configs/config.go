package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   Server    `yaml:"server"`
	Telegram Telegram  `yaml:"telegram"`
	Service  []Service `yaml:"service"`
	Health   []Health  `yaml:"healthCheck"`
	Database Database  `yaml:"postgres"`
}

type Server struct {
	Port string `yaml:"port"`
}

type Telegram struct {
	Token   string `yaml:"telegramToken"`
	AdminID int64  `yaml:"adminId"`
}

type Service struct {
	ServiceID string `yaml:"serviceId"`
	UserID    int64  `yaml:"userId"`
}

type Health struct {
	URL string `yaml:"url"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbName"`
	SSLMode  string `yaml:"sslMode"`
}

func NewConfig() *Config {
	data, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return nil
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil
	}
	return &cfg
}
