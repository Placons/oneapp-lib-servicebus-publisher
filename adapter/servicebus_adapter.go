package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Placons/oneapp-logger/logger"
)

const messageCreationFailure = "Failed to send message to service bus due to statusCode %d"

type ServiceBusAdapter struct {
	logger *logger.StandardLogger
	client HTTPClient
}

func NewServiceBusAdapter(logger *logger.StandardLogger, client HTTPClient) ServiceBusAdapter {
	return ServiceBusAdapter{
		logger: logger,
		client: client,
	}
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// SendMessage sends a message to an event bus queue using a POST http request
// serviceNamespace is the namespace of the azure service bus
// endpoint is the name of the endpoint (topic or queue)
// message is the actual message
func (a ServiceBusAdapter) SendMessage(baseURL string, sasToken string, message interface{}) error {
	url := fmt.Sprintf("%s/messages", baseURL)
	requestByte, _ := json.Marshal(message)
	r, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(requestByte))

	if err != nil {
		a.logger.ErrorWithErr("Failed to create http request", err)
		return err
	}

	r.Header.Add("Authorization", sasToken)

	resp, err := a.client.Do(r)
	if err != nil {
		a.logger.ErrorWithErr("Failed to send http request", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		err = fmt.Errorf(messageCreationFailure, resp.StatusCode)
		a.logger.ErrorWithErr("", err)
		return err
	}
	a.logger.DebugWithFields("Successfully sent message", map[string]interface{}{
		"message": string(requestByte),
		"url":     url,
	})
	return nil
}
