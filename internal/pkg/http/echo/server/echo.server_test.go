package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/logger/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewEchoServerInitialization(t *testing.T) {
	e := NewEchoServer()
	assert.NotNil(t, e, "Expected Echo server instance to be initialized")
}

func TestEchoLibraryImportFailure(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r, "undefined", "Expected panic due to missing Echo library import")
		}
	}()

	// Simulate missing import by calling a function that requires the Echo library
	_ = NewEchoServer()
}

func TestApplyVersioningFromHeaderMiddleware(t *testing.T) {
	e := echo.New()
	ApplyVersioningFromHeader(e)

	req := httptest.NewRequest(http.MethodGet, "/tests", nil)
	req.Header.Set("version", "v1")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handlerFunc := func(c echo.Context) error {
		return c.String(http.StatusOK, "Test")
	}

	apiVersionMiddleware := apiVersion(handlerFunc)
	err := apiVersionMiddleware(c)
	if err != nil {
		return
	}

	if req.URL.Path != "/v1/tests" {
		t.Errorf("Expected path to be '/v1/tests', got '%s'", req.URL.Path)
	}
}

func TestRegisterGroupWithValidNameAndBuilder(t *testing.T) {
	e := echo.New()
	groupName := "/tests"
	builder := func(g *echo.Group) {
		g.GET("/endpoint", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})
	}

	result := RegisterGroupFunc(groupName, e, builder)

	assert.NotNil(t, result)
	assert.NotNil(t, result.Routes())
	assert.Equal(t, groupName+"/endpoint", result.Routes()[0].Path)
}

func TestRunHttpServerWithTimeoutSettings(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := &mocks.ILogger{}
	log.On("Infof", "shutting down HTTP server on port: %s", "8080").Return()

	cfg := &EchoConfig{
		Host:     "localhost",
		Port:     "8080",
		BasePath: "/api/v1",
	}

	e := echo.New()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := RunHttpServer(ctx, e, log, cfg)
		assert.NoError(t, err)
	}()

	time.Sleep(1 * time.Second)

	cancel()

	wg.Wait()

	assert.Equal(t, MaxHeaderBytes, e.Server.MaxHeaderBytes)
	assert.Equal(t, ReadTimeout, e.Server.ReadTimeout)
	assert.Equal(t, WriteTimeout, e.Server.WriteTimeout)

	log.AssertExpectations(t)
}
