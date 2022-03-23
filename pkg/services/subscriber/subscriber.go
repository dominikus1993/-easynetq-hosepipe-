package subscriber

import (
	"context"

	"github.com/dominikus1993/easynetq-hosepipe/pkg/config"
	"github.com/dominikus1993/easynetq-hosepipe/pkg/data"
	"github.com/dominikus1993/easynetq-hosepipe/pkg/rabbitmq"
)

type rabbitMqSubscriber struct {
	client rabbitmq.RabbitMqSubscriber
	config *config.ErrorMessageSubscriberConfig
}

func (f *rabbitMqSubscriber) Subscribe(c context.Context) <-chan *data.HosepipeMessage {

	return f.client.Subscribe(c, f.config., f.config.Queue, f.config.Topic)
}

func NewrRabbitMqSubscriber(client rabbitmq.RabbitMqSubscriber, cfg *config.ErrorMessageSubscriberConfig) *rabbitMqSubscriber {
	return &rabbitMqSubscriber{client: client, config: cfg}
}
