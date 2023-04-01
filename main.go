package main

import (
	"context"

	"github.com/dominikus1993/easynetq-hosepipe/pkg/rabbitmq"
	"github.com/dominikus1993/easynetq-hosepipe/pkg/services/republisher"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var configuration HosepipeConfig

type HosepipeConfig struct {
	AmqpConnection string `mapstructure:"AMQP_CONNECTION"`
	ErrorQueue     string `mapstructure:"ERROR_QUEUE"`
}

func LoadConfig(path string) (config HosepipeConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func init() {
	cfg, err := LoadConfig(".")
	if err != nil {
		log.WithError(err).Fatalln("error when loading config")
	}
	configuration = cfg
}

func main() {
	client, err := rabbitmq.NewRabbitMqClient(configuration.AmqpConnection)
	if err != nil {
		log.WithError(err).Panicln("error when creating client")
	}
	defer client.Close()
	subscriber, err := rabbitmq.NewRabbitMqSubscriber(client)
	if err != nil {
		log.WithError(err).Panicln("error when creating subscriber")
	}
	defer subscriber.Close()
	publisher, err := rabbitmq.NewRabbitMqPublisher(client)
	if err != nil {
		log.WithError(err).Panicln("error when creating publisher")
	}
	defer publisher.CloseChannel()

	rep := republisher.NewrRabbitMqPublisher(publisher)

	ctx := context.Background()
	StartRepublish(ctx, subscriber, rep)
}

func StartRepublish(ctx context.Context, subscriber rabbitmq.RabbitMqSubscriber, republisher republisher.RabbitMqRePublisher) {
	log.Infoln("Start")
	for rabbitError := range subscriber.Subscribe(ctx, configuration.ErrorQueue) {
		err := republisher.RePublish(ctx, rabbitError)
		if err != nil {
			log.WithError(err).Panicln("error when republishing")
		}
	}
}
