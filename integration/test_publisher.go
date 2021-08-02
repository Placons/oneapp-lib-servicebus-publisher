package integration

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockPublisher(t *testing.T, queue string, port string) *httptest.Server {
	return mockPublisher(t, queue, port, http.StatusCreated)
}

func mockPublisher(t *testing.T, queue string, port string, s int) *httptest.Server {
	publisherServeMux := http.NewServeMux()
	publisherServeMux.HandleFunc(fmt.Sprintf("/%s/messages", queue), func(res http.ResponseWriter, req *http.Request) {
		assert.Equal(t, http.MethodPost, req.Method)
		res.WriteHeader(s)
	})

	cert, err := tls.LoadX509KeyPair("testdata/cert/server.crt", "testdata/cert/server.key")
	if err != nil {
		t.Errorf("Bad server certs, error was: %s", err)
	}
	certs := []tls.Certificate{cert}

	publisherMockServer := httptest.NewUnstartedServer(publisherServeMux)
	publisherMockServer.TLS = &tls.Config{Certificates: certs}
	publisherMockServer.Listener.Close()
	publisherMockServer.Listener = createListener(port)
	publisherMockServer.StartTLS()

	return publisherMockServer
}

func BrokenMockPublisher(t *testing.T, queue string, port string) *httptest.Server {
	return mockPublisher(t, queue, port, http.StatusInternalServerError)
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
