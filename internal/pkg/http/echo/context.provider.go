package echo

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// NewContext creates a new context that listens for system interrupt and termination signals.
// It returns the context which will be canceled when a signal is caught.
func NewContext() context.Context {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		waitForSignal(ctx, cancel)
	}()

	return ctx
}

// waitForSignal waits for the context to be done and logs the cancellation.
func waitForSignal(ctx context.Context, cancel context.CancelFunc) {
	<-ctx.Done()
	log.Info("context canceled")
	cancel()
}
