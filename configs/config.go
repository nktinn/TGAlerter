package configs

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server   Server    `yaml:"server"`
	Telegram Telegram  `yaml:"telegram"`
	Service  []Service `yaml:"service"`
	Health   []Health  `yaml:"healthCheck"`
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
