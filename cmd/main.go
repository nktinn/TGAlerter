package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"

	"github.com/nktinn/TGAlerter/configs"
	"github.com/nktinn/TGAlerter/internal/handler"
	"github.com/nktinn/TGAlerter/internal/repository"
	"github.com/nktinn/TGAlerter/internal/server"
	"github.com/nktinn/TGAlerter/internal/service"
)

func main() {
	// Create log file
	file, fileErr := createLogFile()

	// Start logger
	var logger zerolog.Logger
	if fileErr != nil {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Info().Msgf("[%s] Program started", time.Now().Format("20060102150405"))
		logger.Error().Msgf("[%s] Error while creating log file: %s", time.Now().Format("20060102150405"), fileErr.Error())
	} else {
		multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, file)
		logger = zerolog.New(multi).With().Timestamp().Logger()
		logger.Info().Msgf("[%s] Program started", time.Now().Format("20060102150405"))
		logger.Info().Msgf("[%s] Log file created successfully", time.Now().Format("20060102150405"))
	}

	// Start NATS connection
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		logger.Error().Msgf("[%s] Error while connecting NATS: %s", time.Now().Format("20060102150405"), fileErr.Error())
		return
	}
	logger.Info().Msgf("[%s] NATS connection successful", time.Now().Format("20060102150405"))

	// Close NATS connection
	defer func(nc *nats.Conn) {
		err := nc.Drain()
		defer nc.Close()
		if err != nil {
			logger.Error().Msgf("[%s] Error while draining NATS: %s", time.Now().Format("20060102150405"), fileErr.Error())
			return
		}
	}(nc)

	// Config
	cfg := configs.NewConfig()
	logger.Info().Msgf("[%s] Config read", time.Now().Format("20060102150405"))

	// Database
	db, dbErr := repository.NewPostgresDB(cfg.Database)
	if dbErr != nil {
		logger.Error().Msgf("[%s] Error while connecting to database: %s", time.Now().Format("20060102150405"), dbErr.Error())
		//return
	}
	logger.Info().Msgf("[%s] Database connection successful", time.Now().Format("20060102150405"))

	// Start telegram bot
	telegramBot, telegramErr := telebot.NewBot(telebot.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * 60 * 1000},
	})
	if telegramErr != nil {
		logger.Error().Msgf("[%s] Error while connecting telebot: %s", time.Now().Format("20060102150405"), telegramErr.Error())
		return
	}
	logger.Info().Msgf("[%s] Telebot connection successful", time.Now().Format("20060102150405"))

	// Services and handlers
	repo := repository.NewRepository(db)
	services := service.NewService(repo, telegramBot, cfg.Telegram, cfg.Service, &logger)
	handlers := handler.NewHandler(services, &logger, nc)

	// Subscribe NATS
	sub, err := nc.Subscribe("alert", func(msg *nats.Msg) {
		logger.Info().Msgf("[%s] Received a message: %s", time.Now().Format("20060102150405"), string(msg.Data))
		services.Alerter.SendAlert(msg)
	})
	if err != nil {
		logger.Error().Msgf("[%s] Error while subscribing to NATS: %s", time.Now().Format("20060102150405"), fileErr.Error())
		return
	}
	logger.Info().Msgf("[%s] Subscribed to NATS", time.Now().Format("20060102150405"))

	defer func(sub *nats.Subscription) {
		err := sub.Unsubscribe()
		if err != nil {
			logger.Error().Msgf("[%s] Error while unsubscribing from NATS: %s", time.Now().Format("20060102150405"), fileErr.Error())
			return
		}
		logger.Info().Msgf("[%s] Unsubscribed from NATS", time.Now().Format("20060102150405"))
	}(sub)

	srv := new(server.Server)
	// Run listener
	go func() {
		if err := srv.RunFiber(cfg.Server, handlers.InitPostRoutes()); err != nil {
			logger.Error().Msgf("[%s] Error while running listener: %s", time.Now().Format("20060102150405"), err.Error())
			return
		}
	}()

	// Run healthCheckers
	go services.Alerter.HealthCheckWorker(cfg.Health)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(); err != nil {
		logger.Info().Msgf("[%s] Program stopped with error while shutting down server", time.Now().Format("20060102150405"))
	} else {
		logger.Info().Msgf("[%s] Program stopped", time.Now().Format("20060102150405"))
	}
}

func createLogFile() (*os.File, error) {
	file, err := os.Create("logs/log-" + time.Now().Format("2006-01-02_15-04-05") + ".log")
	return file, err
}
