package controller

import (
	"boilerplate/server/response"
	"boilerplate/service"
	"boilerplate/utils/config"
	"boilerplate/utils/database"
)

type Controller struct {
	config *config.Config
	// add service
	pagination *response.Pagination
	response   *response.Response
}

func NewController(cfg *config.Config, db *database.Database, pagination *response.Pagination, resp *response.Response) *Controller {
	// assign svc variable to other service
	// svc := service.NewService(cfg, db, pagination)
	service.NewService(cfg, db, pagination)

	return &Controller{
		config: cfg,
		// add service
		pagination: pagination,
		response:   resp,
	}
}
