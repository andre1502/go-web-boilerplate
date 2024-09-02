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

	"github.com/gin-gonic/gin"
)

func (m *Middleware) JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestTkn := m.ExtractToken(c)
		claims, err := token.TokenValid(requestTkn, []byte(m.Config.TokenConfig.SecretKey))
		if err != nil {
			m.Response.Json(c, http.StatusUnauthorized, "", cerror.Fail(cerror.FuncName(), "invalid_token", nil, err))
			return
		}

		key := fmt.Sprintf(constant.MEMBER_TOKEN_KEY, claims.UserId)

		tkn, err := m.Db.Redis.Get(key)
		if err != nil {
			m.Response.Json(c, http.StatusUnauthorized, "", cerror.Fail(cerror.FuncName(), "invalid_token", nil, err))
			return
		}

		if tkn != requestTkn {
			logger.Sugar.Warnf("Not valid comparison of token: [%s] and requestToken: [%s].", tkn, requestTkn)
			m.Response.Json(c, http.StatusUnauthorized, "", cerror.Fail(cerror.FuncName(), "invalid_token", nil, nil))
			return
		}

		c.Set("user_id", claims.UserId)

		c.Next()
	}
}

func (m *Middleware) ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if !utils.IsEmptyString(token) {
		return token
	}

	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}
