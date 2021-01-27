package publisher

import (
	"fmt"
	"time"

	"github.com/Placons/oneapp-lib-servicebus-publisher/adapter"
	"github.com/Placons/oneapp-lib-servicebus-publisher/sas"
	"github.com/Placons/oneapp-logger/logger"
)

type Publisher struct {
	logger            *logger.StandardLogger
	serviceBusAdapter ServiceBusAdapter
	generator         SasTokenGenerator
	config            ServiceBusConfig
}

func NewPublisher(logger *logger.StandardLogger, client adapter.HTTPClient, config ServiceBusConfig) Publisher {
	return Publisher{
		logger:            logger,
		serviceBusAdapter: adapter.NewServiceBusAdapter(logger, client),
		generator:         sas.NewSasGenerator(realClock{}),
		config:            config,
	}
}

func (p Publisher) Publish(message interface{}) error {
	var (
		baseURL       = p.config.EndpointBaseURL
		namespace     = p.config.Namespace
		endpoint      = p.config.Endpoint
		signingKey    = p.config.SigningKey
		expiry        = p.config.SigningKeyExpiresMS
		sharedKeyName = p.config.SharedKeyName
	)

	sasToken := p.generator.Generate(fmt.Sprintf("%s.servicebus.windows.net/%s", namespace, endpoint), signingKey, expiry, sharedKeyName)

	err := p.serviceBusAdapter.SendMessage(baseURL, sasToken, message)
	if err != nil {
		p.logger.ErrorWithErrAndFields("Failed to publish message to endpoint", err, map[string]interface{}{
			"endpoint": endpoint,
			"message":  message,
		})
		return err
	}
	p.logger.DebugWithFields("Successfully publish message to endpoint", map[string]interface{}{
		"endpoint": endpoint,
		"message":  message,
	})
	return nil
}

type ServiceBusAdapter interface {
	SendMessage(url string, sasToken string, message interface{}) error
}

type SasTokenGenerator interface {
	Generate(resourceUri string, signingKey string, expiresInMins int, policyName string) string
}

type ServiceBusConfig struct {
	EndpointBaseURL     string
	Namespace           string
	Endpoint            string
	SharedKeyName       string
	SigningKey          string
	SigningKeyExpiresMS int
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }
