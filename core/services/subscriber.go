package services

import (
	"context"
	"hossipe/core/model"
)

type ErrorMessageSubscriber interface {
	Subscribe(c context.Context) chan model.HosepipeMessage
	Close()
}
