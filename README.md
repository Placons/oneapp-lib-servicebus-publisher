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

In your service mock out the publisher by initializing an httptest.Server like this
```go
	func MockPublisher(t *testing.T, queue string) *httptest.Server {
		publisherServeMux := http.NewServeMux()
		publisherServeMux.HandleFunc(fmt.Sprintf("/%s/messages", queue), func(res http.ResponseWriter, req *http.Request) {
			assert.Equal(t, http.MethodPost, req.Method)
			res.WriteHeader(http.StatusCreated)
		})

		cert, err := tls.LoadX509KeyPair("testdata/cert/server.crt", "testdata/cert/server.key")
		if err != nil {
			t.Errorf("Bad server certs, error was: %s", err)
		}
		certs := []tls.Certificate{cert}

		publisherMockServer := httptest.NewUnstartedServer(publisherServeMux)
		publisherMockServer.TLS = &tls.Config{Certificates: certs}
		publisherMockServer.Listener.Close()
		publisherMockServer.Listener = createListener("7000")
		publisherMockServer.StartTLS()

		return publisherMockServer
    }
```

Use an insecure httpClient in order to bypass certificate verification
```go
	func HttpClientNoCertVerify() *http.Client {
		return &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}
	}
```