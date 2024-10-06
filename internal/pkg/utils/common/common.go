package common

import (
	"os"

	"github.com/NekKkMirror/go-app/internal/pkg/constants/environment"
)

// IsTestEnvironment checks if the application is running in a test environment.
func IsTestEnvironment() bool {
	return GetEnv("APP_ENV", environment.Development) == environment.Test
}

// IsDevelopmentEnvironment checks if the application is running in a development environment.
func IsDevelopmentEnvironment() bool {
	return GetEnv("APP_ENV", environment.Development) == environment.Development
}

// GetEnv retrieves environment variables or returns a default value.
func GetEnv(string, defaultValue string) string {
	if value, exists := os.LookupEnv(string); exists {
		return value
	}
	return defaultValue
}
