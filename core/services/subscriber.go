package services

import (
	"context"

	"github.com/dominikus1993/easynetq-hosepipe/core/model"
)

type ErrorMessageSubscriber interface {
	Subscribe(c context.Context) chan model.HosepipeMessage
	Close()
}
