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
	Namespace:           "my-name-space",
	Endpoint:            "my-queue",
	SharedKeyName:       "my-shared-keyname",
	SigningKey:          "my-signing-key",
	SigningKeyExpiresMS: 1234,
}

func TestShouldPublish(t *testing.T) {
	mockServiceBusAdapter := new(FakeServiceBusAdapter)
	mockGenerator := new(FakeSasGenerator)

	mockGenerator.On("Generate", "my-name-space.servicebus.windows.net/my-queue", "my-signing-key", 1234, "my-shared-keyname").Return("some-sas-token", nil)
	mockServiceBusAdapter.On("SendMessage", "my-name-space", "my-queue", "some-sas-token", message).Return(nil)

	publisher := NewPublisher(logger.NewStandardLogger("test"), mockServiceBusAdapter, config)
	publisher.generator = mockGenerator

	err := publisher.Publish(message)
	assert.NoError(t, err)

	mockGenerator.AssertExpectations(t)
	mockServiceBusAdapter.AssertExpectations(t)
}

func TestShouldReturnErrorWhenSendMessageReturns(t *testing.T) {
	mockServiceBusAdapter := new(FakeServiceBusAdapter)
	mockGenerator := new(FakeSasGenerator)

	mockGenerator.On("Generate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("some-sas-token", nil)
	mockServiceBusAdapter.On("SendMessage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("An expected error"))

	publisher := NewPublisher(logger.NewStandardLogger("test"), mockServiceBusAdapter, config)
	publisher.generator = mockGenerator

	err := publisher.Publish(message)
	assert.Error(t, err)

	mockGenerator.AssertExpectations(t)
	mockServiceBusAdapter.AssertExpectations(t)
}

type FakeServiceBusAdapter struct {
	mock.Mock
}

func (m *FakeServiceBusAdapter) SendMessage(serviceNamespace string, queue string, sasToken string, message interface{}) error {
	args := m.Called(serviceNamespace, queue, sasToken, message)
	return args.Error(0)
}

type FakeSasGenerator struct {
	mock.Mock
}

func (m *FakeSasGenerator) Generate(resourceURI string, signingKey string, expiresInMins int, policyName string) string {
	args := m.Called(resourceURI, signingKey, expiresInMins, policyName)
	return args.String(0)
}
