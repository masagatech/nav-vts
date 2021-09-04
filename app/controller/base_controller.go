package controller

import "github.com/masagatech/nav-vts/app/models"

type base_controller struct {
	App *models.App
}

func (o *base_controller) super(app *models.App) {
	o.App = app
}
