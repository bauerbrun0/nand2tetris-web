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

func (h *Handlers) HandleUpdateChip(w http.ResponseWriter, r *http.Request) {
	projectId, err := strconv.ParseInt(r.PathValue("projectId"), 10, 32)
	if err != nil {
		h.Application.WriteJSONBadRequestError(w, r, "invalid project id")
		return
	}
	chipId, err := strconv.ParseInt(r.PathValue("chipId"), 10, 32)
	if err != nil {
		h.Application.WriteJSONBadRequestError(w, r, "invalid chip id")
		return
	}

	userId := h.Application.GetAuthenticatedUserInfo(r).ID

	var updateChipRequest *apidata.UpdateChipRequest
	err = h.Application.ReadJSON(w, r, &updateChipRequest)
	if err != nil {
		h.Application.WriteJSONBadRequestError(w, r, err.Error())
		return
	}

	v := &validator.Validator{
		Validate: validator.NewValidator(),
	}

	if updateChipRequest.Name != nil {
		v.CheckFieldTag(updateChipRequest.Name, "required", "name", "name is required")
		v.CheckFieldTag(updateChipRequest.Name, "min=2", "name", "name must be at least 2 characters long")
		v.CheckFieldTag(updateChipRequest.Name, "max=100", "name", "name must not be more than 100 characters long")
		v.CheckFieldTag(updateChipRequest.Name, "no_whitespace", "name", "name must not contain whitespace")
		v.CheckFieldBool(
			!regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(*updateChipRequest.Name),
			"name",
			"name cannot contain special characters",
		)
		v.CheckFieldBool(
			len(*updateChipRequest.Name) > 0 && ((*updateChipRequest.Name)[0] < '0' || (*updateChipRequest.Name)[0] > '9'),
			"name",
			"name cannot start with a number",
		)
	}

	chip, err := h.Application.ChipService.UpdateChip(
		int32(chipId),
		int32(projectId),
		userId,
		updateChipRequest.Name,
		updateChipRequest.Hdl,
	)

	if err != nil {
		if errors.Is(err, services.ErrChipNotFound) {
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

	err = h.Application.WriteJSON(w, http.StatusOK, chip, nil)
	if err != nil {
		h.Application.ServerError(w, r, err)
		return
	}
}
