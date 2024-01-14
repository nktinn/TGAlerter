package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/nktinn/TGAlerter/internal/model"
)

func (h *Handler) getAlert(c *fiber.Ctx) error {
	var alert model.Alert
	if err := c.BodyParser(&alert); err != nil {
		h.logger.Error().Msgf("[%s] Error while parsing json body: %s", time.Now().Format("20060102150405"), err.Error())
		return err
	}
	if err := h.services.Alerter.SendAlert(alert); err != nil {
		h.logger.Error().Msgf("[%s] Error while sending alert to telegram: %s", time.Now().Format("20060102150405"), err.Error())
		return err
	}
	h.logger.Info().Msgf("[%s] Alert sent: %s -- %s", time.Now().Format("20060102150405"), alert.ServiceID, alert.Message)
	return nil
}
