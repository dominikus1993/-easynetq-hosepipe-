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
	} else {
		props := msg.BasicProperties

		if checkIfShouldRepublishMessage(props.Headers) {
			log.Traceln("Message published")
			rMsg := amqp.Publishing{ContentType: "application/json", CorrelationId: props.CorrelationID, Type: props.Type, Headers: props.Headers, Body: body}
			f.channel.Publish("", msg.RoutingKey, false, false, rMsg)
		}

	}

}

func checkIfShouldRepublishMessage(headers map[string]interface{}) bool {
	if val, ok := headers["retry"].(int); ok {
		if val < 10 {
			headers["retry"] = val + 1
			return true
		}
		return false
	}
	headers["retry"] = 1
	return true
}

func (f *rabbitMqPublisher) Close() {
	f.channel.Close()
}

func NewrRabbitMqPublisher(client *rabbitmq.RabbitMqClient, cfg *ErrorMessageSubscriberConfig) services.ErrorMessagePublisher {
	publisher := &rabbitMqPublisher{client: client}
	publisher.crateChannel()
	return publisher
}
