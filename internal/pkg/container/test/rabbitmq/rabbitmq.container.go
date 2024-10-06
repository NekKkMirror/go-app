package rabbitmqcontainer

import (
	"context"
	"fmt"
	"testing"
	"time"

	testcontainerutils "github.com/NekKkMirror/go-app/internal/pkg/container/test/utils"
	myrmq "github.com/NekKkMirror/go-app/internal/pkg/rabbitmq"
	"github.com/cenkalti/backoff/v4"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Options holds configuration for the RabbitMQ container.
type Options struct {
	Host           string
	Port           nat.Port
	UserName       string
	Password       string
	VirtualHost    string
	ImageName      string
	Tag            string
	Name           string
	ManagementPort nat.Port
}

// Start initializes a RabbitMQ container and returns an amqp connection and any error occurred.
func Start(ctx context.Context, t *testing.T) (*amqp.Connection, error) {
	options := getDefaultRabbitMQOptions()
	containerReq := getContainerRequest(options)

	// Start RabbitMQ container
	rabbitContainer, err := testcontainerutils.StartContainer(ctx, containerReq)
	if err != nil {
		return nil, testcontainerutils.TerminateOnErr(ctx, rabbitContainer, err, "failed to start RabbitMQ container")
	}

	// Ensure container cleanup
	t.Cleanup(func() {
		if err := testcontainerutils.TerminateContainer(ctx, rabbitContainer); err != nil {
			t.Fatalf("failed to cleanup container: %s", err)
		}
	})

	options.Host, err = rabbitContainer.Host(ctx)
	if err != nil {
		return nil, err
	}
	options.ManagementPort, err = rabbitContainer.MappedPort(ctx, options.ManagementPort)
	if err != nil {
		return nil, err
	}
	logrus.Infof("RabbitMQ management interface running at http://%s:%d", options.Host, options.ManagementPort.Int())

	conn, err := createRabbitMQConnection(ctx, rabbitContainer, options)
	if err != nil {
		return nil, testcontainerutils.TerminateOnErr(ctx, rabbitContainer, err, "failed to create RabbitMQ connection")
	}

	// Create default exchange
	err = createDefaultExchange(conn, options)
	if err != nil {
		return nil, testcontainerutils.TerminateOnErr(ctx, rabbitContainer, err, "failed to create default exchange")
	}

	return conn, nil
}

// getDefaultRabbitMQOptions returns the default configuration for RabbitMQ container.
func getDefaultRabbitMQOptions() *Options {
	port, err := nat.NewPort("tcp", "5672")
	if err != nil {
		panic(errors.Wrap(err, "failed to create new port"))
	}

	managementPort, err := nat.NewPort("tcp", "15672")
	if err != nil {
		panic(errors.Wrap(err, "failed to create new management port"))
	}

	return &Options{
		Port:           port,
		ManagementPort: managementPort,
		Host:           "localhost",
		UserName:       "guest",
		Password:       "guest",
		ImageName:      "rabbitmq",
		Name:           "rabbitmq-testcontainer",
		Tag:            "3-management",
		VirtualHost:    "/",
	}
}

// getContainerRequest builds and returns a testcontainers.ContainerRequest using the provided options.
func getContainerRequest(opts *Options) testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", opts.ImageName, opts.Tag),
		ExposedPorts: []string{opts.Port.Port(), opts.ManagementPort.Port()},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(opts.Port),
			wait.ForListeningPort(opts.ManagementPort),
		),
		Env: map[string]string{
			"RABBITMQ_DEFAULT_USER":  opts.UserName,
			"RABBITMQ_DEFAULT_PASS":  opts.Password,
			"RABBITMQ_DEFAULT_VHOST": opts.VirtualHost,
		},
	}
}

// createRabbitMQConnection establishes a RabbitMQ connection using provided RabbitMQ container and options.
func createRabbitMQConnection(ctx context.Context, container testcontainers.Container, opts *Options) (*amqp.Connection, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second
	const maxRetries = 5

	var (
		conn *amqp.Connection
		err  error
		cfg  *myrmq.RMQConfig
	)

	err = backoff.Retry(func() error {

		opts.Port, err = container.MappedPort(ctx, opts.Port)
		if err != nil {
			return errors.Wrap(err, "failed to get exposed container port")
		}

		cfg = &myrmq.RMQConfig{
			Host:         opts.Host,
			Port:         opts.Port.Int(),
			User:         opts.UserName,
			Password:     opts.Password,
			ExchangeName: "default",
			Kind:         "direct",
			ContentType:  "application/json",
			MaxRetries:   maxRetries,
			RetryDelay:   2 * time.Second,
		}

		conn, err = myrmq.NewRMQConn(cfg, ctx)
		return err
	}, backoff.WithMaxRetries(bo, uint64(maxRetries-1)))

	if err != nil {
		return nil, errors.Wrap(err, "failed to create connection after retries")
	}

	return conn, nil
}

// createDefaultExchange creates a default exchange in RabbitMQ.
func createDefaultExchange(conn *amqp.Connection, opts *Options) error {
	channel, err := conn.Channel()
	if err != nil {
		return errors.Wrap(err, "failed to open a channel")
	}
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		if err != nil {

		}
	}(channel)

	err = channel.ExchangeDeclare(
		"default", // name
		"direct",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare default exchange")
	}

	return nil
}
