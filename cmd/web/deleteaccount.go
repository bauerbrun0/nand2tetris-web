package main

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
)

func (app *application) handleUserSettingsDeleteAccountPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	data.CheckFieldTag(data.DaEmail, "required", "da-email", data.T("error.field_required"))
	data.CheckFieldTag(data.DaEmail, "email", "da-email", data.T("error.field_invalid_email"))
	data.CheckFieldTag(data.DaEmail, "max=128", "da-email", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "128"}))
	data.CheckFieldBool(data.DaEmail == data.UserInfo.Email, "da-email", data.T("error.type_current_email"))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	err := app.userService.DeleteAccount(data.UserInfo.ID)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.T("toast.successfully_deleted_account"),
			Variant:  "success",
			Duration: 2000,
		},
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
