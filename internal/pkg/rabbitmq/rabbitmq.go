package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// MessageHandler defines a function type for custom message handling logic.
type MessageHandler[T any] func(ctx context.Context, msg T, dependencies T) error

type RMQConfig struct {
	Host              string
	Port              int
	User              string
	Password          string
	ExchangeName      string
	Kind              string
	ContentType       string
	DefaultPriority   uint8
	DefaultExpiration string
	ReplyToQueue      string
	MaxRetries        int
	RetryDelay        time.Duration
}

// NewRMQConn establishes a connection to RabbitMQ using the provided configuration and context.
func NewRMQConn(cfg *RMQConfig, ctx context.Context) (*amqp.Connection, error) {
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second // Maximum time to retry
	maxRetries := 5                      // Number of retries (including the initial attempt)

	var conn *amqp.Connection
	var err error

	err = backoff.Retry(func() error {
		conn, err = amqp.Dial(connAddr)
		if err != nil {
			log.Errorf("Failed to connect to RabbitMQ: %v. Connection information: %s", err, connAddr)
			return err
		}
		return nil
	}, backoff.WithMaxRetries(bo, uint64(maxRetries-1)))

	if err != nil {
		return nil, err
	}

	log.Info("Connected to RabbitMQ")

	go func() {
		select {
		case <-ctx.Done():
			err := conn.Close()
			if err != nil {
				log.Error("Failed to close RabbitMQ connection")
			}
			log.Info("RabbitMQ connection is closed")
		}
	}()

	return conn, nil
}
