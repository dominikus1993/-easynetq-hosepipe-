package services

import (
	"hossipe/infrastructure/rabbitmq"

	"github.com/streadway/amqp"
)

type rabbitMqPublisher struct {
	client  *rabbitmq.RabbitMqClient
	channel *amqp.Channel
}
