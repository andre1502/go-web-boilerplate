package middleware

import (
	"boilerplate/server/response"
	"boilerplate/utils/config"
	cconstant "boilerplate/utils/constant"
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

func NewMiddleware(cfg *config.Config, db *database.Database, lcl *locale.Locale, resp *response.Response) *Middleware {
	return &Middleware{
		Config: cfg,
		Db:     db,
		Locale: lcl,
		Pagination: &response.Pagination{
			Page:     0,
			PageSize: cconstant.PAGE_SIZE,
		},
		Response: resp,
	}
}
