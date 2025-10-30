package handlers

import (
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/ui/layouts"
	"github.com/bauerbrun0/nand2tetris-web/ui/pages"
)

func (h *Handlers) Projects(w http.ResponseWriter, r *http.Request) {
	data := h.NewPageData(r)
	data.SveltePage = pages.SveltePageProjects
	data.ShowFooter = true
	h.Render(r.Context(), w, r, layouts.BaseLayout(data))
}
