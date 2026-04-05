package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextWithLogger(t *testing.T) {
	ctx := ContextWithLogger(context.Background(), Global())
	log := FromContext(ctx)
	assert.NotNil(t, log)
}

func TestContextWithRequestID(t *testing.T) {
	ctx := ContextWithRequestID(context.Background(), "req123")
	id := RequestIDFromContext(ctx)
	assert.Equal(t, "req123", id)
}
