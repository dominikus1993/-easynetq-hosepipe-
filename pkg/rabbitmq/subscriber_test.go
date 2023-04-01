package rabbitmq

import (
	"context"
	"fmt"
	"testing"

	"github.com/dominikus1993/integrationtestcontainers-go/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

var event string = `{
	"RoutingKey":"#",
	"Exchange":"notifications",
	"Queue":"sender",
	"Exception":"System.Exception",
	"Message":"Test",
	"DateTime":"2022-08-09T14:23:08.0690627Z",
	"BasicProperties":{
	   "ContentType":null,
	   "ContentEncoding":null,
	   "Headers":{
		  "x-datadog-sampling-priority":"0"
	   },
	   "DeliveryMode":0,
	   "Priority":0,
	   "CorrelationId":"14a6caf0-68f1-42de-b0b6-461a0dca4473",
	   "ReplyTo":null,
	   "Expiration":null,
	   "MessageId":null,
	   "Timestamp":0,
	   "Type":"Some.NameSpace, Messaging.Messages",
	   "UserId":null,
	   "AppId":null,
	   "ClusterId":null,
	   "ContentTypePresent":false,
	   "ContentEncodingPresent":false,
	   "HeadersPresent":true,
	   "DeliveryModePresent":false,
	   "PriorityPresent":false,
	   "CorrelationIdPresent":true,
	   "ReplyToPresent":false,
	   "ExpirationPresent":false,
	   "MessageIdPresent":false,
	   "TimestampPresent":false,
	   "TypePresent":true,
	   "UserIdPresent":false,
	   "AppIdPresent":false,
	   "ClusterIdPresent":false
	}
 }`

func TestSubscriber(t *testing.T) {
	exchangeName := "test-2137"
	queueName := "test-69420"
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	config := rabbitmq.NewRabbitMqContainerConfigurationBuilder().Build()
	// Arrange
	ctx := context.Background()

	container, err := rabbitmq.StartContainer(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	url, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("can't download rabbit conectionstring, %w", err))
	}
	conn, err := NewRabbitMqClient(url)
	if err != nil {
		t.Fatal(fmt.Errorf("can't connect to rabbitmq, %w", err))
	}

	t.Cleanup(func() {
		conn.Close()
	})

	subscriber := NewRabbitMqSubscriber(conn)

	publisher := NewRabbitMqPublisher(conn)

	err = conn.DeclareExchange(ctx, exchangeName)
	assert.NoError(t, err)
	err = conn.DeclareQueue(ctx, queueName)
	assert.NoError(t, err)
	err = conn.QueueBind(exchangeName, queueName, "test")
	assert.NoError(t, err)
	err = publisher.Publish(ctx, exchangeName, "test", amqp091.Publishing{Body: []byte(event)})
	assert.NoError(t, err)

	result := subscriber.Subscribe(ctx, queueName)

	subject := ToSlice(result)

	assert.NotEmpty(t, subject)

}

func ToSlice[T any](s <-chan T) []T {
	res := make([]T, 0)
	for v := range s {
		res = append(res, v)
	}
	return res
}
