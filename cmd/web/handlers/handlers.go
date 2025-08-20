package handlers

import (
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/application"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers/userhandlers"
)

type Handlers struct {
	User *userhandlers.Handlers
	*application.Application
}

func NewHandlers(app *application.Application) *Handlers {
	return &Handlers{
		User:        NewUserHandlers(app),
		Application: app,
	}
}

func NewUserHandlers(app *application.Application) *userhandlers.Handlers {
	return &userhandlers.Handlers{
		Application: app,
	}
}
