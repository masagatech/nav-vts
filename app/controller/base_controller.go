package controller

import (
	"fmt"

	"github.com/masagatech/nav-vts/app/models"
)

type base_controller struct {
	App *models.App
}

func (o *base_controller) super(app *models.App) {
	fmt.Println("Registering the router")

	o.App = app
}
