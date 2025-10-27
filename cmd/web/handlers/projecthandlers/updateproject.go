package projecthandlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bauerbrun0/nand2tetris-web/internal/apidata"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
)

func (h *Handlers) HandleUpdateProject(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.Application.WriteJSONBadRequestError(w, r, "invalid project id")
		return
	}

	var updateProjectRequest apidata.UpdateProjectRequest
	err = h.Application.ReadJSON(w, r, &updateProjectRequest)
	if err != nil {
		h.Application.WriteJSONBadRequestError(w, r, err.Error())
		return
	}

	v := &validator.Validator{
		Validate: validator.NewValidator(),
	}

	if updateProjectRequest.Title != nil {
		v.CheckFieldTag(updateProjectRequest.Title, "max=100", "title", "title must not be more than 100 characters long")
	}
	if updateProjectRequest.Description != nil {
		v.CheckFieldTag(updateProjectRequest.Description, "max=500", "description", "description must not be more than 500 characters long")
	}

	if !v.Valid() {
		h.Application.WriteJSONBadRequestError(w, r, v.GetFirstFieldError())
		return
	}

	userId := h.Application.GetAuthenticatedUserInfo(r).ID

	project, err := h.Application.ProjectService.UpdateProject(
		int32(id),
		updateProjectRequest.Title,
		updateProjectRequest.Description,
		userId,
	)

	if err != nil {
		if errors.Is(err, models.ErrProjectTitleTaken) {
			h.Application.WriteJSONBadRequestError(w, r, "project title is already taken")
			return
		}
		if errors.Is(err, services.ErrProjectNotFound) {
			h.Application.WriteJSONNotFoundError(w, r)
			return
		}
		h.Application.ServerError(w, r, err)
		return
	}

	h.Application.WriteJSON(w, http.StatusOK, project, nil)
}
