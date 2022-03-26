package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type RabbitMqPublisher interface {
	Publish(ctx context.Context, exchangeName string, topic string, msg amqp.Publishing) error
	CloseChannel()
}

type rabbitMqPublisher struct {
	channel *amqp.Channel
}

func NewRabbitMqPublisher(client RabbitMqClient) (*rabbitMqPublisher, error) {
	channel, err := client.CreateChannel()
	if err != nil {
		return nil, fmt.Errorf("error when creating channel %w", err)
	}
	return &rabbitMqPublisher{channel: channel}, nil
}

func (client *rabbitMqPublisher) Publish(ctx context.Context, exchangeName string, topic string, msg amqp.Publishing) error {
	err := DeclareExchange(ctx, client.channel, exchangeName)
	if err != nil {
		return fmt.Errorf("error when declare exchange %w", err)
	}

	return client.channel.Publish(exchangeName, topic, false, false, msg)
}

func (client *rabbitMqPublisher) CloseChannel() {
	err := client.channel.Close()
	if err != nil {
		log.WithError(err).Fatalln("error when closing channel")
	}
}
