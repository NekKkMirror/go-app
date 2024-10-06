package logger

import (
	"os"

	"sync"

	log "github.com/sirupsen/logrus"
)

// Logger holds the singleton instance of the logger
var Logger ILogger
var once sync.Once

type Config struct {
	LogLevel string
}

type appLogger struct {
	level  string
	logger *log.Logger
}

var loggerLevelMap = map[string]log.Level{
	"debug": log.DebugLevel,
	"info":  log.InfoLevel,
	"warn":  log.WarnLevel,
	"error": log.ErrorLevel,
	"panic": log.PanicLevel,
	"fatal": log.FatalLevel,
	"trace": log.TraceLevel,
}

// GetLevel returns the log level set in config or defaults to DebugLevel
func (l *appLogger) GetLevel() log.Level {
	level, exist := loggerLevelMap[l.level]
	if !exist {
		return log.DebugLevel
	}
	return level
}

// InitLogger initializes the logger with the given config
func InitLogger(cfg *Config) ILogger {
	once.Do(func() {
		l := &appLogger{level: cfg.LogLevel}
		l.logger = log.StandardLogger()

		l.setupFormatter()
		log.SetLevel(l.GetLevel())

		Logger = l
	})
	return Logger
}

func (l *appLogger) setupFormatter() {
	env := os.Getenv("APP_ENV")
	if env == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: false,
			ForceColors:   true,
			FullTimestamp: true,
		})
	}
}

func (l *appLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *appLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *appLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *appLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *appLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *appLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *appLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *appLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *appLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *appLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *appLogger) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l *appLogger) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}
