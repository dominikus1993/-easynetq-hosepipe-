package services

import (
	"context"
	"encoding/json"
	"hossipe/core/model"
	"hossipe/core/services"
	"hossipe/infrastructure/rabbitmq"
	"log"

	"github.com/streadway/amqp"
)

type ErrorMessageSubscriberConfig struct {
	Queue string
}

type rabbitMqSubscriber struct {
	client *rabbitmq.RabbitMqClient
	config *ErrorMessageSubscriberConfig
}

func subscribe(rabbitmq *amqp.Connection, stream chan model.HosepipeMessage) {
	const exchange = "crawl-media"
	ch, err := rabbitmq.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"crawl-media", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		log.Fatal(err)
	}

	q, err := ch.QueueDeclare(
		exchange, // name
		true,     // durable
		false,    // delete when usused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(
		q.Name,   // queue name
		"#",      // routing key
		exchange, // exchange
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name,    // queue
		"crawler", // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	for msg := range msgs {
		var res model.CrawlWebsite
		err := json.Unmarshal(msg.Body, &res)
		if err != nil {
			log.Println(err)
		} else {
			stream <- res
		}
	}

	close(stream)
}

func (f *rabbitMqSubscriber) Subscribe(c context.Context) chan model.HosepipeMessage {
	stream := make(chan model.CrawlWebsite)

	go subscribe(f.rabbitmq, stream)

	return stream
}

func NewrRabbitMqSubscriber(client *rabbitmq.RabbitMqClient, cfg *ErrorMessageSubscriberConfig) services.ErrorMessageSubscriber {
	return &rabbitMqSubscriber{client: client, config: cfg}
}
