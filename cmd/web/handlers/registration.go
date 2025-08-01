package handlers

import (
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/registerpage"
)

func (h *Handlers) UserRegister(w http.ResponseWriter, r *http.Request) {
	basePageData := h.NewPageData(r)
	data := registerpage.RegisterPageData{
		PageData: basePageData,
	}
	h.Render(r.Context(), w, r, registerpage.Page(data))
}

func (h *Handlers) UserRegisterPost(w http.ResponseWriter, r *http.Request) {
	basePageData := h.NewPageData(r)
	data := registerpage.RegisterPageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := h.DecodePostForm(r, &data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data.CheckFieldTag(data.Email, "required", "email", data.T("error.field_required"))
	data.CheckFieldTag(data.Email, "email", "email", data.T("error.field_invalid_email"))
	data.CheckFieldTag(data.Email, "max=128", "email", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "128"}))
	data.CheckFieldTag(data.Username, "required", "username", data.T("error.field_required"))
	data.CheckFieldTag(data.Username, "min=3", "username", data.TTemplate("error.field_not_enough_characters", map[string]string{"Min": "3"}))
	data.CheckFieldTag(data.Username, "max=64", "username", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}))
	err = data.Validate.Var(data.Username, "email")
	if err == nil {
		data.AddFieldError("username", data.T("error.field_is_email"))
	}
	data.CheckFieldTag(data.Username, "no_whitespace", "username", data.T("error.field_contains_whitespace"))
	data.CheckFieldTag(data.Password, "required", "password", data.T("error.field_required"))
	data.CheckFieldTag(data.Password, "no_whitespace", "password", data.T("error.field_contains_whitespace"))
	data.CheckFieldTag(data.Password, "min=8", "password", data.TTemplate("error.field_not_enough_characters", map[string]string{"Min": "8"}))
	data.CheckFieldTag(data.Password, "max=64", "password", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}))
	data.CheckFieldTag(data.PasswordConfirmation, "required", "password-confirmation", data.T("error.field_required"))
	data.CheckFieldBool(data.Password == data.PasswordConfirmation, "password", data.T("error.passwords_do_not_match"))
	data.CheckFieldTag(data.Terms, "required,eq=on", "terms", data.T("error.terms_conditions"))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, registerpage.Page(data))
		return
	}

	_, err = h.UserService.CreateUser(data.Email, data.Username, data.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			data.AddFieldError("email", data.T("error.email_already_used"))
			w.WriteHeader(http.StatusUnprocessableEntity)
			h.Render(r.Context(), w, r, registerpage.Page(data))
			return
		}
		if errors.Is(err, models.ErrDuplicateUsername) {
			data.AddFieldError("username", data.T("error.username_already_used"))
			w.WriteHeader(http.StatusUnprocessableEntity)
			h.Render(r.Context(), w, r, registerpage.Page(data))
			return
		}
		h.ServerError(w, r, err)
	}

	h.SessionManager.Put(r.Context(), "email-to-verify", data.Email)
	http.Redirect(w, r, "/user/verify-email", http.StatusSeeOther)
}
