package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqPublisher interface {
	Publish(ctx context.Context, exchangeName string, topic string, msg amqp.Publishing) error
}

type rabbitMqPublisher struct {
	client RabbitMqClient
}

func NewRabbitMqPublisher(client RabbitMqClient) *rabbitMqPublisher {
	return &rabbitMqPublisher{client: client}
}

func (publisher *rabbitMqPublisher) Publish(ctx context.Context, exchangeName string, topic string, msg amqp.Publishing) error {
	return publisher.client.Publish(ctx, exchangeName, topic, msg)
}
