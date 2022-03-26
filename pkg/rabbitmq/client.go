package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type RabbitMqClient interface {
	CreateChannel() (*amqp.Channel, error)
	Close()
}

type rabbitMqClient struct {
	connection *amqp.Connection
}

func NewRabbitMqClient(connStr string) (*rabbitMqClient, error) {
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, err
	}
	return &rabbitMqClient{connection: conn}, nil
}

func (client *rabbitMqClient) CreateChannel() (*amqp.Channel, error) {
	return client.connection.Channel()
}

func (client *rabbitMqClient) Close() {
	err := client.connection.Close()
	if err != nil {
		log.WithError(err).Fatalln("error when closing connection")
	}
}

func DeclareExchange(ctx context.Context, channel *amqp.Channel, exchangeName string) error {
	return channel.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}
