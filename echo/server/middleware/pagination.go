package middleware

import (
	"boilerplate/utils"
	"boilerplate/utils/constant"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (m *Middleware) Paginate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var page int
		var pageSize int
		var err error

		pageStr := c.QueryParam("page")
		pageSizeStr := c.QueryParam("page_size")

		if utils.IsEmptyString(pageStr) {
			page = 0
		} else if page, err = strconv.Atoi(pageStr); err != nil || page < 0 {
			if utils.IsEmptyString(pageSizeStr) {
				page = 0
			} else {
				page = 1
			}
		}

		if utils.IsEmptyString(pageSizeStr) {
			pageSize = constant.PAGE_SIZE
		} else if pageSize, err = strconv.Atoi(pageSizeStr); err != nil || pageSize < 0 {
			pageSize = constant.PAGE_SIZE
		}

		m.Pagination.Page = page
		m.Pagination.PageSize = pageSize

		return next(c)
	}
}
