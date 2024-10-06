package echo

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// NewContext creates a new context that listens for interrupt signals and cancels the context accordingly.
// The context is created using signal.NotifyContext with os.Interrupt, syscall.SIGTERM, and syscall.SIGINT.
// A goroutine is started to listen for these signals and cancel the context when any of them are received.
// The function returns the created context.
func NewContext() context.Context {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("context canceled")
				cancel()
				return
			}
		}
	}()

	return ctx
}
