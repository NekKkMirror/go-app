package logger

// ILogger defines the interface for our logger
//
//go:generate mockery --name ILogger
type ILogger interface {
	// Debug logs a message at level Debug.
	Debug(args ...interface{})

	// Debugf logs a message at level Debug using the provided format and arguments.
	Debugf(format string, args ...interface{})

	// Info logs a message at level Info.
	Info(args ...interface{})

	// Infof logs a message at level Info using the provided format and arguments.
	Infof(format string, args ...interface{})

	// Warn logs a message at level Warn.
	Warn(args ...interface{})

	// Warnf logs a message at level Warn using the provided format and arguments.
	Warnf(format string, args ...interface{})

	// Error logs a message at level Error.
	Error(args ...interface{})

	// Errorf logs a message at level Error using the provided format and arguments.
	Errorf(format string, args ...interface{})

	// Fatal logs a message at level Fatal then the process will exit with status set to 1.
	Fatal(args ...interface{})

	// Fatalf logs a message at level Fatal using the provided format and arguments then the process will exit with status set to 1.
	Fatalf(format string, args ...interface{})

	// Trace logs a message at level Trace.
	Trace(args ...interface{})

	// Tracef logs a message at level Trace using the provided format and arguments.
	Tracef(format string, args ...interface{})
}
