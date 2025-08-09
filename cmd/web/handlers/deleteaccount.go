package handlers

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/usersettingspage"
)

func (h *Handlers) handleUserSettingsDeleteAccountPost(w http.ResponseWriter, r *http.Request, data *usersettingspage.UserSettingsPageData) {
	data.CheckFieldTag(data.DeleteAccount.Email, "required", "DeleteAccount.Email", data.T("error.field_required"))
	data.CheckFieldTag(data.DeleteAccount.Email, "email", "DeleteAccount.Email", data.T("error.field_invalid_email"))
	data.CheckFieldTag(
		data.DeleteAccount.Email,
		"max=128",
		"DeleteAccount.Email",
		data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "128"}),
	)
	data.CheckFieldBool(data.DeleteAccount.Email == data.UserInfo.Email, "DeleteAccount.Email", data.T("error.type_current_email"))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	verificationMethod, ok := ParseVerificationMethod(data.Verification)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
		return
	}

	switch verificationMethod {
	case VerificationPassword:
		ok := h.validateAndCheckPasswordField(w, r, data, data.DeleteAccount.Password, "DeleteAccount.Password")
		if ok {
			h.deleteAccount(w, r)
		}
	case VerificationGitHub:
		h.sendGithubActionRedirect(w, r, ActionDeleteAccount, "/user/oauth/github/callback/action")
	case VerificationGoogle:
		h.sendGoogleActionRedirect(w, r, ActionDeleteAccount, "/user/oauth/google/callback/action")
	default:
		w.WriteHeader(http.StatusBadRequest)
		h.Render(r.Context(), w, r, usersettingspage.Page(*data))
	}
}

func (h *Handlers) deleteAccount(w http.ResponseWriter, r *http.Request) {
	data := h.NewPageData(r)
	err := h.UserService.DeleteAccount(data.UserInfo.ID)

	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.T("toast.successfully_deleted_account"),
			Variant:  "success",
			Duration: 2000,
		},
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
