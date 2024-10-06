package echo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContextCreation(t *testing.T) {
	ctx := NewContext()
	assert.NotNil(t, ctx, "Expected context to be created and not nil")
}

func TestNewContextCreation(t *testing.T) {
	ctx := NewContext()

	if ctx == nil {
		t.Fatal("Expected non-nil context")
	}

	select {
	case <-ctx.Done():
		t.Fatal("Expected context to be active, but it was canceled")
	default:
		// Context is active
	}
}

func TestNewContextNoInterrupt(t *testing.T) {
	ctx := NewContext()

	select {
	case <-ctx.Done():
		t.Fatal("Expected context to remain active, but it was canceled")
	case <-time.After(1 * time.Second):
		// No interrupt signal received, context remains active
	}
}
