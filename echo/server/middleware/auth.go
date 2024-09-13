package middleware

import (
	"boilerplate/utils"
	"boilerplate/utils/constant"
	cerror "boilerplate/utils/error"
	"boilerplate/utils/logger"
	"boilerplate/utils/token"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (m *Middleware) JwtAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error

		requestTkn := m.ExtractToken(c)
		claims, err := token.TokenValid(requestTkn, []byte(m.Config.TokenConfig.SecretKey))
		if err != nil {
			return m.Response.Json(c, http.StatusUnauthorized, "", cerror.Fail(cerror.FuncName(), "invalid_token", nil, err))
		}

		key := fmt.Sprintf(constant.MEMBER_TOKEN_KEY, claims.UserId)

		tkn, err := m.Db.Redis.Get(key)
		if err != nil {
			return m.Response.Json(c, http.StatusUnauthorized, "", cerror.Fail(cerror.FuncName(), "invalid_token", nil, err))
		}

		if tkn != requestTkn {
			logger.Sugar.Errorf("Not valid comparison of token: [%s] and requestToken: [%s].", tkn, requestTkn)
			return m.Response.Json(c, http.StatusUnauthorized, "", cerror.Fail(cerror.FuncName(), "invalid_token", nil, nil))
		}

		c.Set("user_id", claims.UserId)

		return next(c)
	}
}

func (m *Middleware) ExtractToken(c echo.Context) string {
	token := c.QueryParam("token")
	if !utils.IsEmptyString(token) {
		return token
	}

	bearerToken := c.Request().Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}
