package route

import (
	"boilerplate/server/controller"
	"boilerplate/server/middleware"
	"boilerplate/server/validation"

	"github.com/labstack/echo/v4"
)

type Route struct {
	ech         *echo.Echo
	routerGroup *echo.Group
	middleware  *middleware.Middleware
	controller  *controller.Controller
}

func NewRoute(ech *echo.Echo, middleware *middleware.Middleware, validation *validation.Validation) *Route {
	route := &Route{
		ech:         ech,
		routerGroup: ech.Group("/api"),
		middleware:  middleware,
		controller:  controller.NewController(middleware, validation),
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
