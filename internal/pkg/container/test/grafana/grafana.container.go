package grafanacontainer

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

type Options struct {
	Host          string
	Port          nat.Port
	AdminUser     string
	AdminPassword string
	ImageName     string
	Tag           string
	Name          string
}

// Start starts a Grafana container using the provided context
func Start(ctx context.Context, t *testing.T) (*ContainerInfo, error) {
	options := getDefaultGrafanaOptions()

	containerReq := getContainerRequest(options)

	// Start Grafana container
	grafanaContainer, err := testcontainerutils.StartContainer(ctx, containerReq)
	if err != nil {
		return nil, testcontainerutils.TerminateOnErr(ctx, grafanaContainer, err, "failed to start Grafana container")
	}

	// Ensure container cleanup
	t.Cleanup(func() {
		if err := testcontainerutils.TerminateContainer(ctx, grafanaContainer); err != nil {
			t.Fatalf("failed to cleanup container: %s", err)
		}
	})

	host, err := grafanaContainer.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := grafanaContainer.MappedPort(ctx, options.Port)
	if err != nil {
		return nil, err
	}
	connectionString := fmt.Sprintf("http://%s:%s", options.Host, port)
	logrus.Infof("Grafana running at %s", connectionString)

	return &ContainerInfo{
		Host:          host,
		Port:          port.Int(),
		AdminUser:     options.AdminUser,
		AdminPassword: options.AdminPassword,
	}, nil
}

type ContainerInfo struct {
	Host          string
	Port          int
	AdminUser     string
	AdminPassword string
}

// getDefaultGrafanaOptions returns a pointer to a new Options struct with default values for a Grafana container.
func getDefaultGrafanaOptions() *Options {
	port, err := nat.NewPort("tcp", "3000")
	if err != nil {
		panic(errors.Wrap(err, "failed to create new port"))
	}

	return &Options{
		Port:          port,
		Host:          "localhost",
		AdminUser:     "admin",
		AdminPassword: "admin",
		ImageName:     "grafana/grafana",
		Name:          "grafana-testcontainer",
		Tag:           "latest",
	}
}

// getContainerRequest creates a testcontainers.ContainerRequest for Grafana.
func getContainerRequest(opts *Options) testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", opts.ImageName, opts.Tag),
		ExposedPorts: []string{opts.Port.Port()},
		WaitingFor:   wait.ForHTTP("/login").WithPort(opts.Port),
		Env: map[string]string{
			"GF_SECURITY_ADMIN_USER":     opts.AdminUser,
			"GF_SECURITY_ADMIN_PASSWORD": opts.AdminPassword,
		},
	}
}
