package broker

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"sync"
)

var conn *amqp.Connection
var channel *amqp.Channel
var once sync.Once

func InitMessageBroker() *amqp.Channel {
	once.Do(func() {
		var err error
		conn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
		if err != nil {
			log.Fatal("failed to create RabbitMQ connection")
		}
		channel, err = conn.Channel()
		if err != nil {
			log.Fatal("failed to create RabbitMQ channel")
		}
		log.Print("RabbitMQ connection initialized")
	})
	return channel
}

func Close() {
	var err error
	err = conn.Close()
	err = channel.Close()
	if err != nil {
		log.Fatal("failed to close RabbitMQ connection")
	}
	log.Print("RabbitMQ connection closed")
}
