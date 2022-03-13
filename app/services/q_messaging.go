package services

import (
	"github.com/streadway/amqp"
)

type MessagingQ struct {
	amqp *amqp.Connection
}

func (q *MessagingQ) SendToQ(qName string, message interface{}) {

	q.amqp.Channel()


}
