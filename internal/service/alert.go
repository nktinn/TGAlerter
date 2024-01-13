package service

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/nktinn/TGAlerter/configs"
	"github.com/nktinn/TGAlerter/internal/model"
	"github.com/nktinn/TGAlerter/internal/repository"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
	"time"
)

type Alert struct {
	repo        *repository.Repository
	telegramBot *telebot.Bot
	telegramCfg configs.Telegram
	serviceCfg  []configs.Service
	logger      zerolog.Logger
}

func NewAlert(repo *repository.Repository, telegramBot *telebot.Bot, telegramCfg configs.Telegram, serviceCfg []configs.Service, logger zerolog.Logger) *Alert {
	return &Alert{
		repo:        repo,
		telegramBot: telegramBot,
		telegramCfg: telegramCfg,
		serviceCfg:  serviceCfg,
		logger:      logger,
	}
}

func (a *Alert) SendAlert(alert model.Alert) error {
	var msg string
	switch alert.AlertType {
	case 1:
		msg = fmt.Sprintf("[EASY ERROR]\n%s\n%s", alert.ServiceId, alert.Message)
	case 2:
		msg = fmt.Sprintf("[MEDIUM ERROR]\n%s\n%s", alert.ServiceId, alert.Message)
	case 3:
		msg = fmt.Sprintf("[HARD ERROR]\n%s\n%s", alert.ServiceId, alert.Message)
	default:
		msg = fmt.Sprintf("[nothing]\n%s\n%s", alert.ServiceId, alert.Message)
	}

	for _, service := range a.serviceCfg {
		if service.ServiceId == alert.ServiceId {
			_, err := a.telegramBot.Send(&telebot.Chat{ID: service.UserId}, msg)
			if err != nil {
				a.logger.Error().Msgf("["+time.Now().Format("20060102150405")+"] "+"Error while sending alert: %s -- %s -- %s", err.Error(), alert.ServiceId, alert.Message)
				return err
			}
			return nil
		}
	}
	_, err := a.telegramBot.Send(&telebot.Chat{ID: a.telegramCfg.AdminId}, msg)
	if err != nil {
		a.logger.Error().Msgf("["+time.Now().Format("20060102150405")+"] "+"Error while sending alert: %s -- %s -- %s", err.Error(), alert.ServiceId, alert.Message)
		return err
	}
	return nil
}

func (a *Alert) HealthCheck(client *resty.Client, url string) error {
	resp, err := client.R().Get(url)
	if err != nil {
		a.logger.Error().Msgf("["+time.Now().Format("20060102150405")+"] "+"Error while health checking \"%s\": %s", url, err.Error())
		return err
	}
	if resp.StatusCode() != 200 {
		msg := fmt.Sprintf("["+time.Now().Format("2006-01-02 15-04-05")+"]\n"+"Health check \"%s\" error with status code: %d", url, resp.StatusCode())
		a.logger.Error().Msgf("["+time.Now().Format("20060102150405")+"] "+"Health check \"%s\" status: %d", url, resp.StatusCode())
		_, err := a.telegramBot.Send(&telebot.Chat{ID: a.telegramCfg.AdminId}, msg)
		if err != nil {
			a.logger.Error().Msgf("["+time.Now().Format("20060102150405")+"] "+"Error sending health check \"%s\" result: %s", url, err.Error())
			return err
		}
		a.logger.Info().Msgf("["+time.Now().Format("20060102150405")+"] "+"Health check \"%s\" result sent", url)
		return nil
	} else {
		a.logger.Info().Msgf("["+time.Now().Format("20060102150405")+"] "+"Health check \"%s\" status 200", url)
	}
	return err
}
