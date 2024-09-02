package server

import (
	"boilerplate/server/route"
	"boilerplate/server/validation"
	"boilerplate/utils/constant"
	"boilerplate/utils/logger"
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (s *Server) newGinEngine() {
	switch strings.ToLower(s.Config.Environment) {
	case strings.ToLower(constant.PRD):
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	validation := validation.NewValidation()

	s.Gin = gin.New()
	s.Gin.SetTrustedProxies(nil)
	s.Gin.TrustedPlatform = "X-Forwarded-For"
	s.Gin.HandleMethodNotAllowed = true
	s.Gin.ForwardedByClientIP = true

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(validation.GetJsonTagName())
		v.RegisterValidation("empty_string", validation.EmptyString())
	}

	s.Gin.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodOptions, http.MethodPost, http.MethodPut, http.MethodGet, http.MethodDelete},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}))

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	s.Gin.Use(ginzap.GinzapWithConfig(logger.Logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			fields := []zapcore.Field{}
			// log request ID
			if requestID := c.Request.Header.Get("X-Request-Id"); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}

			// log request body
			var body []byte
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ = io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)
			fields = append(fields, zap.String("body", string(body)))

			return fields
		}),
	}))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	s.Gin.Use(ginzap.RecoveryWithZap(logger.Logger, true))

	s.Gin.Use(s.middleware.Timeout())
	s.Gin.Use(s.middleware.Language())
	s.Gin.Use(s.middleware.Paginate())

	s.Router = route.NewRoute(s.middleware, s.Gin)
}
