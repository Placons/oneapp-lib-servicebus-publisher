package integration

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockPublisher(t *testing.T, queue string, port string) *httptest.Server {
	return mockPublisher(t, queue, port, http.StatusCreated, make(map[string]string), nil)
}

func MockPublisherWithProps(t *testing.T, queue string, port string, properties map[string]string) *httptest.Server {
	return mockPublisher(t, queue, port, http.StatusCreated, properties, nil)
}

func MockPublisherWithPropsAndBody(t *testing.T, queue string, port string, properties map[string]string, body *string) *httptest.Server {
	return mockPublisher(t, queue, port, http.StatusCreated, properties, body)
}

func BrokenMockPublisher(t *testing.T, queue string, port string) *httptest.Server {
	return mockPublisher(t, queue, port, http.StatusInternalServerError, make(map[string]string), nil)
}

func mockPublisher(t *testing.T, queue string, port string, s int, properties map[string]string, body *string) *httptest.Server {
	publisherServeMux := http.NewServeMux()
	publisherServeMux.HandleFunc(fmt.Sprintf("/%s/messages", queue), func(res http.ResponseWriter, req *http.Request) {
		if body != nil {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Errorf("could not read request: %v", err)
			}
			assert.JSONEq(t, *body, string(b))
		}

		assertProperties(t, req, properties)
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

func assertProperties(t *testing.T, req *http.Request, properties map[string]string) {
	for h, v := range properties {
		got := req.Header.Get(h)
		if len(got) == 0 {
			t.Errorf("Property %s was not found in request header", h)
			return
		}

		assert.Equal(t, v, got, "Unexpected value found in request header")
	}
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
