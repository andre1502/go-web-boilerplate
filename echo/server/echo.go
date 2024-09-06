package server

import (
	"boilerplate/server/route"
	"boilerplate/utils/constant"
	cerror "boilerplate/utils/error"
	"boilerplate/utils/logger"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func (s *Server) newEchoEngine() {
	switch strings.ToLower(s.config.Environment) {
	case strings.ToLower(constant.PRD):
		s.Echo.Debug = false
	default:
		s.Echo.Debug = true
	}

	s.Echo.HTTPErrorHandler = s.httpErrorHandler

	s.Echo.IPExtractor = echo.ExtractIPFromXFFHeader(
		echo.TrustLoopback(false),
		echo.TrustLinkLocal(false),
		echo.TrustPrivateNet(false),
	)

	s.Echo.Validator = s.validation

	s.Echo.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		Skipper:          echoMiddleware.DefaultSkipper,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodOptions, http.MethodPost, http.MethodPut, http.MethodGet, http.MethodDelete},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
		MaxAge:           86400,
	}))

	s.Echo.Use(echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogError:     true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogRequestID: true,
		HandleError:  true,
		LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
			logger.Sugar.Debug(v)

			return nil
		},
	}))

	s.Echo.Use(echoMiddleware.RecoverWithConfig(echoMiddleware.RecoverConfig{
		Skipper:             echoMiddleware.DefaultSkipper,
		StackSize:           100 << 10, // 100KB
		DisableStackAll:     false,
		DisablePrintStack:   false,
		DisableErrorHandler: false,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logger.Sugar.Error(err, string(stack))

			return err
		},
	}))

	s.Echo.Use(echoMiddleware.TimeoutWithConfig(echoMiddleware.TimeoutConfig{
		Skipper:      echoMiddleware.DefaultSkipper,
		ErrorMessage: s.middleware.Locale.Localize("http_request_timeout", nil),
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			logger.Sugar.Error(err)
		},
		Timeout: constant.TIMEOUT_SECONDS * time.Second,
	}))

	s.Echo.Pre(echoMiddleware.RemoveTrailingSlash())
	s.Echo.Pre(s.middleware.Language)
	s.Echo.Pre(s.middleware.Paginate)

	s.router = route.NewRoute(s.Echo, s.config, s.db, s.validation, s.response, s.middleware)
}

func (s *Server) httpErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	var message any

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message

		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				code = herr.Code
				message = herr.Message
			}
		}
	} else {
		message = s.getOriginalError(err).Error()
	}

	if _, ok := message.(string); ok {
		message = err.Error()
	}

	if !c.Response().Committed {
		s.middleware.Response.Json(c, code, nil, cerror.Fail(cerror.FuncName(), "http_error", map[string]any{"code": code, "message": message}, err))
	}
}

func (s *Server) getOriginalError(err error) error {
	var originalErr = err

	for originalErr != nil {
		var internalErr = errors.Unwrap(originalErr)

		if internalErr == nil {
			break
		}

		originalErr = internalErr
	}

	return originalErr
}
