package publisher

import (
	"errors"
	"testing"

	"github.com/Placons/oneapp-logger/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var message = map[string]string{
	"I am here": "to poke you",
}
var config = ServiceBusConfig{
	Endpoint:            "sb://my-namespace.servicebus.windows.net/",
	QueueName:           "my-queue",
	SharedKeyName:       "my-shared-keyname",
	SigningKey:          "my-signing-key",
	SigningKeyExpiresMS: 1234,
}

func TestShouldPublish(t *testing.T) {
	mockServiceBusAdapter := new(FakeServiceBusAdapter)
	mockGenerator := new(FakeSasGenerator)

	mockGenerator.On("Generate", "https://my-namespace.servicebus.windows.net/my-queue", "my-signing-key", 1234, "my-shared-keyname").Return("some-sas-token", nil)
	mockServiceBusAdapter.On("SendMessage", "https://my-namespace.servicebus.windows.net/my-queue", "some-sas-token", message, make(map[string]string)).Return(nil)

	publisher := Publisher{logger.NewStandardLogger("test"), mockServiceBusAdapter, mockGenerator, config}

	err := publisher.Publish(message)
	assert.NoError(t, err)

	mockGenerator.AssertExpectations(t)
	mockServiceBusAdapter.AssertExpectations(t)
}

func TestShouldPublishWithProperties(t *testing.T) {
	mockServiceBusAdapter := new(FakeServiceBusAdapter)
	mockGenerator := new(FakeSasGenerator)
	props := map[string]string{"some-key": "some-value", "another-key": "another-value"}
	mockGenerator.On("Generate", "https://my-namespace.servicebus.windows.net/my-queue", "my-signing-key", 1234, "my-shared-keyname").Return("some-sas-token", nil)
	mockServiceBusAdapter.On("SendMessage", "https://my-namespace.servicebus.windows.net/my-queue", "some-sas-token", message, props).Return(nil)

	publisher := Publisher{logger.NewStandardLogger("test"), mockServiceBusAdapter, mockGenerator, config}

	err := publisher.PublishWithProps(message, props)
	assert.NoError(t, err)

	mockGenerator.AssertExpectations(t)
	mockServiceBusAdapter.AssertExpectations(t)
}

func TestShouldReturnErrorWhenSendMessageReturns(t *testing.T) {
	mockServiceBusAdapter := new(FakeServiceBusAdapter)
	mockGenerator := new(FakeSasGenerator)

	mockGenerator.On("Generate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("some-sas-token", nil)
	mockServiceBusAdapter.On("SendMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("An expected error"))

	publisher := Publisher{logger.NewStandardLogger("test"), mockServiceBusAdapter, mockGenerator, config}

	err := publisher.Publish(message)
	assert.Error(t, err)

	mockGenerator.AssertExpectations(t)
	mockServiceBusAdapter.AssertExpectations(t)
}

type FakeServiceBusAdapter struct {
	mock.Mock
}

func (m *FakeServiceBusAdapter) SendMessage(url string, sasToken string, message interface{}, properties map[string]string) error {
	args := m.Called(url, sasToken, message, properties)
	return args.Error(0)
}

type FakeSasGenerator struct {
	mock.Mock
}

func (m *FakeSasGenerator) Generate(resourceURI string, signingKey string, expiresInMins int, policyName string) string {
	args := m.Called(resourceURI, signingKey, expiresInMins, policyName)
	return args.String(0)
}
