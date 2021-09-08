package connectors

import (
	"fmt"
	"strconv"

	"github.com/masagatech/nav-vts/app/models"
	"github.com/streadway/amqp"
)

func NewRabbtMq(config *models.Config) *amqp.Connection {

	fmt.Println("amqp://" + config.Rabbitmq.User + ":" + config.Rabbitmq.Pwd + "@" + config.Rabbitmq.Host + ":" + strconv.Itoa(config.Rabbitmq.Port) + "/")
	conn, err := amqp.Dial("amqp://" + config.Rabbitmq.User + ":" + config.Rabbitmq.Pwd + "@" + config.Rabbitmq.Host + ":" + strconv.Itoa(config.Rabbitmq.Port) + "/")
	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to RabbitMQ")
	}
	return conn
}
