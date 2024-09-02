package middleware

import (
	"boilerplate/utils/constant"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) Paginate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var page int
		var pageSize int
		var err error

		defaultPageStr := strconv.Itoa(constant.PAGE_SIZE)

		if page, err = strconv.Atoi(c.DefaultQuery("page", "0")); err != nil || page < 0 {
			page = 0
		}

		if pageSize, err = strconv.Atoi(c.DefaultQuery("page_size", defaultPageStr)); err != nil || pageSize < 0 {
			pageSize = constant.PAGE_SIZE
		}

		m.Pagination.Page = page
		m.Pagination.PageSize = pageSize

		c.Next()
	}
}

func (m *Middleware) CalculateTotalPage() {
	m.Pagination.TotalPage = int(math.Ceil(float64(m.Pagination.TotalRecord) / float64(m.Pagination.PageSize)))
}
