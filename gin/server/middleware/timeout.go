package middleware

import (
	"boilerplate/utils/constant"
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) Timeout() gin.HandlerFunc {
	return timeout.New(timeout.WithTimeout(constant.TIMEOUT_SECONDS*time.Second), timeout.WithHandler(func(c *gin.Context) {
		c.Next()
	}), timeout.WithResponse(func(c *gin.Context) {
		c.String(http.StatusServiceUnavailable, "")
	}))
}
