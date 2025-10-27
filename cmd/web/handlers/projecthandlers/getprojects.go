package projecthandlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/bauerbrun0/nand2tetris-web/internal/apidata"
	"github.com/bauerbrun0/nand2tetris-web/internal/validator"
)

func (h *Handlers) HandleGetProjects(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	v := &validator.Validator{
		Validate: validator.NewValidator(),
	}

	v.CheckFieldTag(pageStr, "required,numeric,min=1", "page", "page must be a positive integer.")
	v.CheckFieldTag(pageSizeStr, "required,numeric,min=1,max=100", "page_size", "page_size must be a positive integer between 1 and 100.")

	if !v.Valid() {
		h.Application.WriteJSONBadRequestError(w, r, v.GetFirstFieldError())
		return
	}

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil {
		h.Application.WriteJSONBadRequestError(w, r, "invalid page parameter")
		return
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		h.Application.WriteJSONBadRequestError(w, r, "invalid page_size parameter")
		return
	}

	h.Logger.Info("get projects", slog.Int("page", int(page)), slog.Int("page_size", int(pageSize)))

	userId := h.Application.GetAuthenticatedUserInfo(r).ID
	projects, totalCount, err := h.Application.ProjectService.GetPaginatedProjects(int32(page), int32(pageSize), userId)

	response := apidata.ProjectsResponse{
		Projects:   projects,
		TotalCount: totalCount,
		Page:       int32(page),
		PageSize:   int32(pageSize),
		TotalPages: (totalCount + int32(pageSize) - 1) / int32(pageSize),
	}

	err = h.Application.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		h.Application.WriteJSONServerError(w, r, err)
		return
	}
}
