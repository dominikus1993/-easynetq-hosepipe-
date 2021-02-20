package model

type HosepipeMessage struct {
	RoutingKey string
	Exchange   string
	Queue      string
	Message    string
}
