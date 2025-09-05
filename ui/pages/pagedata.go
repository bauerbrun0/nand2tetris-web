package pages

import (
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Toast struct {
	Message  string `json:"message"`
	Variant  string `json:"variant"`
	Duration int32  `json:"duration"`
}

type Account string

const (
	AccountGitHub Account = "GitHub"
	AccountGoogle Account = "Google"
)

type UserInfo struct {
	ID             int32
	Username       string
	Email          string
	EmailVerified  bool
	Created        time.Time
	IsPasswordSet  bool
	LinkedAccounts []Account
}

type SveltePage string

var (
	SveltePageNone              = SveltePage("none")
	SveltePageHardwareSimulator = SveltePage("hardware-simulator")
)

type PageData struct {
	IsAuthenticated bool            `json:"-"`
	UserInfo        *UserInfo       `json:"-"`
	CurrentYear     int             `json:"-"`
	InitialToasts   []Toast         `json:"initialToasts"`
	Localizer       *i18n.Localizer `json:"-"`
	CSRFToken       string          `json:"-"`
	SveltePage      SveltePage      `json:"-"`
	ShowFooter      bool            `json:"-"`
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
