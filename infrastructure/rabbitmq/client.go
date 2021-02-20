package rabbitmq

import (
	"github.com/streadway/amqp"
)

type RabbitMqClient struct {
	Connection *amqp.Connection
}

func NewRabbitMqClient(connStr string) (*RabbitMqClient, error) {
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, err
	}
	return &RabbitMqClient{Connection: conn}, nil
}

func (client *RabbitMqClient) Close() {
	client.Connection.Close()
}
