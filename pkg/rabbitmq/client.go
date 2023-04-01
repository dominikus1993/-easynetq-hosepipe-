package rabbitmq

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type RabbitMqClient interface {
	Close()
	Publish(ctx context.Context, exchangeName string, topic string, msg amqp.Publishing) error
	Get(queue string, autoAck bool) (msg amqp.Delivery, ok bool, err error)
}

type rabbitMqClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMqClient(connStr string) (*rabbitMqClient, error) {
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &rabbitMqClient{connection: conn, channel: channel}, nil
}

func (client *rabbitMqClient) Publish(ctx context.Context, exchangeName string, topic string, msg amqp.Publishing) error {
	return client.channel.PublishWithContext(ctx, exchangeName, topic, false, false, msg)
}

func (client *rabbitMqClient) Get(queue string, autoAck bool) (msg amqp.Delivery, ok bool, err error) {
	return client.channel.Get(queue, autoAck)
}

func (client *rabbitMqClient) Close() {
	err := client.connection.Close()
	if err != nil {
		log.WithError(err).Fatalln("error when closing connection")
	}
}

func (client *rabbitMqClient) DeclareExchange(ctx context.Context, exchangeName string) error {
	return client.channel.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

func (client *rabbitMqClient) DeclareQueue(ctx context.Context, queueName string) error {
	_, err := client.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	return err
}

func (client *rabbitMqClient) QueueBind(exchange, queue, topic string) error {
	return client.channel.QueueBind(queue, topic, exchange, true, amqp091.Table{})
}
