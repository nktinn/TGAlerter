package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nktinn/TGAlerter/internal/service"
	"github.com/rs/zerolog"
)

type Handler struct {
	services *service.Service
	logger   zerolog.Logger
}

func NewHandler(services *service.Service, logger zerolog.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitPostRoutes() *fiber.App {
	app := fiber.New()
	app.Post("/alert", h.getAlert)
	return app
}
