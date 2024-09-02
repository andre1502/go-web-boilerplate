package service

import (
	"boilerplate/server/middleware"
	"boilerplate/utils/config"
	"boilerplate/utils/database"
)

type Service struct {
	middleware *middleware.Middleware
	config     *config.Config
	redis      *database.Redis
}

func NewService(middleware *middleware.Middleware) *Service {
	return &Service{
		middleware: middleware,
		config:     middleware.Config,
		redis:      middleware.Db.Redis,
	}
}
