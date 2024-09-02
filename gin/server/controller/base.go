package controller

import (
	"boilerplate/server/middleware"
	"boilerplate/server/response"
	"boilerplate/utils/config"
)

type Controller struct {
	middleware *middleware.Middleware
	config     *config.Config
	response   *response.Response
}

func NewController(middleware *middleware.Middleware) *Controller {
	return &Controller{
		middleware: middleware,
		config:     middleware.Config,
		response:   response.NewResponse(middleware.Locale),
	}
}
