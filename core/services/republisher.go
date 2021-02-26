package services

import (
	"context"
	"hossipe/core/model"
)

type ErrorMessagePublisher interface {
	RePublish(msg model.HosepipeMessage, c context.Context) error
	Close()
}
