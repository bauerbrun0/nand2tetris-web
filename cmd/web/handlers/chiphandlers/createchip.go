package chiphandlers

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/bauerbrun0/nand2tetris-web/internal/apidata"
	"github.com/bauerbrun0/nand2tetris-web/internal/models"
	"github.com/bauerbrun0/nand2tetris-web/internal/services"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
)

func (h *Handlers) HandleCreateChip(w http.ResponseWriter, r *http.Request) {
	projectId, err := strconv.ParseInt(r.PathValue("projectId"), 10, 32)
	if err != nil {
		h.Application.WriteJSONBadRequestError(w, r, "invalid project id")
		return
	}

	var createChipRequest apidata.CreateChipRequest
	err = h.Application.ReadJSON(w, r, &createChipRequest)
	if err != nil {
		h.Application.WriteJSONBadRequestError(w, r, err.Error())
		return
	}

	v := &validator.Validator{
		Validate: validator.NewValidator(),
	}

	if createChipRequest.Name == nil {
		h.Application.WriteJSONBadRequestError(w, r, "name is required")
		return
	}

	v.CheckFieldTag(projectId, "number,gte=0", "projectId", "projectId must be a positive integer")
	v.CheckFieldTag(createChipRequest.Name, "required", "name", "name is required")
	v.CheckFieldTag(createChipRequest.Name, "min=2", "name", "name must be at least 2 characters long")
	v.CheckFieldTag(createChipRequest.Name, "max=100", "name", "name must not be more than 100 characters long")
	v.CheckFieldTag(createChipRequest.Name, "no_whitespace", "name", "name must not contain whitespace")
	v.CheckFieldBool(
		!regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(*createChipRequest.Name),
		"name",
		"name cannot contain special characters",
	)
	v.CheckFieldBool(
		len(*createChipRequest.Name) > 0 && ((*createChipRequest.Name)[0] < '0' || (*createChipRequest.Name)[0] > '9'),
		"name",
		"name cannot start with a number",
	)

	if !v.Valid() {
		h.Application.WriteJSONBadRequestError(w, r, v.GetFirstFieldError())
		return
	}

	userId := h.Application.GetAuthenticatedUserInfo(r).ID

	chip, err := h.Application.ChipService.CreateChip(
		*createChipRequest.Name,
		int32(projectId),
		userId,
	)

	if err != nil {
		if errors.Is(err, services.ErrProjectNotFound) {
			h.Application.WriteJSONNotFoundError(w, r)
			return
		}
		if errors.Is(err, models.ErrChipNameTaken) {
			h.Application.WriteJSONBadRequestError(w, r, "chip name is already taken")
			return
		}
		h.Application.ServerError(w, r, err)
		return
	}

	err = h.Application.WriteJSON(w, http.StatusCreated, chip, nil)
	if err != nil {
		h.Application.ServerError(w, r, err)
		return
	}
}
