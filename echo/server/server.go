package server

import (
	"boilerplate/server/middleware"
	"boilerplate/server/response"
	"boilerplate/server/route"
	"boilerplate/server/validation"
	"boilerplate/utils/config"
	"boilerplate/utils/database"
	"boilerplate/utils/locale"
	"fmt"

	"github.com/labstack/echo/v4"
)

type Server struct {
	AppPort    string
	config     *config.Config
	db         *database.Database
	locale     *locale.Locale
	Echo       *echo.Echo
	router     *route.Route
	validation *validation.Validation
	response   *response.Response
	middleware *middleware.Middleware
}

func NewServer(cfg *config.Config) *Server {
	db := database.NewDatabase(cfg)
	lcl := locale.NewLocale(cfg)
	resp := response.NewResponse(lcl)
	vld := validation.NewValidation()

	server := &Server{
		AppPort:    fmt.Sprintf(":%d", cfg.Port),
		config:     cfg,
		db:         db,
		locale:     lcl,
		Echo:       echo.New(),
		validation: vld,
		response:   resp,
		middleware: middleware.NewMiddleware(cfg, db, lcl, resp),
	}

	server.newEchoEngine()

	return server
}
