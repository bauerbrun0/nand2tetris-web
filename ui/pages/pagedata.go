package pages

import "github.com/bauerbrun0/nand2tetris-web/internal/models"

type PageData struct {
	IsAuthenticated bool
	UserInfo        *models.GetUserInfoRow
	CurrentYear     int
}
