package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (d *Controller) Default(c echo.Context) error {
	return d.response.Json(c, http.StatusOK, nil, nil)
}

func (d *Controller) HealthCheck(c echo.Context) error {
	return d.response.Json(c, http.StatusOK, map[string]any{"app_name": d.config.AppName}, nil)
}
