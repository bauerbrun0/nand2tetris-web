package handlers

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/ui/pages/homepage"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages/landingpage"
)

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	data := h.NewPageData(r)
	if h.IsAuthenticated(r) {
		h.Render(r.Context(), w, r, homepage.Page(data))
		return
	}
	h.Render(r.Context(), w, r, landingpage.Page(data))
}
