package logger

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	lokicontainer "github.com/NekKkMirror/go-app/internal/pkg/container/test/loki"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoggerWithoutLoki(t *testing.T) {
	cfg := Config{
		LogLevel:      "info",
		LokiURL:       "",
		Labels:        map[string]string{"job": "test"},
		SendInterval:  time.Second,
		MaxBatchSize:  10,
		MaxBufferSize: 100,
		RetryCount:    3,
		RetryDelay:    500 * time.Millisecond,
	}

	testLogger := InitLogger(&cfg)

	appLogger, ok := testLogger.(*appLogger)
	assert.True(t, ok)

	logEntry := appLogger.instance.WithFields(log.Fields{
		"serviceName": "testService",
		"moduleName":  "testModule",
	})
	logEntry.Info("Test log message")

	assert.NotNil(t, testLogger)
	assert.Equal(t, log.InfoLevel, appLogger.instance.GetLevel())
}

func TestLoggerWithLoki(t *testing.T) {
	ctx := context.Background()
	lokiURL, err := lokicontainer.Start(ctx, t)
	assert.NoError(t, err)

	cfg := Config{
		LogLevel:      "info",
		LokiURL:       lokiURL,
		Labels:        map[string]string{"job": "test", "service": "testService"},
		SendInterval:  time.Second,
		MaxBatchSize:  10,
		MaxBufferSize: 100,
		RetryCount:    3,
		RetryDelay:    500 * time.Millisecond,
	}

	testLogger := InitLogger(&cfg)

	appLogger, ok := testLogger.(*appLogger)
	assert.True(t, ok)

	logEntry := log.NewEntry(appLogger.instance).WithFields(log.Fields{
		"serviceName": "testService",
		"moduleName":  "testModule",
	})
	logEntry.Message = "Test log message"
	hook := NewLokiHook(&cfg)

	err = hook.Fire(logEntry)
	assert.NoError(t, err)

	time.Sleep(20 * time.Second)

	query := `{job="test",service="testService"}`
	start := time.Now().Add(-5 * time.Minute)
	end := time.Now()

	respBody, err := hook.QueryRange(query, 10, start, end)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	assert.NoError(t, err)

	data, ok := result["data"].(map[string]interface{})
	assert.True(t, ok)

	streams, ok := data["result"].([]interface{})
	assert.True(t, ok)
	assert.Greater(t, len(streams), 0)

	assert.Contains(t, string(respBody), "Test log message")
	assert.Contains(t, string(respBody), "testModule")
}
