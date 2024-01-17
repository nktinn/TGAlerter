package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"

	"github.com/nktinn/TGAlerter/configs"
	"github.com/nktinn/TGAlerter/internal/model"
	"github.com/nktinn/TGAlerter/internal/repository"
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

func (a *Alert) GetRoute(serviceID string) int64 {
	return a.repo.GetRoute(serviceID)
}

func (a *Alert) SendAlert(msgPost *nats.Msg) error {
	var msgToSend string

	var alert model.Alert
	if err := json.Unmarshal(msgPost.Data, &alert); err != nil {
		a.logger.Error().Msgf("[%s] Error while parsing json body: %s", time.Now().Format("20060102150405"), err.Error())
		return err
	}

	switch alert.AlertType {
	case 1:
		msgToSend = fmt.Sprintf("[EASY ERROR]\n%s\n%s", alert.ServiceID, alert.Message)
	case 2:
		msgToSend = fmt.Sprintf("[MEDIUM ERROR]\n%s\n%s", alert.ServiceID, alert.Message)
	case 3:
		msgToSend = fmt.Sprintf("[HARD ERROR]\n%s\n%s", alert.ServiceID, alert.Message)
	default:
		msgToSend = fmt.Sprintf("[nothing]\n%s\n%s", alert.ServiceID, alert.Message)
	}

	//userID := a.GetRoute(alert.ServiceID)
	userID := int64(0)
	if userID != 0 {
		_, err := a.telegramBot.Send(&telebot.Chat{ID: userID}, msgToSend)
		if err != nil {
			a.logger.Error().Msgf("[%s] Error while sending alert: %s -- %s -- %s",
				time.Now().Format("20060102150405"), err.Error(), alert.ServiceID, alert.Message)
			return err
		}
		return nil
	}
	_, err := a.telegramBot.Send(&telebot.Chat{ID: a.telegramCfg.AdminID}, msgToSend)
	if err != nil {
		a.logger.Error().Msgf("[%s] Error while sending alert: %s -- %s -- %s",
			time.Now().Format("20060102150405"), err.Error(), alert.ServiceID, alert.Message)
		return err
	}
	a.logger.Info().Msgf("[%s] Alert sent", time.Now().Format("20060102150405"))
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
