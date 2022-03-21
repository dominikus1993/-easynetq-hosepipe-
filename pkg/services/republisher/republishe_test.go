package republisher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	headers         map[string]interface{}
	shouldRepublish bool
	newHeaders      map[string]interface{}
}

func TestCheckIfShouldRepublishMessage(t *testing.T) {
	testData := []testData{
		{headers: map[string]interface{}{"retry": 1}, shouldRepublish: true, newHeaders: map[string]interface{}{"retry": 2}},
		{headers: map[string]interface{}{"retry": 10}, shouldRepublish: false, newHeaders: map[string]interface{}{"retry": 10}},
		{headers: map[string]interface{}{}, shouldRepublish: true, newHeaders: map[string]interface{}{"retry": 1}},
	}
	for _, tt := range testData {
		t.Run("", func(t *testing.T) {
			shouldRepublish, newHeaders := checkIfShouldRepublishMessage(tt.headers)
			assert.Equal(t, tt.shouldRepublish, shouldRepublish)
			assert.Equal(t, tt.newHeaders, newHeaders)
		})
	}
}
