package services

import (
	"context"
	"easynetq-hossipe/core/model"
)

type ErrorMessageSubscriber interface {
	Subscribe(c context.Context) chan model.HosepipeMessage
}
