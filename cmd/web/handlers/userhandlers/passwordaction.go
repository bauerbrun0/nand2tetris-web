package userhandlers

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
)

func (h *Handlers) validateAndCheckPasswordField(
	w http.ResponseWriter,
	r *http.Request,
	data *usersettingspage.UserSettingsPageData,
	password string,
	passwordFieldName string,
) bool {
	data.CheckFieldTag(password, "required", passwordFieldName, data.T("error.field_required"))
	data.CheckFieldTag(password, "max=64", passwordFieldName, data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return false
	}

	ok, err := h.UserService.CheckUserPassword(data.UserInfo.ID, password)
	if err != nil {
		h.ServerError(w, r, err)
		return false
	}

	if !ok {
		data.AddFieldError(passwordFieldName, data.T("error.incorrect_password"))
		w.WriteHeader(http.StatusUnauthorized)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return false
	}

	return true
}
