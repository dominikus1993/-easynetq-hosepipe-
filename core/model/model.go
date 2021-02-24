package model

import (
	"errors"
)

type HosepipeMessage struct {
	RoutingKey string `json:"RoutingKey,omitempty"`
	Exchange   string `json:"Exchange,omitempty"`
	Queue      string `json:"Queue,omitempty"`
	Message    string `json:"Message,omitempty"`
	Exception  string `json:"Exception,omitempty"`
	DateTime   string `json:"DateTime,omitempty"`
}

func (hm *HosepipeMessage) GetMessageBody() ([]byte, error) {
	if hm.Message == "" {
		return nil, errors.New("Message is null")
	}
	return []byte(hm.Message), nil
}
