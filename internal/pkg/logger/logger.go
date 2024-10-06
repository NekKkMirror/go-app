package logger

import (
	"sync"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/constants/environment"
	"github.com/NekKkMirror/go-app/internal/pkg/utils/common"
	log "github.com/sirupsen/logrus"
)

// Config holds the configuration parameters for initializing the logger.
type Config struct {
	LogLevel string
	LokiURL  string
	Labels   map[string]string

	// Parameters for Loki should be moved to lokihook.go
	SendInterval  time.Duration
	MaxBatchSize  int
	MaxBufferSize int
	RetryCount    int
	RetryDelay    time.Duration
}

// appLogger wraps around logrus.Logger providing additional functionality or interfaces if needed.
type appLogger struct {
	instance *log.Logger
}

// Singleton instance of the logger
var (
	Logger ILogger
	once   sync.Once
)

// InitLogger initializes the logger singleton based on the provided configuration.
func InitLogger(cfg *Config) ILogger {
	once.Do(func() {
		loggerInstance := log.New()
		loggerInstance.SetFormatter(getFormatter())
		loggerInstance.SetLevel(parseLogLevel(cfg.LogLevel))

		l := &appLogger{instance: loggerInstance}

		if cfg.LokiURL != "" {
			hook := NewLokiHook(cfg)
			loggerInstance.AddHook(hook)
		}

		Logger = l
	})
	return Logger
}

// getFormatter returns the appropriate log formatter depending on the environment.
func getFormatter() log.Formatter {
	if common.GetEnv("APP_ENV", environment.Development) == environment.Production {
		return &log.JSONFormatter{}
	}
	return &log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	}
}

// parseLogLevel converts a string representation of a log level into a log.Level.
func parseLogLevel(level string) log.Level {
	switch level {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}

// Methods for use with the logger

// Debug logs a message at level Debug.
func (l *appLogger) Debug(args ...interface{}) {
	l.instance.Debug(args...)
}

// Debugf logs a message at level Debug using the provided format and arguments.
func (l *appLogger) Debugf(format string, args ...interface{}) {
	l.instance.Debugf(format, args...)
}

// Info logs a message at level Info.
func (l *appLogger) Info(args ...interface{}) {
	l.instance.Info(args...)
}

// Infof logs a message at level Info using the provided format and arguments.
func (l *appLogger) Infof(format string, args ...interface{}) {
	l.instance.Infof(format, args...)
}

// Warn logs a message at level Warn.
func (l *appLogger) Warn(args ...interface{}) {
	l.instance.Warn(args...)
}

// Warnf logs a message at level Warn using the provided format and arguments.
func (l *appLogger) Warnf(format string, args ...interface{}) {
	l.instance.Warnf(format, args...)
}

// Error logs a message at level Error.
func (l *appLogger) Error(args ...interface{}) {
	l.instance.Error(args...)
}

// Errorf logs a message at level Error using the provided format and arguments.
func (l *appLogger) Errorf(format string, args ...interface{}) {
	l.instance.Errorf(format, args...)
}

// Fatal logs a message at level Fatal and then the process will exit with status set to 1.
func (l *appLogger) Fatal(args ...interface{}) {
	l.instance.Fatal(args...)
}

// Fatalf logs a message at level Fatal using the provided format and arguments, then the process will exit with status set to 1.
func (l *appLogger) Fatalf(format string, args ...interface{}) {
	l.instance.Fatalf(format, args...)
}

// Trace logs a message at level Trace.
func (l *appLogger) Trace(args ...interface{}) {
	l.instance.Trace(args...)
}

// Tracef logs a message at level Trace using the provided format and arguments.
func (l *appLogger) Tracef(format string, args ...interface{}) {
	l.instance.Tracef(format, args...)
}
