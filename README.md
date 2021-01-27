# oneapp-lib-servicebus-publisher

This is a library to publish a message to an azure service endpoint (topic or queue).

### Usage

Import oneapp-lib-servicebus-publisher

```import ("github.com/Placons/oneapp-lib-servicebus-publisher/publisher")```

or

```go get github.com/Placons/oneapp-lib-servicebus-publisher/publisher```

From your service initialize a publisher with all the dependant services

```go
	l := logger.NewStandardLogger("oneapp-my-service")
	serviceBusClient := &http.Client{Timeout: time.Duration(1000) * time.Millisecond}

	// these values will be derived from the services configuration
	config := publisher.ServiceBusConfig{
		Namespace:           "my-name-space",
		Endpoint:            "my-queue",
		SharedKeyName:       "my-shared-keyname",
		SigningKey:          "my-signing-key",
		SigningKeyExpiresMS: 1234,
		EndpointURL:         "https://my-name-space.servicebus.windows.net/my-queue"
	}
	p := publisher.NewPublisher(l, serviceBusClient, config)
```

 and publish a message

```go
	p.Publish(map[string]string{"I am here": "to poke you"})
```
