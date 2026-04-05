package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCancel(t *testing.T) {
	b := &Booking{Status: BookingActive}
	b.Cancel()
	assert.Equal(t, BookingCancelled, b.Status)

	// идемпотентность
	b.Cancel()
	assert.Equal(t, BookingCancelled, b.Status)
}
