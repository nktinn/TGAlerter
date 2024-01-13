package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nktinn/TGAlerter/internal/model"
	"time"
)

func (h *Handler) getAlert(c *fiber.Ctx) error {
	var alert model.Alert
	if err := c.BodyParser(&alert); err != nil {
		h.logger.Error().Msgf("["+time.Now().Format("20060102150405")+"] "+"Error while parsing json body: %s", err.Error())
		return err
	}
	if err := h.services.Alerter.SendAlert(alert); err != nil {
		h.logger.Error().Msgf("["+time.Now().Format("20060102150405")+"] "+"Error while sending alert to telegram: %s", err.Error())
		return err
	}
	h.logger.Info().Msgf("["+time.Now().Format("20060102150405")+"] "+"Alert sent: %s -- %s", alert.ServiceId, alert.Message)
	return nil
}
