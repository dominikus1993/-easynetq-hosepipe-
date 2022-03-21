package subscriber

import (
	"context"
	"encoding/json"

	"github.com/dominikus1993/easynetq-hosepipe/pkg/config"
	"github.com/dominikus1993/easynetq-hosepipe/pkg/data"
	"github.com/dominikus1993/easynetq-hosepipe/pkg/rabbitmq"
	log "github.com/sirupsen/logrus"
)

type rabbitMqSubscriber struct {
	client rabbitmq.RabbitMqClient
	config *config.ErrorMessageSubscriberConfig
}

func (subscriber *rabbitMqSubscriber) subscribe(stream chan data.HosepipeMessage) {
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
			log.WithError(err).Errorln("Error in subscribe method")
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

func NewrRabbitMqSubscriber(client *rabbitmq.RabbitMqClient, cfg *config.ErrorMessageSubscriberConfig) *rabbitMqSubscriber {
	return &rabbitMqSubscriber{client: client, config: cfg}
}
