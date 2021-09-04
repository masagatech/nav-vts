package controller

import "github.com/masagatech/nav-vts/app/models"

type Master_controller struct {
	base_controller
}

// Initialize controller constructor
func (o *Master_controller) Initr(app *models.App) {
	// setting app context
	o.super(app)

	// intialize router for controller

	// set api version group
	_ = o.App.Fiber.Group("/api/v1/master")
	// actual route
}
