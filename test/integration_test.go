package test

import (
	"testing"

	publisher "github.com/Placons/oneapp-lib-servicebus-publisher/v2"
	testpublisher "github.com/Placons/oneapp-lib-servicebus-publisher/v2/integration"
	"github.com/Placons/oneapp-logger/logger"
	"github.com/stretchr/testify/assert"
)

func TestShouldPublishMessage(t *testing.T) {
	publisherMockServer := testpublisher.MockPublisher(t, "test-queue", "7000")
	defer publisherMockServer.Close()

	l := logger.NewStandardLogger("oneapp-my-service")

	c := publisher.ServiceBusConfig{
		Endpoint:      "sb://localhost:7000",
		QueueName:     "test-queue",
		SharedKeyName: "my-shared-key",
		SigningKey:    "my-signing-key",
	}
	p := publisher.NewPublisher(l, testpublisher.HttpClientNoCertVerify(), c)

	err := p.Publish("some message")
	assert.NoError(t, err)
}
func TestShouldPublishMessageWithProperties(t *testing.T) {
	props := map[string]string{"some-property": "some-value"}

	publisherMockServer := testpublisher.MockPublisherWithProps(t, "test-queue", "7000", props)
	defer publisherMockServer.Close()

	l := logger.NewStandardLogger("oneapp-my-service")

	c := publisher.ServiceBusConfig{
		Endpoint:      "sb://localhost:7000",
		QueueName:     "test-queue",
		SharedKeyName: "my-shared-key",
		SigningKey:    "my-signing-key",
	}
	p := publisher.NewPublisher(l, testpublisher.HttpClientNoCertVerify(), c)

	err := p.PublishWithProps("some message", props)
	assert.NoError(t, err)
}

func TestReturnErrorPublishMessageFails(t *testing.T) {
	publisherMockServer := testpublisher.MockPublisher(t, "test-queue", "7000")
	defer publisherMockServer.Close()

	l := logger.NewStandardLogger("oneapp-my-service")

	c := publisher.ServiceBusConfig{
		Endpoint:      "sb://localhost:7000",
		QueueName:     "some-unknown-queue",
		SharedKeyName: "my-shared-key",
		SigningKey:    "my-signing-key",
	}
	p := publisher.NewPublisher(l, testpublisher.HttpClientNoCertVerify(), c)

	err := p.Publish("some message")
	assert.Error(t, err)
}
