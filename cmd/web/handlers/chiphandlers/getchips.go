package chiphandlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bauerbrun0/nand2tetris-web/internal/services"
)

func (h *Handlers) HandleGetChips(w http.ResponseWriter, r *http.Request) {
	projectId, err := strconv.ParseInt(r.PathValue("projectId"), 10, 32)
	if err != nil {
		h.Application.WriteJSONBadRequestError(w, r, "invalid project id")
		return
	}
	userId := h.Application.GetAuthenticatedUserInfo(r).ID

	chips, err := h.Application.ChipService.GetChips(int32(projectId), userId)
	if err != nil {
		if errors.Is(err, services.ErrProjectNotFound) {
			h.Application.WriteJSONNotFoundError(w, r)
			return
		}
		h.Application.ServerError(w, r, err)
		return
	}

	err = h.Application.WriteJSON(w, http.StatusOK, chips, nil)
	if err != nil {
		h.Application.ServerError(w, r, err)
		return
	}
}
