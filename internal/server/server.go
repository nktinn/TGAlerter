package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/nktinn/TGAlerter/configs"
)

type Server struct {
	fiberApp  *fiber.App
	serverCfg configs.Server
}

func (s *Server) RunFiber(serverCfg configs.Server, f *fiber.App) error {
	s.fiberApp = f
	s.serverCfg = serverCfg
	return f.Listen(serverCfg.Port)
}

func (s *Server) Shutdown() error {
	return s.fiberApp.Shutdown()
}
