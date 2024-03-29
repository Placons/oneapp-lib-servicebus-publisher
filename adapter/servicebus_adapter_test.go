package adapter

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Placons/oneapp-logger/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var createdHTTPResponse = &http.Response{
	StatusCode: 201,
	Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
	Header:     make(http.Header),
}

var errorHTTPResponse = &http.Response{
	StatusCode: 500,
	Body:       ioutil.NopCloser(bytes.NewBufferString("Some error body")),
	Header:     make(http.Header),
}

func TestShouldSendMessageSuccessfully(t *testing.T) {
	mockClient := new(FakeHTTPClient)
	mockClient.On("Do", mock.Anything).Return(createdHTTPResponse, nil)

	adapter := NewServiceBusAdapter(logger.NewStandardLogger("test"), mockClient)

	err := adapter.SendMessage("some-url", "my-sas",
		DummyMessage{
			Greeting: "Hello from my go client!",
		}, nil)

	assert.NoError(t, err)
}

func TestReturnErrorWhenResponseStatusCodeNotCreated(t *testing.T) {
	mockClient := new(FakeHTTPClient)
	mockClient.On("Do", mock.Anything).Return(errorHTTPResponse, nil)

	adapter := NewServiceBusAdapter(logger.NewStandardLogger("test"), mockClient)

	err := adapter.SendMessage("some-url", "my-sas",
		DummyMessage{
			Greeting: "Hello from my go client!",
		}, nil)

	assert.Error(t, err)
	assert.EqualError(t, err, "(send-message) failed to send message to service bus due to statusCode: 500")
}

func TestServiceBusAdapterReturnErrorWhenHttpRequestCouldNotBeSent(t *testing.T) {
	mockClient := new(FakeHTTPClient)
	mockClient.On("Do", mock.Anything).Return(createdHTTPResponse, errors.New("An expected error"))

	adapter := NewServiceBusAdapter(logger.NewStandardLogger("test"), mockClient)

	err := adapter.SendMessage("some-url", "my-sas", nil, nil)

	assert.Error(t, err)
	assert.EqualError(t, err, "(send-message) failed to send http request: An expected error")
}

type FakeHTTPClient struct {
	mock.Mock
}

func (s *FakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := s.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

type DummyMessage struct {
	Greeting string
}
