package ctxi18n

import (
	"context"

	"github.com/bauerbrun0/nand2tetris-web/internal/appctx"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func GetLocalizer(ctx context.Context) *i18n.Localizer {
	localizer, ok := ctx.Value(appctx.LocalizerKey).(*i18n.Localizer)
	if !ok {
		panic("localizer not found in context")
	}
	return localizer
}

func T(ctx context.Context, key string) string {
	localizer := GetLocalizer(ctx)

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})

	if err != nil {
		return key
	}

	return msg
}

func TTemplate(ctx context.Context, key string, templateData map[string]string) string {
	localizer := GetLocalizer(ctx)

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: templateData,
	})

	if err != nil {
		return key
	}

	return msg
}

func TConfig(ctx context.Context, cfg *i18n.LocalizeConfig) string {
	localizer := GetLocalizer(ctx)
	return localizer.MustLocalize(cfg)
}
