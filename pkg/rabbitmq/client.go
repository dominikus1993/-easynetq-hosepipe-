package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type RabbitMqClient interface {
	GetConnection() *amqp.Connection
	GetChannel() *amqp.Channel
	Close()
	DeclareExchange(ctx context.Context, exchangeName string) error
	Publish(ctx context.Context, exchangeName string, topic string, msg amqp.Publishing) error
}

type rabbitMqClient struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMqClient(connStr string) (*rabbitMqClient, error) {
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &rabbitMqClient{Connection: conn, Channel: ch}, nil
}

func (client *rabbitMqClient) GetChannel() *amqp.Channel {
	return client.Channel
}

func (client *rabbitMqClient) GetConnection() *amqp.Connection {
	return client.Connection
}

func (client *rabbitMqClient) Close() {
	err := client.Channel.Close()
	if err != nil {
		log.WithError(err).Fatalln("channel close error")
	}
	err = client.Connection.Close()
	if err != nil {
		log.WithError(err).Fatalln("client connection close error")
	}
}

func (client *rabbitMqClient) DeclareExchange(ctx context.Context, exchangeName string) error {
	return client.Channel.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

func (client *rabbitMqClient) Publish(ctx context.Context, exchangeName string, topic string, msg amqp.Publishing) error {
	err := client.DeclareExchange(ctx, exchangeName)
	if err != nil {
		return fmt.Errorf("error when declare exchange %w", err)
	}

	return client.Channel.Publish(exchangeName, topic, false, false, msg)
}
