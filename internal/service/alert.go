package service

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/nktinn/TGAlerter/configs"
	"github.com/nktinn/TGAlerter/internal/model"
	"github.com/nktinn/TGAlerter/internal/repository"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
	"time"
)

type Alert struct {
	repo        *repository.Repository
	restyClient *resty.Client
	telegramBot *telebot.Bot
	telegramCfg configs.Telegram
	serviceCfg  []configs.Service
	logger      *zerolog.Logger
}

func NewAlert(repo *repository.Repository, telegramBot *telebot.Bot, telegramCfg configs.Telegram,
	serviceCfg []configs.Service, logger *zerolog.Logger) *Alert {
	return &Alert{
		repo:        repo,
		telegramBot: telegramBot,
		telegramCfg: telegramCfg,
		serviceCfg:  serviceCfg,
		logger:      logger,
		restyClient: resty.New(),
	}
}

func (a *Alert) SendAlert(alert model.Alert) error {
	var msg string
	switch alert.AlertType {
	case 1:
		msg = fmt.Sprintf("[EASY ERROR]\n%s\n%s", alert.ServiceID, alert.Message)
	case 2:
		msg = fmt.Sprintf("[MEDIUM ERROR]\n%s\n%s", alert.ServiceID, alert.Message)
	case 3:
		msg = fmt.Sprintf("[HARD ERROR]\n%s\n%s", alert.ServiceID, alert.Message)
	default:
		msg = fmt.Sprintf("[nothing]\n%s\n%s", alert.ServiceID, alert.Message)
	}

	for _, service := range a.serviceCfg {
		if service.ServiceID == alert.ServiceID {
			_, err := a.telegramBot.Send(&telebot.Chat{ID: service.UserID}, msg)
			if err != nil {
				a.logger.Error().Msgf("[%s] Error while sending alert: %s -- %s -- %s",
					time.Now().Format("20060102150405"), err.Error(), alert.ServiceID, alert.Message)
				return err
			}
			return nil
		}
	}
	_, err := a.telegramBot.Send(&telebot.Chat{ID: a.telegramCfg.AdminID}, msg)
	if err != nil {
		a.logger.Error().Msgf("[%s] Error while sending alert: %s -- %s -- %s",
			time.Now().Format("20060102150405"), err.Error(), alert.ServiceID, alert.Message)
		return err
	}
	return nil
}

func (a *Alert) HealthCheckWorker(healthCfg []configs.Health) {
	for _, url := range healthCfg {
		go func(u string) {
			ticker := time.NewTicker(10 * time.Minute)
			defer ticker.Stop()
			for {
				if err := a.HealthCheck(u); err != nil {
					a.logger.Error().Msgf("[%s] Error while health checking \"%s\": %s", time.Now().Format("20060102150405"), u, err.Error())
				}
				<-ticker.C
			}
		}(url.URL)
	}
}

func (a *Alert) HealthCheck(url string) error {
	resp, err := a.restyClient.R().Get(url)
	if err != nil {
		a.logger.Error().Msgf("[%s] Error while health checking \"%s\": %s", time.Now().Format("20060102150405"), url, err.Error())
		return err
	}
	if resp.StatusCode() != fiber.StatusOK {
		msg := fmt.Sprintf("[%s]\nHealth check \"%s\" error with status code: %d", time.Now().Format("20060102150405"), url, resp.StatusCode())
		a.logger.Error().Msgf("[%s] Health check \"%s\" status: %d", time.Now().Format("20060102150405"), url, resp.StatusCode())
		_, err := a.telegramBot.Send(&telebot.Chat{ID: a.telegramCfg.AdminID}, msg)
		if err != nil {
			a.logger.Error().Msgf("[%s] Error sending health check \"%s\" result: %s", time.Now().Format("20060102150405"), url, err.Error())
			return err
		}
		a.logger.Info().Msgf("[%s] Health check \"%s\" result sent", time.Now().Format("20060102150405"), url)
		return nil
	} else {
		a.logger.Info().Msgf("[%s] Health check \"%s\" status 200", time.Now().Format("20060102150405"), url)
	}
	return err
}
