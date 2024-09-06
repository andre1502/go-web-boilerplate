package service

import (
	"boilerplate/repository"
	"boilerplate/server/response"
	"boilerplate/utils/config"
	"boilerplate/utils/database"
)

type Service struct {
	config *config.Config
	redis  *database.Redis
	// add repository
}

func NewService(cfg *config.Config, db *database.Database, pagination *response.Pagination) *Service {
	// assign repo variable to other repository
	// repo := repository.NewRepository(cfg, db, pagination)
	repository.NewRepository(cfg, db, pagination)

	return &Service{
		config: cfg,
		redis:  db.Redis,
		// add repository
	}
}
