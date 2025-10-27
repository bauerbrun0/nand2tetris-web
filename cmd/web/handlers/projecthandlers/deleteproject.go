package projecthandlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bauerbrun0/nand2tetris-web/internal/services"
)

func (h *Handlers) HandleDeleteProject(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.Application.WriteJSONBadRequestError(w, r, "invalid project id")
		return
	}

	project, err := h.Application.ProjectService.DeleteProject(int32(id), h.Application.GetAuthenticatedUserInfo(r).ID)
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
