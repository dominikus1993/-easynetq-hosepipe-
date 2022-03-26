package main

import (
	"context"
	"flag"

	"github.com/dominikus1993/easynetq-hosepipe/pkg/rabbitmq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var amqpConnection string

func init() {
	flag.String("rabbitmq", "amqp://guest:guest@localhost:5672/", "help message for flagname")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	amqpConnection = viper.GetString("rabbitmq")
}

func main() {
	client, err := rabbitmq.NewRabbitMqClient(amqpConnection)
	if err != nil {
		log.WithError(err).Panicln("error when creating client")
	}
	defer client.Close()
	subscriber, err := rabbitmq.NewRabbitMqSubscriber(client)
	if err != nil {
		log.WithError(err).Panicln("error when creating subscriber")
	}
	defer subscriber.CloseChannel()
	republiser, err := rabbitmq.NewRabbitMqPublisher(client)
	if err != nil {
		log.WithError(err).Panicln("error when creating publisher")
	}
	defer republiser.CloseChannel()

	ctx := context.Background()
	for rabbitError := range subscriber.Subscribe(ctx, "easynetq-hosepipe") {
		republiser.Publish(ctx, "easynetq-hosepipe", rabbitError.Topic, rabbitError.Message)
	}
}

func modify(headers map[string]interface{}) {
	headers["21"] = 2112
}
