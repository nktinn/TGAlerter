package main

import (
	"github.com/nktinn/TGAlerter/configs"
	"github.com/nktinn/TGAlerter/internal/handler"
	"github.com/nktinn/TGAlerter/internal/repository"
	"github.com/nktinn/TGAlerter/internal/server"
	"github.com/nktinn/TGAlerter/internal/service"
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//Create log file
	file, fileErr := createLogFile()

	//zerolog logger
	var logger zerolog.Logger
	if fileErr != nil {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Info().Msgf("[" + time.Now().Format("20060102150405") + "] " + "Program started")
		logger.Error().Msgf("["+time.Now().Format("20060102150405")+"] "+"Error while creating log file: %s", fileErr.Error())
	} else {
		multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, file)
		logger = zerolog.New(multi).With().Timestamp().Logger()
		logger.Info().Msgf("[" + time.Now().Format("20060102150405") + "] " + "Program started")
		logger.Info().Msgf("[" + time.Now().Format("20060102150405") + "] " + "Log file created successfully")
	}

	//cfg
	cfg := configs.NewConfig()
	logger.Info().Msgf("[" + time.Now().Format("20060102150405") + "] " + "Config read")

	//telebot
	telegramBot, telegramErr := telebot.NewBot(telebot.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * 60 * 1000},
	})
	if telegramErr != nil {
		logger.Error().Msgf("["+time.Now().Format("20060102150405")+"] "+"Error while connecting telebot: %s", telegramErr.Error())
		return
	}
	logger.Info().Msgf("[" + time.Now().Format("20060102150405") + "] " + "Telebot connection successfull")

	repos := repository.NewRepository()
	services := service.NewService(repos, telegramBot, cfg.Telegram, cfg.Service, logger)
	handlers := handler.NewHandler(services, logger)

	srv := new(server.Server)
	app := handlers.InitPostRoutes()

	//Run listener
	go func() {
		if err := srv.RunFiber(cfg.Server, app); err != nil {
			logger.Error().Msgf("["+time.Now().Format("20060102150405")+"] "+"Error while running listener: %s", err.Error())
			return
		}
	}()

	//Run healthCheckers
	restyClient := srv.RunResty()
	for _, url := range cfg.Health {
		go func(u string) {
			ticker := time.NewTicker(10 * time.Minute)
			defer ticker.Stop()
			for {
				if err := services.Alerter.HealthCheck(restyClient, u); err != nil {
					logger.Error().Msgf("["+time.Now().Format("20060102150405")+"] "+"Error while health checking \"%s\": %s", u, err.Error())
				}
				<-ticker.C
			}
		}(url.Url)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info().Msgf("[" + time.Now().Format("20060102150405") + "] " + "Program stopped")
}

func createLogFile() (*os.File, error) {
	file, err := os.Create("logs/log-" + time.Now().Format("2006-01-02_15-04-05") + ".log")
	return file, err
}
