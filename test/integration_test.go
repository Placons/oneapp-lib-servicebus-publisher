package test

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	publisher "github.com/Placons/oneapp-lib-servicebus-publisher/v2"
	"github.com/Placons/oneapp-logger/logger"
	"github.com/stretchr/testify/assert"
)

var sa []*httptest.Server
var (
	requests      map[string]int
	requestsMutex = sync.RWMutex{}
)

func TestShouldPublishMessage(t *testing.T) {
	publisherMockServer := MockPublisher(t, "test-queue")
	defer publisherMockServer.Close()

	l := logger.NewStandardLogger("oneapp-my-service")

	c := publisher.ServiceBusConfig{
		Endpoint:      "sb://localhost:7000",
		QueueName:     "test-queue",
		SharedKeyName: "my-shared-key",
		SigningKey:    "my-signing-key",
	}
	p := publisher.NewPublisher(l, HttpClientNoCertVerify(), c)

	err := p.Publish("some message")
	assert.NoError(t, err)
}

func TestReturnErrorPublishMessageFails(t *testing.T) {
	publisherMockServer := MockPublisher(t, "test-queue")
	defer publisherMockServer.Close()

	l := logger.NewStandardLogger("oneapp-my-service")

	c := publisher.ServiceBusConfig{
		Endpoint:      "sb://localhost:7000",
		QueueName:     "some-unknown-queue",
		SharedKeyName: "my-shared-key",
		SigningKey:    "my-signing-key",
	}
	p := publisher.NewPublisher(l, HttpClientNoCertVerify(), c)

	err := p.Publish("some message")
	assert.Error(t, err)
}

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

func HttpClientNoCertVerify() *http.Client {
	return &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}
}

func createListener(port string) net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Fatal(err)
	}
	return l
}
