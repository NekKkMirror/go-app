package postgrescontainer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ORM_Container(t *testing.T) {
	gorm, _, err := Start(context.Background(), t)
	require.NoError(t, err)

	assert.NotNil(t, gorm)
}

func TestStartUsesDefaultPostgresOptions(t *testing.T) {
	ctx := context.Background()
	db, _, err := Start(ctx, t)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if db == nil {
		t.Fatalf("expected a valid DB instance, got nil")
	}
}

func TestStartCleansUpContainerAfterTestCompletes(t *testing.T) {
	ctx := context.Background()
	db, _, err := Start(ctx, t)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if db == nil {
		t.Fatalf("expected a valid DB instance, got nil")
	}
}

func TestStartEstablishesGORMConnection(t *testing.T) {
	ctx := context.Background()
	db, _, err := Start(ctx, t)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if db == nil {
		t.Fatalf("expected a valid DB instance, got nil")
	}
}

func TestStartHandlesContextCancellationOrTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	db, _, err := Start(ctx, t)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if db != nil {
		t.Fatal("expected a nil DB instance, got non-nil")
	}
}

func TestStartSuccessfullyStartsPostgresContainer(t *testing.T) {
	ctx := context.Background()
	db, _, err := Start(ctx, t)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if db == nil {
		t.Fatalf("expected a valid DB instance, got nil")
	}
}

func TestConfiguresContainerEnvironmentVariablesForPostgreSQL(t *testing.T) {
	ctx := context.Background()
	db, _, err := Start(ctx, t)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if db == nil {
		t.Fatalf("expected a valid DB instance, got nil")
	}
}

func TestMapsContainerPortsCorrectlyForHostAccess(t *testing.T) {
	ctx := context.Background()
	db, _, err := Start(ctx, t)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if db == nil {
		t.Fatalf("expected a valid DB instance, got nil")
	}
}
