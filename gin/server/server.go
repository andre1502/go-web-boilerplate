package server

import (
	"boilerplate/server/middleware"
	"boilerplate/server/response"
	"boilerplate/server/route"
	"boilerplate/utils/config"
	"boilerplate/utils/database"
	"boilerplate/utils/locale"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	AppPort    string
	config     *config.Config
	db         *database.Database
	locale     *locale.Locale
	Gin        *gin.Engine
	router     *route.Route
	response   *response.Response
	middleware *middleware.Middleware
}

func NewServer(cfg *config.Config) *Server {
	db := database.NewDatabase(cfg)
	lcl := locale.NewLocale(cfg)
	resp := response.NewResponse(lcl)

	server := &Server{
		AppPort:    fmt.Sprintf(":%d", cfg.Port),
		config:     cfg,
		db:         db,
		locale:     lcl,
		Gin:        gin.New(),
		response:   resp,
		middleware: middleware.NewMiddleware(cfg, db, lcl, resp),
	}

	server.newGinEngine()

	return server
}
