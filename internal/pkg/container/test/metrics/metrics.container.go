package metricscontainer

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

// Options contains configuration for the metrics container.
type Options struct {
	Host      string
	Port      nat.Port
	ImageName string
	Tag       string
	Name      string
}

// StartContainer launches a metrics container.
func StartContainer(ctx context.Context, t *testing.T) (*ContainerInfo, error) {
	options := getDefaultMetricsOptions()

	containerReq := getContainerRequest(options)

	// Start metrics server container
	metricsContainer, err := testcontainerutils.StartContainer(ctx, containerReq)
	if err != nil {
		return nil, testcontainerutils.TerminateOnErr(ctx, metricsContainer, err, "failed to start metrics container")
	}

	// Ensure container cleanup
	t.Cleanup(func() {
		if err := testcontainerutils.TerminateContainer(ctx, metricsContainer); err != nil {
			t.Fatalf("failed to cleanup container: %s", err)
		}
	})

	host, err := metricsContainer.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := metricsContainer.MappedPort(ctx, options.Port)
	if err != nil {
		return nil, err
	}
	connectionString := fmt.Sprintf("http://%s:%s", host, port.Port())
	logrus.Infof("Metrics running at %s", connectionString)

	return &ContainerInfo{
		Host: host,
		Port: port.Int(),
	}, nil
}

// ContainerInfo holds information about the metrics container.
type ContainerInfo struct {
	Host string
	Port int
}

// getDefaultMetricsOptions returns default options for the metrics container.
func getDefaultMetricsOptions() *Options {
	port, err := nat.NewPort("tcp", "9090")
	if err != nil {
		panic(errors.Wrap(err, "failed to create new port"))
	}

	return &Options{
		Port:      port,
		Host:      "localhost",
		ImageName: "prom/prometheus",
		Name:      "metrics-testcontainer",
		Tag:       "latest",
	}
}

// getContainerRequest creates a testcontainers.ContainerRequest for the metrics container.
func getContainerRequest(opts *Options) testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", opts.ImageName, opts.Tag),
		ExposedPorts: []string{opts.Port.Port()},
		WaitingFor:   wait.ForHTTP("/metrics").WithPort(opts.Port),
	}
}
