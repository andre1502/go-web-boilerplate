package middleware

import (
	"boilerplate/server/response"
	"boilerplate/utils/config"
	"boilerplate/utils/constant"
	"boilerplate/utils/database"
	"boilerplate/utils/locale"
)

type Middleware struct {
	Config     *config.Config
	Db         *database.Database
	Locale     *locale.Locale
	Pagination *response.Pagination
	Response   *response.Response
}

func NewMiddleware(config *config.Config, db *database.Database) *Middleware {
	locale := locale.NewLocale(config)

	return &Middleware{
		Config: config,
		Db:     db,
		Locale: locale,
		Pagination: &response.Pagination{
			Page:      0,
			TotalPage: constant.PAGE_SIZE,
		},
		Response: response.NewResponse(locale),
	}
}
