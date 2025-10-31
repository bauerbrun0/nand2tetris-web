package userhandlers

import (
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/verifyemailpage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/verifyemailsendcodepage"
)

func (h *Handlers) UserVerifyEmail(w http.ResponseWriter, r *http.Request) {
	basePageData := h.NewPageData(r)
	email := h.SessionManager.PopString(r.Context(), "email-to-verify")
	data := verifyemailpage.VerifyEmailPageData{
		PageData: basePageData,
		Email:    email,
	}
	h.Render(r.Context(), w, r, verifyemailpage.Page(data))
}

func (h *Handlers) UserVerifyEmailPost(w http.ResponseWriter, r *http.Request) {
	basePageData := h.NewPageData(r)
	data := verifyemailpage.VerifyEmailPageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := h.DecodePostForm(r, &data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data.CheckFieldTag(data.Code, "required", "code", data.TTemplate("error.x_field_required", map[string]string{"Field": "Code"}))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, verifyemailpage.Page(data))
		return
	}

	ok, err := h.UserService.VerifyEmail(data.Code)
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	if !ok {
		data.AddFieldError("code", data.T("error.verification_code_invalid"))
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, verifyemailpage.Page(data))
		return
	}
	h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.T("toast.successfully_registered"),
			Variant:  "success",
			Duration: 3000,
		},
	})
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (h *Handlers) UserVerifyEmailResendCode(w http.ResponseWriter, r *http.Request) {
	basePageData := h.NewPageData(r)
	data := verifyemailsendcodepage.VerifyEmailSendCodePageData{
		PageData: basePageData,
	}
	h.Render(r.Context(), w, r, verifyemailsendcodepage.Page(data))
}

func (h *Handlers) UserVerifyEmailResendCodePost(w http.ResponseWriter, r *http.Request) {
	basePageData := h.NewPageData(r)
	data := verifyemailsendcodepage.VerifyEmailSendCodePageData{
		PageData: basePageData,
	}
	data.Validate = validator.NewValidator()

	err := h.DecodePostForm(r, &data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	data.CheckFieldTag(data.Email, "required", "email", data.TTemplate("error.x_field_required", map[string]string{"Field": "Email"}))
	data.CheckFieldTag(data.Email, "email", "email", data.T("error.email_field_invalid_email"))

	if !data.Valid() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.Render(r.Context(), w, r, verifyemailsendcodepage.Page(data))
		return
	}

	_, err = h.UserService.ResendEmailVerificationCode(data.Email)

	if err != nil && !errors.Is(err, models.ErrUserDoesNotExist) && !errors.Is(err, services.ErrEmailAlreadyVerified) {
		h.ServerError(w, r, err)
		return
	}

	h.SessionManager.Put(r.Context(), "email-to-verify", data.Email)
	h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  data.TTemplate("toast.new_email_sent_to", map[string]string{"Email": data.Email}),
			Variant:  "info",
			Duration: 3000,
		},
	})
	http.Redirect(w, r, "/user/verify-email", http.StatusSeeOther)
}
