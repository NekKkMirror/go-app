package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

//go:generate mockery --name ILogger
type ILogger interface {
	GetLevel() log.Level
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
}

var (
	Logger ILogger
)

type Config struct {
	LogLevel string
	logger   *log.Logger
}

type appLogger struct {
	level  string
	logger *log.Logger
}

var loggerLeveMap = map[string]log.Level{
	"debug": log.DebugLevel,
	"info":  log.InfoLevel,
	"warn":  log.WarnLevel,
	"error": log.ErrorLevel,
	"panic": log.PanicLevel,
	"fatal": log.FatalLevel,
	"trace": log.TraceLevel,
}

func (l *appLogger) GetLevel() log.Level {
	level, exist := loggerLeveMap[l.level]
	if !exist {
		return log.DebugLevel
	}

	return level
}

func InitLogger(cfg *Config) ILogger {
	l := &appLogger{level: cfg.LogLevel}

	l.logger = log.StandardLogger()

	logLevel := l.GetLevel()

	env := os.Getenv("APP_END")

	if env == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: false,
			ForceColors:   true,
			FullTimestamp: true,
		})
	}

	log.SetLevel(logLevel)

	return l
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
