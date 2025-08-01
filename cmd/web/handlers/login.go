package handlers

import (
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/loginpage"
)

func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) {
	basePageData := h.NewPageData(r)
	data := loginpage.LoginPageData{
		PageData: basePageData,
	}
	h.Render(r.Context(), w, r, loginpage.Page(data))
}

func (h *Handlers) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	basePageData := h.NewPageData(r)
	data := loginpage.LoginPageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := h.DecodePostForm(r, &data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data.CheckFieldTag(data.Username, "required", "username", data.T("error.field_required"))
	data.CheckFieldTag(data.Password, "required", "password", data.T("error.field_required"))
	data.CheckFieldTag(data.Password, "max=64", "password", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, loginpage.Page(data))
		return
	}

	user, err := h.UserService.AuthenticateUser(data.Username, data.Password)
	if err != nil && errors.Is(err, services.ErrInvalidCredentials) {
		data.AddFieldError("password", data.T("error.invalid_credentials"))
		w.WriteHeader(http.StatusUnauthorized)
		h.Render(r.Context(), w, r, loginpage.Page(data))
		return
	}

	if err != nil && errors.Is(err, services.ErrEmailNotVerified) {
		h.SessionManager.Put(r.Context(), "email-to-verify", user.Email)
		http.Redirect(w, r, "/user/verify-email", http.StatusSeeOther)
		return
	}

	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	err = h.SessionManager.RenewToken(r.Context())
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	h.SessionManager.Put(r.Context(), "authenticatedUserId", user.ID)
	h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.user_welcome", map[string]string{"Username": user.Username}),
			Variant:  "simple",
			Duration: 2000,
		},
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
