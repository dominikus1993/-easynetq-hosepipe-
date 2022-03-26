package republisher

import (
	"context"

	"github.com/dominikus1993/easynetq-hosepipe/pkg/data"
	"github.com/dominikus1993/easynetq-hosepipe/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type rabbitMqPublisher struct {
	client rabbitmq.RabbitMqPublisher
}

func (f *rabbitMqPublisher) RePublish(ctx context.Context, msg *data.HosepipeMessage) error {
	body, err := msg.GetMessageBody()
	if err != nil {
		fields := log.Fields{"Exchange": msg.Exchange, "Queue": msg.Queue, "Topic": msg.RoutingKey, "Exception": msg.Exception}
		log.WithFields(fields).WithError(err).WithContext(ctx).Errorln("Error when trying resend message")
		return err
	}

	props := msg.BasicProperties

	shouldRepublish, newHeaders := checkIfShouldRepublishMessage(props.Headers)
	if shouldRepublish {
		log.WithContext(ctx).Traceln("Message published")
		rMsg := amqp.Publishing{ContentType: "application/json", CorrelationId: props.CorrelationID, Type: props.Type, Headers: newHeaders, Body: body}
		return f.client.Publish(ctx, msg.Exchange, msg.RoutingKey, rMsg)
	}
	return nil
}

func checkIfShouldRepublishMessage(headers map[string]interface{}) (bool, map[string]interface{}) {
	if val, ok := headers["retry"].(int); ok {
		if val < 10 {
			headers["retry"] = val + 1
			return true, headers
		}
		return false, headers
	}
	headers["retry"] = 1
	return true, headers
}

func NewrRabbitMqPublisher(client rabbitmq.RabbitMqPublisher) *rabbitMqPublisher {
	return &rabbitMqPublisher{client: client}
}
