package projecthandlers

import (
	"errors"
	"net/http"

	"github.com/bauerbrun0/nand2tetris-web/internal/apidata"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
)

func (h *Handlers) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	var createProjectRequest apidata.CreateProjectRequest
	err := h.Application.ReadJSON(w, r, &createProjectRequest)
	if err != nil {
		h.Application.WriteJSONBadRequestError(w, r, err.Error())
		return
	}

	v := &validator.Validator{
		Validate: validator.NewValidator(),
	}

	v.CheckFieldTag(createProjectRequest.Title, "required", "title", "title is required")
	v.CheckFieldTag(createProjectRequest.Title, "max=100", "title", "title must not be more than 100 characters long")
	v.CheckFieldTag(createProjectRequest.Description, "max=500", "description", "description must not be more than 500 characters long")
	if !v.Valid() {
		h.Application.WriteJSONBadRequestError(w, r, v.GetFirstFieldError())
		return
	}

	userId := h.Application.GetAuthenticatedUserInfo(r).ID

	project, err := h.Application.ProjectService.CreateProject(
		createProjectRequest.Title,
		createProjectRequest.Description,
		userId,
	)

	if err != nil {
		if errors.Is(err, models.ErrProjectTitleTaken) {
			h.Application.WriteJSONBadRequestError(w, r, "project title is already taken")
			return
		}
		h.Application.ServerError(w, r, err)
		return
	}

	h.Application.WriteJSON(w, http.StatusCreated, project, nil)
}
