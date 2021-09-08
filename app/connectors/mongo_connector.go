package connectors

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/masagatech/nav-vts/app/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDb(config *models.Config) *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://" + config.Database.Host + ":" + strconv.Itoa(config.Database.Port))
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client.Database(config.Database.Database)

}
