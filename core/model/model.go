package model

type HosepipeMessage struct {
	RoutingKey string `json:"RoutingKey,omitempty"`
	Exchange   string `json:"Exchange,omitempty"`
	Queue      string `json:"Queue,omitempty"`
	Message    string `json:"Message,omitempty"`
}
