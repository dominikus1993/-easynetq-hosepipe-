package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dominikus1993/easynetq-hosepipe/pkg/data"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

//var json = jsoniter.ConfigCompatibleWithStandardLibrary

type RabbitMqSubscriber interface {
	Subscribe(ctx context.Context, queue string) <-chan *data.HosepipeMessage
	Close()
}

type rabbitMqSubscriber struct {
	channel *amqp.Channel
}

func NewRabbitMqSubscriber(client RabbitMqClient) (*rabbitMqSubscriber, error) {
	channel, err := client.CreateChannel()
	if err != nil {
		return nil, fmt.Errorf("error when creating channel %w", err)
	}
	return &rabbitMqSubscriber{channel: channel}, nil
}

func (client *rabbitMqSubscriber) Close() {
	err := client.channel.Close()
	if err != nil {
		log.WithError(err).Fatalln("error when closing channel")
	}
}

func (client *rabbitMqSubscriber) Subscribe(ctx context.Context, queue string) <-chan *data.HosepipeMessage {
	res := make(chan *data.HosepipeMessage)

	q, err := client.channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.WithError(err).Fatalln("error when declaring queue")
	}

	go func(stream chan<- *data.HosepipeMessage) {
		defer close(stream)
		for {
			msg, ok, err := client.channel.Get(q.Name, true)
			if err != nil {
				log.WithError(err).Warnln("Error when tryig get message from rabbitmq")
				break
			}
			if !ok {
				break
			}
			var data data.HosepipeMessage
			err = json.Unmarshal(msg.Body, &data)
			if err != nil {
				log.WithError(err).Fatalln("error when unmarshaling message")
				continue
			}
			stream <- &data

		}
	}(res)
	return res
}
