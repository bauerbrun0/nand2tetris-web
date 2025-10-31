package apidata

import "time"

type CreateChipRequest struct {
	Name *string `json:"name"`
}

type Chip struct {
	ID        int32     `json:"id"`
	ProjectID int32     `json:"projectId"`
	Name      string    `json:"name"`
	Hdl       string    `json:"hdl"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

type UpdateChipRequest struct {
	Name *string `json:"name"`
	Hdl  *string `json:"hdl"`
}
