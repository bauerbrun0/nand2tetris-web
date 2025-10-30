package projecthandlers

import (
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/services"
)

func (h *Handlers) HandleGetProjectBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	project, err := h.Application.ProjectService.GetProjectBySlug(slug, h.Application.GetAuthenticatedUserInfo(r).ID)
	if err != nil {
		if errors.Is(err, services.ErrProjectNotFound) {
			h.Application.WriteJSONNotFoundError(w, r)
			return
		}
		h.Application.WriteJSONServerError(w, r, err)
		return
	}

	err = h.Application.WriteJSON(w, http.StatusOK, project, nil)
	if err != nil {
		h.Application.WriteJSONServerError(w, r, err)
		return
	}
}
