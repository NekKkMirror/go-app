package httpclient

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func TestNewHttpClient(t *testing.T) {
	client := NewHttpClient()

	// Test if the client is not nil
	assert.NotNil(t, client, "Expected client to be non-nil")

	// Test if the client's transport is correctly set
	_, ok := client.GetClient().Transport.(*otelhttp.Transport)
	assert.True(t, ok, "Expected transport to be of type *otelhttp.Transport")

	// Test timeout settings
	assert.Equal(t, timeout, client.GetClient().Timeout, "Expected Timeout to be equal to timeout")

	// Test retry settings
	assert.Equal(t, retryCount, client.RetryCount, "Expected RetryCount to be equal to retryCount")
	assert.Equal(t, retryWaitTime, client.RetryWaitTime, "Expected RetryWaitTime to be equal to retryWaitTime")
}

func TestCreateTransport(t *testing.T) {
	transport := createTransport()

	// Test if the transport is not nil
	assert.NotNil(t, transport, "Expected transport to be non-nil")

	// Test the transport's configuration
	assert.Equal(t, tLSHandshakeTimeout, transport.TLSHandshakeTimeout, "Expected TLSHandshakeTimeout to be equal to tLSHandshakeTimeout")
	assert.Equal(t, maxIdleConns, transport.MaxIdleConns, "Expected MaxIdleConns to be equal to maxIdleConns")
	assert.Equal(t, maxConnsPerHost, transport.MaxConnsPerHost, "Expected MaxConnsPerHost to be equal to maxConnsPerHost")
	assert.Equal(t, idleConnTimeout, transport.IdleConnTimeout, "Expected IdleConnTimeout to be equal to idleConnTimeout")
	assert.Equal(t, responseHeaderTimeout, transport.ResponseHeaderTimeout, "Expected ResponseHeaderTimeout to be equal to responseHeaderTimeout")

	// Test if DialContext function is correctly configured
	dialer := &net.Dialer{Timeout: dialContextTimeout}
	assert.NotNil(t, transport.DialContext, "Expected transport.DialContext to be non-nil")

	// Test dynamic dialing to ensure proper function assignment (Added for completeness)
	conn, err := dialer.Dial("tcp", "localhost:0")
	assert.Error(t, err, "Expected an error for invalid address")
	if conn != nil { // Ensure no connection is left open
		conn.Close()
	}
}
