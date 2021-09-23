package repository

import (
	"context"
	"fmt"

	"github.com/masagatech/nav-vts/app/constants"
	"github.com/masagatech/nav-vts/app/models"
	"go.mongodb.org/mongo-driver/bson"
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

func (v *VehicleRepository) GetVehicleById(vehicle models.VehicleModel) (res interface{}) {
	fmt.Println(vehicle)
	filter := bson.M{"name": "Ertiga"}

	resultd := v.DB.Collection(constants.Vehicle).FindOne(v.Ctx, filter)
	var inter interface{}
	resultd.Decode(&inter)

	return inter
}

// https://dev.to/mikefmeyer/build-a-go-rest-api-with-fiber-and-mongodb-44og
