package logger

import (
	"os"
	"testing"

	"github.com/NekKkMirror/go-app/internal/pkg/logger/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAppLogger_GetLevel_UnknownLevel(t *testing.T) {
	mockLogger := new(mocks.ILogger)
	mockLogger.On("GetLevel").Return(logrus.DebugLevel)

	appLogger := &appLogger{level: "unknown"}
	level := appLogger.GetLevel()

	assert.Equal(t, logrus.DebugLevel, level)
}

func TestAppLogger_GetLevel_KnownLevel(t *testing.T) {
	mockLogger := new(mocks.ILogger)
	mockLogger.On("GetLevel").Return(logrus.InfoLevel)

	appLogger := &appLogger{level: "info"}
	level := appLogger.GetLevel()

	assert.Equal(t, logrus.InfoLevel, level)
}

func TestInitLogger_ProductionEnvironment(t *testing.T) {
	cfg := &Config{LogLevel: "debug"}
	err := os.Setenv("APP_END", "production")
	if err != nil {
		return
	}

	logger := InitLogger(cfg)

	assert.NotNil(t, logger)
}

func TestInitLogger_DevelopmentEnvironment(t *testing.T) {
	cfg := &Config{LogLevel: "info"}

	err := os.Setenv("APP_END", "development")
	if err != nil {
		return
	}
	logger := InitLogger(cfg)

	assert.NotNil(t, logger)
}
