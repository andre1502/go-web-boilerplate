package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (d *Controller) Default(c *gin.Context) {
	d.response.Json(c, http.StatusOK, nil, nil)
}

func (d *Controller) HealthCheck(c *gin.Context) {
	d.response.Json(c, http.StatusOK, map[string]any{"app_name": d.config.AppName}, nil)
}

func (d *Controller) RouteNotFound(c *gin.Context) {
	d.response.Json(c, http.StatusNotFound, nil, nil)
}

func (d *Controller) MethodNotAllowed(c *gin.Context) {
	d.response.Json(c, http.StatusMethodNotAllowed, nil, nil)
}
