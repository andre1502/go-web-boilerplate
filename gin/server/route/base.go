package route

import (
	"boilerplate/server/controller"
	"boilerplate/server/middleware"
	"boilerplate/server/response"
	"boilerplate/utils/config"
	"boilerplate/utils/database"

	"github.com/gin-gonic/gin"
)

type Route struct {
	gin        *gin.Engine
	apiGroup   *gin.RouterGroup
	middleware *middleware.Middleware
	controller *controller.Controller
	// add controller
}

func NewRoute(gin *gin.Engine, cfg *config.Config, db *database.Database, resp *response.Response, mdl *middleware.Middleware) *Route {
	// assign ctrl variable to other controller
	ctrl := controller.NewController(cfg, db, mdl.Pagination, resp)

	route := &Route{
		gin:        gin,
		apiGroup:   gin.Group("/api"),
		middleware: mdl,
		controller: ctrl,
		// add controller
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
