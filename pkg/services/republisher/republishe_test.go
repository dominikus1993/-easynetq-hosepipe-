package republisher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	headers         RabbitMqHeaders
	shouldRepublish bool
	newHeaders      RabbitMqHeaders
}

func TestCheckIfShouldRepublishMessage(t *testing.T) {
	testData := []testData{
		{headers: RabbitMqHeaders{"retry": 1}, shouldRepublish: true, newHeaders: RabbitMqHeaders{"retry": 2}},
		{headers: RabbitMqHeaders{"retry": 10}, shouldRepublish: false, newHeaders: RabbitMqHeaders{"retry": 10}},
		{headers: RabbitMqHeaders{}, shouldRepublish: true, newHeaders: RabbitMqHeaders{"retry": 1}},
	}
	for _, tt := range testData {
		t.Run("", func(t *testing.T) {
			shouldRepublish, newHeaders := checkIfShouldRepublishMessage(tt.headers)
			assert.Equal(t, tt.shouldRepublish, shouldRepublish)
			assert.Equal(t, tt.newHeaders, newHeaders)
		})
	}
}
