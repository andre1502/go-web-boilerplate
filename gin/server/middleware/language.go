package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (m *Middleware) Language() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang, exists := c.GetQuery("lang")
		lang = strings.TrimSpace(lang)

		if (len([]rune(lang)) == 0) || (!exists) {
			lang = "en-US"
		}

		m.Locale.Localizer = i18n.NewLocalizer(m.Locale.Bundle, lang)
		m.Locale.Lang = lang

		c.Next()
	}
}
