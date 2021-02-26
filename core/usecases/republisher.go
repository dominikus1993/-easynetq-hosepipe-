package usecases

import (
	"context"
	"hossipe/core/model"
	"hossipe/core/services"
	"sync"

	log "github.com/sirupsen/logrus"
)

type RePublishErrorMessageUseCase interface {
	Start(c context.Context)
}

type rabbitmqRePublishErrorMessageUseCase struct {
	subscriber services.ErrorMessageSubscriber
	publisher  services.ErrorMessagePublisher
}

func (r *rabbitmqRePublishErrorMessageUseCase) publish(context context.Context, msg model.HosepipeMessage, wg *sync.WaitGroup) {
	defer wg.Done()
	err := r.publisher.RePublish(msg, context)
	if err != nil {
		log.WithError(err).Errorln("Error when trying re-publish messages")
	}
}

func (us *rabbitmqRePublishErrorMessageUseCase) Start(c context.Context) {
	defer us.subscriber.Close()
	defer us.publisher.Close()
	pubwg := &sync.WaitGroup{}

	consumerChannel := us.subscriber.Subscribe(c)

	for msg := range consumerChannel {
		pubwg.Add(1)
		go us.publish(c, msg, pubwg)
	}

	pubwg.Wait()
}
