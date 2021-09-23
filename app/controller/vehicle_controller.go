package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/masagatech/nav-vts/app/models"
	"github.com/masagatech/nav-vts/app/repository"
)

type Vehicle_controller struct {
	base_controller
}

// Initialize controller constructor
func (o *Vehicle_controller) Initr(app *models.App) {
	// setting app context
	o.super(app)

	// intialize router for controller

	// set api version group
	d := o.App.Fiber.Group("vehicle")
	d.Post("/", o.UpsertVehicle)
	d.Get("/get/", o.GetVehicle)
}

func (o *Vehicle_controller) UpsertVehicle(ctx *fiber.Ctx) error {
	vehicleRepository := repository.VehicleRepository{DB: o.App.DB, Ctx: ctx.Context()}
	var vehicleModel = models.VehicleModel{}
	if err := ctx.BodyParser(&vehicleModel); err != nil {
		return err
	}
	vehicleRepository.InsertVehicle(vehicleModel)
	return nil
}

func (o *Vehicle_controller) GetVehicle(ctx *fiber.Ctx) error {

	vehicleRepository := repository.VehicleRepository{DB: o.App.DB, Ctx: ctx.Context()}
	var vehicleModel = models.VehicleModel{}
	// if err := ctx.BodyParser(&vehicleModel); err != nil {
	// 	return ctx.Send([]byte(err.Error()))
	// }
	result := vehicleRepository.GetVehicleById(vehicleModel)

	return ctx.Status(fiber.StatusOK).JSON(result)

}
