package emails

import "github.com/nicksnyder/go-i18n/v2/i18n"

type EmailData struct {
	CurrentYear int
	Localizer   *i18n.Localizer
	BaseUrl     string
}

func (ed *EmailData) T(key string) string {
	msg, err := ed.Localizer.Localize(&i18n.LocalizeConfig{MessageID: key})
	if err != nil {
		return key
	}
	return msg
}

func (ed *EmailData) TTemplate(key string, templateData map[string]string) string {
	msg, err := ed.Localizer.Localize(&i18n.LocalizeConfig{MessageID: key, TemplateData: templateData})
	if err != nil {
		return key
	}

	return msg
}

func (ed *EmailData) TConfig(cfg *i18n.LocalizeConfig) string {
	return ed.Localizer.MustLocalize(cfg)
}
