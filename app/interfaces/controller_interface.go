package interfaces

import "github.com/masagatech/nav-vts/app/models"

type Controller interface {
	Initr(app *models.App)
}
