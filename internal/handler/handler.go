package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"

	"github.com/nktinn/TGAlerter/internal/service"
)

type Handler struct {
	services *service.Service
	logger   *zerolog.Logger
	nc       *nats.Conn
}

func NewHandler(services *service.Service, logger *zerolog.Logger, nc *nats.Conn) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
		nc:       nc,
	}
}

func (h *Handler) InitPostRoutes() *fiber.App {
	app := fiber.New()
	app.Post("/alert", h.getAlert)
	return app
}
