package pages

import "github.com/bauerbrun0/nand2tetris-web/internal/models"

type Toast struct {
	Message  string `json:"message"`
	Variant  string `json:"variant"`
	Duration int32  `json:"duration"`
}

type PageData struct {
	IsAuthenticated bool                   `json:"-"`
	UserInfo        *models.GetUserInfoRow `json:"-"`
	CurrentYear     int                    `json:"-"`
	InitialToasts   []Toast                `json:"initialToasts"`
}
