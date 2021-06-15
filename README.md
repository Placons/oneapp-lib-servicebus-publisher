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
		Endpoint:            "sb://my-name-space.servicebus.windows.net",
		SharedKeyName:       "my-shared-keyname",
		SigningKey:          "my-signing-key",
		SigningKeyExpiresMS: 1234,
		QueueName:           "my-queue",
	}
	p := publisher.NewPublisher(l, serviceBusClient, config)
```

 and publish a message

```go
	p.Publish(map[string]string{"I am here": "to poke you"})
```

### Testing

Copy server.crt, server.key files from this library and put them in your project under a testdata/cert folder

In your service import testpublisher package

```testpublisher "github.com/Placons/oneapp-lib-servicebus-publisher/v2/integration"```

mock out the publisher by initializing an httptest.Server like this
```go
		publisherMockServer := testpublisher.MockPublisher(t, "test-queue", "7000")
    }
```

Use an insecure httpClient in order to bypass certificate verification
```go
	testpublisher.HttpClientNoCertVerify()
```