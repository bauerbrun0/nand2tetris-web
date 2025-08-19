package userhandlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/resetpasswordentercodepage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/resetpasswordpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/resetpasswordsendcodepage"
)

func (h *Handlers) UserResetPasswordSendCode(w http.ResponseWriter, r *http.Request) {
	basePageData := h.NewPageData(r)
	data := resetpasswordsendcodepage.ResetPasswordSendCodePageData{
		PageData: basePageData,
	}
	h.Render(r.Context(), w, r, resetpasswordsendcodepage.Page(data))
}

func (h *Handlers) UserResetPasswordSendCodePost(w http.ResponseWriter, r *http.Request) {
	basePageData := h.NewPageData(r)
	data := resetpasswordsendcodepage.ResetPasswordSendCodePageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := h.DecodePostForm(r, &data)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	data.CheckFieldTag(data.Email, "required", "email", data.TTemplate("error.x_field_required", map[string]string{"Field": "Email"}))
	data.CheckFieldTag(data.Email, "email", "email", data.T("error.email_field_invalid_email"))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, resetpasswordsendcodepage.Page(data))
		return
	}

	_, err = h.UserService.SendPasswordResetCode(data.Email)
	if err != nil && !errors.Is(err, models.ErrUserDoesNotExist) {
		h.ServerError(w, r, err)
		return
	}

	h.SessionManager.Put(r.Context(), "reset-password-email", data.Email)
	h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.email_sent_to", map[string]string{"Email": data.Email}),
			Variant:  "info",
			Duration: 3000,
		},
	})
	http.Redirect(w, r, "/user/reset-password/enter-code", http.StatusSeeOther)
}

func (h *Handlers) UserResetPasswordEnterCode(w http.ResponseWriter, r *http.Request) {
	email := h.SessionManager.PopString(r.Context(), "reset-password-email")
	basePageData := h.NewPageData(r)
	data := resetpasswordentercodepage.ResetPasswordEnterCodePageData{
		PageData: basePageData,
		Email:    email,
	}
	h.Render(r.Context(), w, r, resetpasswordentercodepage.Page(data))
}

func (h *Handlers) UserResetPasswordEnterCodePost(w http.ResponseWriter, r *http.Request) {
	basePageData := h.NewPageData(r)
	data := resetpasswordentercodepage.ResetPasswordEnterCodePageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := h.DecodePostForm(r, &data)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	data.CheckFieldTag(data.Code, "required", "code", data.TTemplate("error.x_field_required", map[string]string{"Field": "Code"}))
	data.CheckFieldTag(data.Email, "required", "email", data.TTemplate("error.x_field_required", map[string]string{"Field": "Email"}))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, resetpasswordentercodepage.Page(data))
		return
	}

	valid, err := h.UserService.VerifyPasswordResetCode(data.Code)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	if !valid {
		data.AddFieldError("code", data.T("error.provided_code_invalid"))
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, resetpasswordentercodepage.Page(data))
		return
	}

	h.SessionManager.Put(r.Context(), "password-reset-code", data.Code)
	http.Redirect(w, r, "/user/reset-password", http.StatusSeeOther)
}

func (h *Handlers) UserResetPassword(w http.ResponseWriter, r *http.Request) {
	code := h.SessionManager.PopString(r.Context(), "password-reset-code")
	if code == "" {
		http.Redirect(w, r, " /user/reset-password/send-code", http.StatusSeeOther)
	}
	basePageData := h.NewPageData(r)
	data := resetpasswordpage.ResetPasswordPageData{
		PageData: basePageData,
		Code:     code,
	}

	h.Render(r.Context(), w, r, resetpasswordpage.Page(data))
}

func (h *Handlers) UserResetPasswordPost(w http.ResponseWriter, r *http.Request) {
	basePageData := h.NewPageData(r)
	data := resetpasswordpage.ResetPasswordPageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := h.DecodePostForm(r, &data)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	data.CheckFieldTag(data.Code, "required", "code", data.TTemplate("error.x_field_required", map[string]string{"Field": "Code"}))
	data.CheckFieldTag(
		data.NewPassword, "required", "new-password", data.TTemplate("error.x_field_required", map[string]string{"Field": "New password"}),
	)
	data.CheckFieldTag(data.NewPassword, "no_whitespace", "new-password", data.T("error.password_contains_whitespace"))
	data.CheckFieldTag(data.NewPassword, "min=8", "new-password", data.TTemplate("error.field_not_enough_characters", map[string]string{"Min": "8"}))
	data.CheckFieldTag(data.NewPassword, "max=64", "new-password", data.TTemplate("error.field_too_many_characters", map[string]string{"Max": "64"}))
	data.CheckFieldTag(data.NewPasswordConfirmation, "required", "new-password-confirmation", data.T("error.field_required"))
	data.CheckFieldBool(data.NewPassword == data.NewPasswordConfirmation, "new-password", data.T("error.passwords_do_not_match"))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, resetpasswordpage.Page(data))
		return
	}

	request, err := h.UserService.ResetPassword(data.NewPassword, data.Code)

	if err != nil && errors.Is(err, services.ErrPasswordResetCodeInvalid) {
		data.AddFieldError("code", data.T("error.provided_password_reset_code_invalid"))
		w.WriteHeader(http.StatusUnauthorized)
		h.Render(r.Context(), w, r, resetpasswordpage.Page(data))
		return
	}

	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	err = h.SessionManager.Iterate(r.Context(), func(ctx context.Context) error {
		userID := h.SessionManager.GetInt32(ctx, "authenticatedUserId")
		if userID == request.UserID {
			return h.SessionManager.Destroy(ctx)
		}
		return nil
	})
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.T("toast.successfully_changed_password"),
			Variant:  "success",
			Duration: 3000,
		},
	})
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
