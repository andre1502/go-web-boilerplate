package route

import (
	"boilerplate/server/controller"
	"boilerplate/server/middleware"
	"boilerplate/server/response"
	"boilerplate/server/validation"
	"boilerplate/utils/config"
	"boilerplate/utils/database"

	"github.com/labstack/echo/v4"
)

type Route struct {
	ech         *echo.Echo
	routerGroup *echo.Group
	middleware  *middleware.Middleware
	controller  *controller.Controller
	// add controller
}

func NewRoute(ech *echo.Echo, cfg *config.Config, db *database.Database, vld *validation.Validation, resp *response.Response, mdl *middleware.Middleware) *Route {
	// assign ctrl variable to other controller
	ctrl := controller.NewController(cfg, db, vld, mdl.Pagination, resp)

	route := &Route{
		ech:         ech,
		routerGroup: ech.Group("/api"),
		middleware:  mdl,
		controller:  ctrl,
		// add controller
	}

	route.addRoutes()

	return route
}

func (r *Route) addRoutes() {
	r.defaultRoutes()
}

func (r *Route) defaultRoutes() {
	r.ech.GET("", r.controller.Default)
	r.ech.POST("", r.controller.Default)
	r.ech.GET("/health", r.controller.HealthCheck)
	r.ech.POST("/health", r.controller.HealthCheck)
}
