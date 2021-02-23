package services

import (
	"context"
	"hossipe/core/model"
	"hossipe/core/services"
	"hossipe/infrastructure/rabbitmq"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type rabbitMqPublisher struct {
	client  *rabbitmq.RabbitMqClient
	channel *amqp.Channel
}

func (publisher *rabbitMqPublisher) crateChannel() {
	ch, err := publisher.client.Connection.Channel()
	if err != nil {
		log.Fatalln("Error when trying create channel", err)
	}
	publisher.channel = ch
}

func (f *rabbitMqPublisher) RePublish(msg model.HosepipeMessage, c context.Context) {

}

func (f *rabbitMqPublisher) Close() {
	f.channel.Close()
}

func NewrRabbitMqPublisher(client *rabbitmq.RabbitMqClient, cfg *ErrorMessageSubscriberConfig) services.ErrorMessagePublisher {
	publisher := &rabbitMqPublisher{client: client}
	publisher.crateChannel()
	return publisher
}
