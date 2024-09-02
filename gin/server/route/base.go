package route

import (
	"boilerplate/server/controller"
	"boilerplate/server/middleware"

	"github.com/gin-gonic/gin"
)

type Route struct {
	gin        *gin.Engine
	apiGroup   *gin.RouterGroup
	middleware *middleware.Middleware
	controller *controller.Controller
}

func NewRoute(middleware *middleware.Middleware, gin *gin.Engine) *Route {
	route := &Route{
		gin:        gin,
		apiGroup:   gin.Group("/api"),
		middleware: middleware,
		controller: controller.NewController(middleware),
	}

	route.addRoutes()

	return route
}

func (r *Route) addRoutes() {
	r.defaultRoutes()
}

func (r *Route) defaultRoutes() {
	r.gin.GET("", r.controller.Default)
	r.gin.POST("", r.controller.Default)
	r.gin.GET("/health", r.controller.HealthCheck)
	r.gin.POST("/health", r.controller.HealthCheck)
	r.gin.NoRoute(r.controller.RouteNotFound)
	r.gin.NoMethod(r.controller.MethodNotAllowed)
}
