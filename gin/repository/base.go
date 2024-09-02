package repository

import (
	"boilerplate/server/middleware"
	"boilerplate/utils/config"
	"boilerplate/utils/database"

	"gorm.io/gorm"
)

type Repository struct {
	middleware *middleware.Middleware
	config     *config.Config
	db         *database.Database
}

func NewRepository(middleware *middleware.Middleware) *Repository {
	return &Repository{
		middleware: middleware,
		config:     middleware.Config,
		db:         middleware.Db,
	}
}

func (repo *Repository) Paginate(db *gorm.DB) *gorm.DB {
	page := repo.middleware.Pagination.Page

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * repo.middleware.Pagination.PageSize

	return db.Offset(offset).Limit(repo.middleware.Pagination.PageSize)
}
