package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) getAlert(c *fiber.Ctx) error {
	err := h.nc.Publish("alert", []byte(c.Body()))
	if err != nil {
		h.logger.Error().Msgf("[%s] Error while publishing alert to NATS: %s", time.Now().Format("20060102150405"), err.Error())
		return err
	}
	h.logger.Info().Msgf("[%s] Alert sent to NATS", time.Now().Format("20060102150405"))
	return nil
}
