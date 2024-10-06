package testcontainerutils

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
)

// StartContainer starts a new testcontainer with the given request configuration.
func StartContainer(ctx context.Context, req testcontainers.ContainerRequest) (testcontainers.Container, error) {
	genericReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}
	return testcontainers.GenericContainer(ctx, genericReq)
}

// TerminateOnErr stops the container and returns the wrapped error message.
func TerminateOnErr(ctx context.Context, container testcontainers.Container, err error, msg string) error {
	if container != nil {
		if stopErr := TerminateContainer(ctx, container); stopErr != nil {
			return errors.Wrap(stopErr, fmt.Sprintf("%s: %v", msg, err))
		}
	}
	return errors.Wrap(err, msg)
}

// TerminateContainer stops and terminates the given testcontainer, removing all volumes.
func TerminateContainer(ctx context.Context, container testcontainers.Container) error {
	return errors.Wrap(container.Terminate(ctx), "failed to terminate container")
}
