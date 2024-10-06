package httpclient

import (
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	timeout               = 5 * time.Second
	dialContextTimeout    = 5 * time.Second
	tLSHandshakeTimeout   = 5 * time.Second
	maxIdleConns          = 20
	maxConnsPerHost       = 40
	retryCount            = 3
	retryWaitTime         = 300 * time.Millisecond
	idleConnTimeout       = 120 * time.Second
	responseHeaderTimeout = 5 * time.Second
)

// NewHttpClient creates and configures a new Resty HTTP client.
func NewHttpClient() *resty.Client {
	transport := createTransport()

	client := resty.New().
		SetTimeout(timeout).
		SetRetryCount(retryCount).
		SetRetryWaitTime(retryWaitTime).
		SetTransport(otelhttp.NewTransport(transport))

	return client
}

// createTransport configures and returns a new http.Transport instance.
func createTransport() *http.Transport {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: dialContextTimeout,
		}).DialContext,
		TLSHandshakeTimeout:   tLSHandshakeTimeout,
		MaxIdleConns:          maxIdleConns,
		MaxConnsPerHost:       maxConnsPerHost,
		IdleConnTimeout:       idleConnTimeout,
		ResponseHeaderTimeout: responseHeaderTimeout,
	}
}
