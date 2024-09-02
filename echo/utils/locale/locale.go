package locale

import (
	"boilerplate/utils/config"
	"boilerplate/utils/logger"
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Locale struct {
	Bundle    *i18n.Bundle
	Localizer *i18n.Localizer
	Lang      string
}

func NewLocale(config *config.Config) *Locale {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile(fmt.Sprintf("%s/locales/en-US.json", config.RootPath))
	bundle.MustLoadMessageFile(fmt.Sprintf("%s/locales/zh-TW.json", config.RootPath))

	lang := "en-US"

	return &Locale{
		Bundle:    bundle,
		Localizer: i18n.NewLocalizer(bundle, lang),
		Lang:      lang,
	}
}

func (l *Locale) Localize(messageId string, data map[string]any) (message string) {
	var err error

	if data != nil {
		message, err = l.Localizer.Localize(&i18n.LocalizeConfig{
			MessageID:      messageId,
			TemplateData:   data,
			DefaultMessage: l.defaultMessage(messageId),
		})
	} else {
		message, err = l.Localizer.Localize(&i18n.LocalizeConfig{
			MessageID:      messageId,
			DefaultMessage: l.defaultMessage(messageId),
		})
	}

	if err != nil {
		logger.Sugar.Error(err)
	}

	return message
}

func (l *Locale) defaultMessage(messageId string) *i18n.Message {
	return &i18n.Message{
		ID:    messageId,
		Other: messageId,
	}
}
