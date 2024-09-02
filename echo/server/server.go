package server

import (
	"boilerplate/server/middleware"
	"boilerplate/server/route"
	"boilerplate/utils/config"
	"boilerplate/utils/database"
	"fmt"

	"github.com/labstack/echo/v4"
)

type Server struct {
	Config     *config.Config
	AppPort    string
	Echo       *echo.Echo
	Router     *route.Route
	middleware *middleware.Middleware
}

func NewServer(config *config.Config) *Server {
	database := database.NewDatabase(config)

	server := &Server{
		Config:     config,
		AppPort:    fmt.Sprintf(":%d", config.Port),
		middleware: middleware.NewMiddleware(config, database),
	}

	server.newEchoEngine()

	return server
}
