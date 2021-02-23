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

func publish(publisher *rabbitMqPublisher, msg model.HosepipeMessage) {

}

func (f *rabbitMqPublisher) RePublish(msgs chan model.HosepipeMessage, c context.Context) {
	f.crateChannel()
	for msg := range msgs {
		go publish(f, msg)
	}
}

func (f *rabbitMqPublisher) Close() {
	f.channel.Close()
}

func NewrRabbitMqPublisher(client *rabbitmq.RabbitMqClient, cfg *ErrorMessageSubscriberConfig) services.ErrorMessagePublisher {
	return &rabbitMqPublisher{client: client}
}
