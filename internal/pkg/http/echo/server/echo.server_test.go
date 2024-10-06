package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/logger/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewEchoServer(t *testing.T) {
	log := &mocks.ILogger{}
	cfg := &EchoConfig{
		Host:               "localhost",
		Port:               "0",
		BasePath:           "/api/v1",
		DebugErrorResponse: false,
		Timeout:            time.Second * 10,
		CORSAllowOrigins:   []string{"*"},
		CORSAllowMethods:   []string{http.MethodGet, http.MethodPost},
		RateLimit:          5,
	}

	server := NewEchoServer(cfg, log)
	assert.NotNil(t, server)
}
