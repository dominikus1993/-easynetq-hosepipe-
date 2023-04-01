package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/dominikus1993/easynetq-hosepipe/pkg/data"
	log "github.com/sirupsen/logrus"
)

//var json = jsoniter.ConfigCompatibleWithStandardLibrary

type RabbitMqSubscriber interface {
	Subscribe(ctx context.Context, queue string) <-chan *data.HosepipeMessage
}

type rabbitMqSubscriber struct {
	client RabbitMqClient
}

func NewRabbitMqSubscriber(client RabbitMqClient) *rabbitMqSubscriber {
	return &rabbitMqSubscriber{client: client}
}

func (subscriber *rabbitMqSubscriber) Subscribe(ctx context.Context, queue string) <-chan *data.HosepipeMessage {
	res := make(chan *data.HosepipeMessage)
	go func(stream chan<- *data.HosepipeMessage) {
		defer close(stream)
		for {
			msg, ok, err := subscriber.client.Get(queue, true)
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
