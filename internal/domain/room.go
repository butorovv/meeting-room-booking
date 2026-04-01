package domain

import (
	"time"
)

type Room struct {
	ID          string
	Name        string
	Description *string // nullable по API
	Capacity    *int    // nullable по API
	CreatedAt   time.Time
}
