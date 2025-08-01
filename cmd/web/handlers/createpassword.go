package handlers

import (
	"context"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
)

func (h *Handlers) handleUserSettingsCreatePasswordPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	data.CheckFieldTag(data.CpaPassword, "required", "cpa-password", data.T("error.field_required"))
	data.CheckFieldTag(data.CpaPassword, "no_whitespace", "cpa-password", data.T("error.field_contains_whitespace"))
	data.CheckFieldTag(
		data.CpaPassword, "min=8", "cpa-password", data.TTemplate("error.field_not_enough_characters", map[string]string{"Min": "8"}),
	)
	data.CheckFieldTag(
		data.CpaPassword, "max=64", "cpa-password", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}),
	)
	data.CheckFieldTag(data.CpaPasswordConfirmation, "required", "cpa-password-confirmation", data.T("error.field_required"))
	data.CheckFieldBool(
		data.CpaPassword == data.CpaPasswordConfirmation, "cpa-password", data.T("error.passwords_do_not_match"),
	)

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	err := h.UserService.CreatePassword(data.UserInfo.ID, data.CpaPassword)
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
			Message:  data.T("toast.successfully_created_password"),
			Variant:  "success",
			Duration: 3000,
		},
	})

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
