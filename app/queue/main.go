package queue

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// Server ...
type RMQ struct {
	ampqClient *amqp.Connection
	ch         *amqp.Channel
	queue      amqp.Queue
}

func New(ampqClient *amqp.Connection) *RMQ {
	//ch, _ := r.AmpqClient.Channel()
	channel, _ := ampqClient.Channel()

	// , channel *amqp.Channel, queue amqp.Queue
	q, _ := channel.QueueDeclare(
		"test", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)

	return &RMQ{
		ampqClient: ampqClient,
		ch:         channel,
		queue:      q,
	}

}

func (r *RMQ) PublishOnQueue(exchange string, queueName string, message interface{}) {

	msg, err := r.getBytes(message)
	if err != nil {
		fmt.Println(msg)
	}

	_ = r.ch.Publish(
		exchange,  // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
}

func (r *RMQ) getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *RMQ) Listner(exchange string, queueName string) {

	err := r.ch.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	)
	pages, err := r.ch.Consume(exchange, queueName, false, false, false, false, amqp.Table{})
	if err != nil {
		log.Fatalf("basic.consume: %v", err)
	}

	go func() {
		for log := range pages {
			// ... this consumer is responsible for sending pages per log
			fmt.Println(log.ConsumerTag, interface{}(log.Body))
			log.Ack(false)
		}
	}()
}
