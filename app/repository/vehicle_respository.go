package repository

import (
	"context"

	"github.com/masagatech/nav-vts/app/constants"
	"github.com/masagatech/nav-vts/app/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type VehicleRepository struct {
	DB  *mongo.Database
	Ctx context.Context
}

func (v *VehicleRepository) InsertVehicle(vehicle models.VehicleModel) (result *mongo.InsertOneResult, err error) {
	result, err = v.DB.Collection(constants.Vehicle).InsertOne(v.Ctx, vehicle)

	if err != nil {
		return nil, err
	}

	return result, err
}
