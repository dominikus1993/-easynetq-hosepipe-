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
	body, err := msg.GetMessageBody()
	if err != nil {
		log.WithFields(log.Fields{"Exchange": msg.Exchange, "Queue": msg.Queue, "Topic": msg.RoutingKey, "Exception": msg.Exception}).Errorln("Error when trying resend message", err)
	}
	f.channel.Publish("", msg.RoutingKey, false, false, amqp.Publishing{ContentType: "application/json", Body: body})
}

func (f *rabbitMqPublisher) Close() {
	f.channel.Close()
}

func NewrRabbitMqPublisher(client *rabbitmq.RabbitMqClient, cfg *ErrorMessageSubscriberConfig) services.ErrorMessagePublisher {
	publisher := &rabbitMqPublisher{client: client}
	publisher.crateChannel()
	return publisher
}
