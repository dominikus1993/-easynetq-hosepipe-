package model

import "github.com/streadway/amqp"

func a(rabbitmq *amqp.Connection) {
	ch, _ := rabbitmq.Channel()

	msgs, _ := ch.Consume(
		"q.Name",  // queue
		"crawler", // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	for msg := range msgs {
		msg.
	}

}

type HossipeMessage struct {
	Body []byte
}
