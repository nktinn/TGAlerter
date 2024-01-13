package service

import (
	"github.com/go-resty/resty/v2"
	"github.com/nktinn/TGAlerter/configs"
	"github.com/nktinn/TGAlerter/internal/model"
	"github.com/nktinn/TGAlerter/internal/repository"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

type Alerter interface {
	SendAlert(alert model.Alert) error
	HealthCheck(client *resty.Client, url string) error
}

type Service struct {
	Alerter
}

func NewService(repos *repository.Repository, telegramBot *telebot.Bot, telegramCfg configs.Telegram, serviceCfg []configs.Service, logger zerolog.Logger) *Service {
	return &Service{
		Alerter: NewAlert(repos, telegramBot, telegramCfg, serviceCfg, logger),
	}
}
