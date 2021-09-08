package router

import (
	"github.com/masagatech/nav-vts/app/controller"
	"github.com/masagatech/nav-vts/app/interfaces"
	"github.com/masagatech/nav-vts/app/models"
)

type Router struct {
	App *models.App
}

func (r *Router) InitializeRouters() {
	// registering dynamic route
	r.registerRouters(
		&controller.Master_controller{},
		&controller.Test_controller{},
	)

}

func (r *Router) registerRouters(i ...interfaces.Controller) {

	// initialize all routers
	for _, v := range i {
		v.Initr(r.App)
	}

}
