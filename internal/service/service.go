package service

import (
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
	
	"github.com/nktinn/TGAlerter/configs"
	"github.com/nktinn/TGAlerter/internal/model"
	"github.com/nktinn/TGAlerter/internal/repository"
)

type Alerter interface {
	GetRoute(route string) int64
	SendAlert(alert model.Alert) error
	HealthCheck(url string) error
	HealthCheckWorker(healthCfg []configs.Health)
}

type Service struct {
	Alerter
}

func NewService(
	repo *repository.Repository,
	telegramBot *telebot.Bot,
	telegramCfg configs.Telegram,
	serviceCfg []configs.Service,
	logger *zerolog.Logger,
) *Service {
	return &Service{
		Alerter: NewAlert(repo, telegramBot, telegramCfg, serviceCfg, logger),
	}
}
