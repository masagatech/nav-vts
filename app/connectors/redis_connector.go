package connectors

import (
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/masagatech/nav-vts/app/models"
)

func NewRedis(config *models.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + strconv.Itoa(config.Redis.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
