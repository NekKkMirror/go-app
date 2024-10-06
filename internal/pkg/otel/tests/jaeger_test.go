package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/NekKkMirror/go-app/internal/pkg/constants/environment"
	"github.com/NekKkMirror/go-app/internal/pkg/logger/mocks"
	myotel "github.com/NekKkMirror/go-app/internal/pkg/otel"
	"github.com/NekKkMirror/go-app/internal/pkg/utils/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel"
)

func TestTracerProvider(t *testing.T) {
	tests := []struct {
		name          string
		config        myotel.OTLPConfig
		expectedEnv   string
		expectError   bool
		expectedError string
	}{
		{
			name: "valid configuration in development",
			config: myotel.OTLPConfig{
				Server:      "localhost:4318",
				ServiceName: "test-service",
				TracerName:  "test-tracer",
			},
			expectedEnv: "development",
			expectError: false,
		},
		{
			name: "environment set to production",
			config: myotel.OTLPConfig{
				Server:      "localhost:4318",
				ServiceName: "test-service",
				TracerName:  "test-tracer",
			},
			expectedEnv: "production",
			expectError: false,
		},
		{
			name: "invalid configuration",
			config: myotel.OTLPConfig{
				Server:      "",
				ServiceName: "test-service",
				TracerName:  "test-tracer",
			},
			expectedEnv:   "development",
			expectError:   true,
			expectedError: "failed to create exporter: signal: no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLog := new(mocks.ILogger)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			oldEnv := common.GetEnv("APP_ENV", environment.Development)
			defer os.Setenv("APP_ENV", oldEnv) // reset environment variable after test
			os.Setenv("APP_ENV", tt.expectedEnv)

			// Setting expectations based on the implementation
			if tt.config.Server == "" {
				mockLog.On("Error", "Failed to initialize TracerProvider", "error", mock.Anything).Return(nil).Once()
			} else {
				// Logging related to setup
				mockLog.On("Info", "TracerProvider initialized").Return(nil).Once()
				// Logging related to shutdown should only be triggered at the end.
				mockLog.On("Error", "Failed to shutdown TracerProvider", "error", mock.Anything).Maybe()
				mockLog.On("Info", "open-telemetry exited properly").Maybe()
			}

			// Call the TracerProvider
			tracer, err := myotel.TracerProvider(ctx, &tt.config, mockLog)

			if tt.expectError {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Equal(t, otel.Tracer(tt.config.TracerName), tracer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tracer)

				// Validate the tracer's functionality by starting and ending a span
				_, span := tracer.Start(ctx, "test-span")
				assert.NotNil(t, span)
				span.End()
			}

			// Allow background goroutine logger calls to execute
			time.Sleep(1 * time.Second)

			// Check the expectations
			mockLog.AssertExpectations(t)
		})
	}
}
