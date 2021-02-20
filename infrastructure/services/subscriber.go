package services

import (
	"context"
	"encoding/json"
	"hossipe/core/model"
	"hossipe/core/services"
	"hossipe/infrastructure/rabbitmq"

	log "github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

type ErrorMessageSubscriberConfig struct {
	Queue string
}

type rabbitMqSubscriber struct {
	client  *rabbitmq.RabbitMqClient
	channel *amqp.Channel
	config  *ErrorMessageSubscriberConfig
}

func (subscriber *rabbitMqSubscriber) crateChannel() {
	ch, err := subscriber.client.Connection.Channel()
	if err != nil {
		log.Fatalln("Error when trying create channel", err)
	}
	subscriber.channel = ch
}

func (subscriber *rabbitMqSubscriber) subscribe(stream chan model.HosepipeMessage) {
	ch := subscriber.channel

	q, err := ch.QueueDeclare(
		subscriber.config.Queue, // name
		true,                    // durable
		false,                   // delete when usused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name,     // queue
		"hosepipe", // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	for msg := range msgs {
		var res model.HosepipeMessage
		err := json.Unmarshal(msg.Body, &res)
		if err != nil {
			log.Errorln(err)
		} else {
			stream <- res
		}
	}

	close(stream)
}

func (f *rabbitMqSubscriber) Subscribe(c context.Context) chan model.HosepipeMessage {
	stream := make(chan model.HosepipeMessage)
	f.crateChannel()
	go f.subscribe(stream)
	return stream
}

func (f *rabbitMqSubscriber) Close() {
	f.channel.Close()
}

func NewrRabbitMqSubscriber(client *rabbitmq.RabbitMqClient, cfg *ErrorMessageSubscriberConfig) services.ErrorMessageSubscriber {
	return &rabbitMqSubscriber{client: client, config: cfg}
}
