package repository

import (
	"boilerplate/server/response"
	"boilerplate/utils/config"
	"boilerplate/utils/database"

	"gorm.io/gorm"
)

type Repository struct {
	config     *config.Config
	db         *database.Database
	pagination *response.Pagination
}

func NewRepository(cfg *config.Config, db *database.Database, pagination *response.Pagination) *Repository {
	return &Repository{
		config:     cfg,
		db:         db,
		pagination: pagination,
	}
}

func (repo *Repository) Paginate(db *gorm.DB) *gorm.DB {
	page := repo.pagination.Page

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * repo.pagination.PageSize

	return db.Offset(offset).Limit(repo.pagination.PageSize)
}
