package config

type ErrorMessageSubscriberConfig struct {
	Queue string
}

func NewErrorMessageSubscriberConfig(queue string) *ErrorMessageSubscriberConfig {
	return &ErrorMessageSubscriberConfig{Queue: queue}
}
