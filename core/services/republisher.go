package services

import (
	"context"
	"hossipe/core/model"
)

type ErrorMessagePublisher interface {
	RePublish(msgs chan model.HosepipeMessage, c context.Context)
	Close()
}
