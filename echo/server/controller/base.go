package controller

import (
	"boilerplate/server/middleware"
	"boilerplate/server/response"
	"boilerplate/server/validation"
	"boilerplate/utils/config"
)

type Controller struct {
	config     *config.Config
	middleware *middleware.Middleware
	validation *validation.Validation
	response   *response.Response
}

func NewController(middleware *middleware.Middleware, validation *validation.Validation) *Controller {
	return &Controller{
		config:     middleware.Config,
		middleware: middleware,
		validation: validation,
		response:   response.NewResponse(middleware.Locale),
	}
}
