package services

import (
	"context"

	"github.com/dominikus1993/easynetq-hosepipe/core/model"
)

type ErrorMessagePublisher interface {
	RePublish(msg model.HosepipeMessage, c context.Context) error
	Close()
}
