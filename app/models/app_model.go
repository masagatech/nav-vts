package models

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Fiber *fiber.App
	DB    *mongo.Database
	Redis *redis.Client
}
