package data

import (
	"errors"
)

type MessageProperties struct {
	Headers       map[string]interface{} `json:"Headers,omitempty"`
	DeliveryMode  int                    `json:"DeliveryMode,omitempty"`
	CorrelationID string                 `json:"CorrelationId,omitempty"`
	Type          string                 `json:"Type,omitempty"`
}
type HosepipeMessage struct {
	RoutingKey      string            `json:"RoutingKey,omitempty"`
	Exchange        string            `json:"Exchange,omitempty"`
	Queue           string            `json:"Queue,omitempty"`
	Message         string            `json:"Message,omitempty"`
	Exception       string            `json:"Exception,omitempty"`
	DateTime        string            `json:"DateTime,omitempty"`
	BasicProperties MessageProperties `json:"BasicProperties,omitempty"`
}

func (hm *HosepipeMessage) GetMessageBody() ([]byte, error) {
	if hm.Message == "" {
		return nil, errors.New("Message is null")
	}
	return []byte(hm.Message), nil
}
