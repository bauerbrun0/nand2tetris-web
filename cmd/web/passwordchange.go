package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
)

func (app *application) handleUserSettingsChangePasswordPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	data.CheckFieldTag(data.ChpCurrentPassword, "required", "chp-current-password", data.T("error.field_required"))
	data.CheckFieldTag(
		data.ChpCurrentPassword, "max=64", "chp-current-password", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}),
	)

	data.CheckFieldTag(data.ChpNewPassword, "required", "chp-new-password", data.T("error.field_required"))
	data.CheckFieldTag(data.ChpNewPassword, "no_whitespace", "chp-new-password", data.T("error.field_contains_whitespace"))
	data.CheckFieldTag(
		data.ChpNewPassword, "min=8", "chp-new-password", data.TTemplate("error.field_not_enough_characters", map[string]string{"Min": "8"}),
	)
	data.CheckFieldTag(
		data.ChpNewPassword, "max=64", "chp-new-password", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}),
	)
	data.CheckFieldTag(data.ChpNewPasswordConfirmation, "required", "chp-new-password-confirmation", data.T("error.field_required"))
	data.CheckFieldBool(
		data.ChpNewPassword == data.ChpNewPasswordConfirmation, "chp-new-password", data.T("error.passwords_do_not_match"),
	)

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	err := app.userService.ChangePassword(data.UserInfo.ID, data.ChpCurrentPassword, data.ChpNewPassword)
	if err != nil && errors.Is(err, services.ErrInvalidCredentials) {
		data.AddFieldError("chp-current-password", data.T("error.incorrect_password"))
		w.WriteHeader(http.StatusUnauthorized)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.sessionManager.Iterate(r.Context(), func(ctx context.Context) error {
		userID := app.sessionManager.GetInt32(ctx, "authenticatedUserId")
		if userID == data.UserInfo.ID {
			return app.sessionManager.Destroy(ctx)
		}
		return nil
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Destroy(r.Context())
	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.T("toast.successfully_changed_password"),
			Variant:  "success",
			Duration: 3000,
		},
	})

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
