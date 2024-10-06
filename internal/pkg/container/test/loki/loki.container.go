package lokicontainer

import (
	"context"
	"fmt"
	"testing"

	testcontainerutils "github.com/NekKkMirror/go-app/internal/pkg/container/test/utils"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Options holds configuration for the Loki container.
type Options struct {
	Host      string
	Port      nat.Port
	ImageName string
	Tag       string
	Name      string
}

// Start initializes a Loki container and returns a connection string and any error occurred.
func Start(ctx context.Context, t *testing.T) (string, error) {
	options := getDefaultLokiOptions()
	containerReq := getContainerRequest(options)

	// Start Loki container.
	lokiContainer, err := testcontainerutils.StartContainer(ctx, containerReq)
	if err != nil {
		return "", testcontainerutils.TerminateOnErr(ctx, lokiContainer, err, "failed to start Loki container")
	}

	// Ensure container cleanup.
	t.Cleanup(func() {
		if err := testcontainerutils.TerminateContainer(ctx, lokiContainer); err != nil {
			t.Fatalf("failed to cleanup container: %s", err)
		}
	})

	options.Host, err = lokiContainer.Host(ctx)
	if err != nil {
		return "", err
	}
	options.Port, err = lokiContainer.MappedPort(ctx, options.Port)
	if err != nil {
		return "", err
	}
	connectionString := fmt.Sprintf("http://%s:%s", options.Host, options.Port.Port())
	logrus.Infof("Loki running at %s", connectionString)

	return connectionString, nil
}

// getDefaultLokiOptions returns the default configuration for Loki container.
func getDefaultLokiOptions() *Options {
	port, err := nat.NewPort("tcp", "3100")
	if err != nil {
		panic(errors.Wrap(err, "failed to create port 3100"))
	}

	return &Options{
		Port:      port,
		Host:      "localhost",
		ImageName: "grafana/loki",
		Name:      "loki-testcontainer",
		Tag:       "2.4.1",
	}
}

// getContainerRequest builds and returns a testcontainers.ContainerRequest using the provided options.
func getContainerRequest(opts *Options) testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", opts.ImageName, opts.Tag),
		ExposedPorts: []string{opts.Port.Port()},
		WaitingFor: wait.ForAll(
			wait.ForHTTP("/ready").WithPort(opts.Port),
		),
	}
}
