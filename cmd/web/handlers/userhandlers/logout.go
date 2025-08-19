package userhandlers

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (h *Handlers) UserLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := h.SessionManager.RenewToken(r.Context())
	if err != nil {
		h.ServerError(w, r, err)
		return
	}

	h.SessionManager.Remove(r.Context(), "authenticatedUserId")
	h.SessionManager.Put(r.Context(), "initialToasts", []pages.Toast{
		{
			Message:  h.GetLocalizer(r).MustLocalize(&i18n.LocalizeConfig{MessageID: "toast.logout"}),
			Variant:  "simple",
			Duration: 2000,
		},
	})
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
