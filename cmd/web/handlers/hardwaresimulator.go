package handlers

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/ui/layouts"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
)

func (h *Handlers) HardwareSimulator(w http.ResponseWriter, r *http.Request) {
	data := h.NewPageData(r)
	data.SveltePage = pages.SveltePageHardwareSimulator
	data.ShowFooter = false
	h.Render(r.Context(), w, r, layouts.BaseLayout(data))
}
