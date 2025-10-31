package handlers

import (
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/application"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers/chiphandlers"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers/projecthandlers"
	"github.com/bauerbrun0/nand2tetris-web/cmd/web/handlers/userhandlers"
)

type Handlers struct {
	User    *userhandlers.Handlers
	Project *projecthandlers.Handlers
	Chip    *chiphandlers.Handlers
	*application.Application
}

func NewHandlers(app *application.Application) *Handlers {
	return &Handlers{
		User:        NewUserHandlers(app),
		Project:     NewProjectHandlers(app),
		Chip:        NewChipHandlers(app),
		Application: app,
	}
}

func NewUserHandlers(app *application.Application) *userhandlers.Handlers {
	return &userhandlers.Handlers{
		Application: app,
	}
}

func NewProjectHandlers(app *application.Application) *projecthandlers.Handlers {
	return &projecthandlers.Handlers{
		Application: app,
	}
}

func NewChipHandlers(app *application.Application) *chiphandlers.Handlers {
	return &chiphandlers.Handlers{
		Application: app,
	}
}
