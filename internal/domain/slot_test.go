package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsInPast(t *testing.T) {
	past := Slot{StartTime: time.Now().Add(-1 * time.Hour)}
	future := Slot{StartTime: time.Now().Add(1 * time.Hour)}

	assert.True(t, past.IsInPast())
	assert.False(t, future.IsInPast())
}
