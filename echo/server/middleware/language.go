package middleware

import (
	"boilerplate/utils"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (m *Middleware) Language(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lang := c.QueryParam("lang")

		if (len([]rune(lang)) == 0) || (utils.IsEmptyString(lang)) {
			lang = "en-US"
		}

		m.Locale.Localizer = i18n.NewLocalizer(m.Locale.Bundle, lang)
		m.Locale.Lang = lang

		return next(c)
	}
}
