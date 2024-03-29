package publisher

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/Placons/oneapp-lib-servicebus-publisher/v2/adapter"
	"github.com/Placons/oneapp-lib-servicebus-publisher/v2/sas"
	"github.com/Placons/oneapp-logger/logger"
)

type Publisher struct {
	logger            *logger.StandardLogger
	serviceBusAdapter ServiceBusAdapter
	generator         SasTokenGenerator
	config            ServiceBusConfig
}

//goland:noinspection GoUnusedExportedFunction
func NewPublisher(logger *logger.StandardLogger, client adapter.HTTPClient, config ServiceBusConfig) Publisher {
	return Publisher{
		logger:            logger,
		serviceBusAdapter: adapter.NewServiceBusAdapter(logger, client),
		generator:         sas.NewSasGenerator(realClock{}),
		config:            config,
	}
}

// Publish publishes the given message, using publisher's already provided configuration
func (p Publisher) Publish(message interface{}) error {
	return p.doPublish(message, make(map[string]string))
}

// Publish publishes the given message, including custom properties, using publisher's already provided configuration
func (p Publisher) PublishWithProps(message interface{}, properties map[string]string) error {
	return p.doPublish(message, properties)
}

func (p Publisher) doPublish(message interface{}, properties map[string]string) error {
	var (
		endpoint      = p.config.Endpoint
		queueName     = p.config.QueueName
		topicName     = p.config.TopicName
		signingKey    = p.config.SigningKey
		expiry        = p.config.SigningKeyExpiresMS
		sharedKeyName = p.config.SharedKeyName
	)
	l := p.logger.Audit("send-message")
	l.Start("sending message",
		logger.Generic("endpoint", endpoint),
		logger.Generic("queue", queueName),
		logger.Generic("topic", topicName),
		logger.Generic("messageLength", len(fmt.Sprintf("%v", message))),
	)

	publishURL, err := url.Parse(endpoint)
	if err != nil {
		p.logger.ErrorWithErrAndFields("Failed to construct publish url", err, map[string]interface{}{
			"endpoint": endpoint,
			"message":  message,
		})
		return err
	}
	// change schema from sb to https as the publisher uses rest whereas the consumer uses the custom sb-protocol
	publishURL.Scheme = "https"

	var pubURL string
	if queueName != "" {
		pubURL = joinURL(publishURL, queueName)
	} else if topicName != "" {
		pubURL = joinURL(publishURL, topicName)
	}
	if pubURL == "" {
		err = errors.New("pub url is not set")
		p.logger.ErrorWithErrAndFields("Could not construct pub url", err, map[string]interface{}{
			"endpoint": endpoint,
			"message":  message,
		})
		return err
	}
	sasToken := p.generator.Generate(pubURL, signingKey, expiry, sharedKeyName)

	err = p.serviceBusAdapter.SendMessage(pubURL, sasToken, message, properties)
	if err != nil {
		p.logger.ErrorWithErrAndFields("Failed to publish message to endpoint", err, map[string]interface{}{
			"endpoint": endpoint,
			"message":  message,
		})
		return err
	}
	l.End("sent message")
	return nil
}

func joinURL(url *url.URL, path string) string {
	up, err := url.Parse(path)
	if err != nil {
		return ""
	}
	return url.ResolveReference(up).String()
}

type ServiceBusAdapter interface {
	SendMessage(url string, sasToken string, message interface{}, properties map[string]string) error
}

type SasTokenGenerator interface {
	Generate(resourceUri string, signingKey string, expiresInMins int, policyName string) string
}

type ServiceBusConfig struct {
	Endpoint            string
	QueueName           string
	TopicName           string
	SharedKeyName       string
	SigningKey          string
	SigningKeyExpiresMS int
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }
