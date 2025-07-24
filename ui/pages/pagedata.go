package pages

import (
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Toast struct {
	Message  string `json:"message"`
	Variant  string `json:"variant"`
	Duration int32  `json:"duration"`
}

type PageData struct {
	IsAuthenticated bool                   `json:"-"`
	UserInfo        *models.GetUserInfoRow `json:"-"`
	CurrentYear     int                    `json:"-"`
	InitialToasts   []Toast                `json:"initialToasts"`
	Localizer       *i18n.Localizer        `json:"-"`
}

func (pd *PageData) T(key string) string {
	msg, err := pd.Localizer.Localize(&i18n.LocalizeConfig{MessageID: key})
	if err != nil {
		return key
	}
	return msg
}

func (pd *PageData) TTemplate(key string, templateData map[string]string) string {
	msg, err := pd.Localizer.Localize(&i18n.LocalizeConfig{MessageID: key, TemplateData: templateData})
	if err != nil {
		return key
	}

	return msg
}

func (pd *PageData) TConfig(cfg *i18n.LocalizeConfig) string {
	return pd.Localizer.MustLocalize(cfg)
}
