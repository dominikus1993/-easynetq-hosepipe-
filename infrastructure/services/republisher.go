package services

import (
	"context"

	"github.com/dominikus1993/easynetq-hosepipe/core/model"

	"github.com/dominikus1993/easynetq-hosepipe/core/services"

	"github.com/dominikus1993/easynetq-hosepipe/infrastructure/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type rabbitMqPublisher struct {
	client  *rabbitmq.RabbitMqClient
	channel *amqp.Channel
}

func (publisher *rabbitMqPublisher) crateChannel() {
	ch, err := publisher.client.Connection.Channel()
	if err != nil {
		log.WithError(err).Fatalln("Error when trying create channel")
	}
	publisher.channel = ch
}

func (f *rabbitMqPublisher) RePublish(msg model.HosepipeMessage, c context.Context) error {
	body, err := msg.GetMessageBody()
	if err != nil {
		fields := log.Fields{"Exchange": msg.Exchange, "Queue": msg.Queue, "Topic": msg.RoutingKey, "Exception": msg.Exception}
		log.WithFields(fields).WithError(err).WithContext(c).Errorln("Error when trying resend message")
		return err
	}

	props := msg.BasicProperties

	if checkIfShouldRepublishMessage(props.Headers) {
		log.WithContext(c).Traceln("Message published")
		rMsg := amqp.Publishing{ContentType: "application/json", CorrelationId: props.CorrelationID, Type: props.Type, Headers: props.Headers, Body: body}
		return f.channel.Publish("", msg.RoutingKey, false, false, rMsg)
	}
	return nil
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
