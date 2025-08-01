package handlers

import (
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/application"
)

type Handlers struct {
	*application.Application
}

func NewHandlers(app *application.Application) *Handlers {
	return &Handlers{
		Application: app,
	}
}
