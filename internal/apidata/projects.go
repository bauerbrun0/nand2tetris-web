package apidata

import "time"

type CreateProjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type Project struct {
	ID          int32     `json:"id"`
	UserID      int32     `json:"userId"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type ProjectsResponse struct {
	Projects   []Project `json:"projects"`
	TotalCount int32     `json:"totalCount"`
	Page       int32     `json:"page"`
	PageSize   int32     `json:"pageSize"`
	TotalPages int32     `json:"totalPages"`
}
