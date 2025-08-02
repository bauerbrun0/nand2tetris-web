package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
)

func (h *Handlers) handleUserSettingsChangePasswordPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	data.CheckFieldTag(data.ChangePassword.CurrentPassword, "required", "ChangePassword.CurrentPassword", data.T("error.field_required"))
	data.CheckFieldTag(
		data.ChangePassword.CurrentPassword, "max=64", "ChangePassword.CurrentPassword", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}),
	)

	data.CheckFieldTag(data.ChangePassword.NewPassword, "required", "ChangePassword.NewPassword", data.T("error.field_required"))
	data.CheckFieldTag(data.ChangePassword.NewPassword, "no_whitespace", "ChangePassword.NewPassword", data.T("error.field_contains_whitespace"))
	data.CheckFieldTag(
		data.ChangePassword.NewPassword, "min=8", "ChangePassword.NewPassword", data.TTemplate("error.field_not_enough_characters", map[string]string{"Min": "8"}),
	)
	data.CheckFieldTag(
		data.ChangePassword.NewPassword, "max=64", "ChangePassword.NewPassword", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}),
	)
	data.CheckFieldTag(data.ChangePassword.NewPasswordConfirmation, "required", "ChangePassword.NewPasswordConfirmation", data.T("error.field_required"))
	data.CheckFieldBool(
		data.ChangePassword.NewPassword == data.ChangePassword.NewPasswordConfirmation, "ChangePassword.NewPassword", data.T("error.passwords_do_not_match"),
	)

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	err := h.UserService.ChangePassword(data.UserInfo.ID, data.ChangePassword.CurrentPassword, data.ChangePassword.NewPassword)
	if err != nil && errors.Is(err, services.ErrInvalidCredentials) {
		data.AddFieldError("ChangePassword.CurrentPassword", data.T("error.incorrect_password"))
		w.WriteHeader(http.StatusUnauthorized)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	err = h.SessionManager.Iterate(r.Context(), func(ctx context.Context) error {
		userID := h.SessionManager.GetInt32(ctx, "authenticatedUserId")
		if userID == data.UserInfo.ID {
			return h.SessionManager.Destroy(ctx)
		}
		return nil
	})
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	h.SessionManager.Destroy(r.Context())
	h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.T("toast.successfully_changed_password"),
			Variant:  "success",
			Duration: 3000,
		},
	})

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
