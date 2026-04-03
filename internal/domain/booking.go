package domain

import (
	"time"
)

type BookingStatus string

const (
	BookingActive    BookingStatus = "active"
	BookingCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID             string
	SlotID         string
	UserID         string
	Status         BookingStatus
	ConferenceLink *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (b *Booking) Cancel() {
	if b.Status == BookingActive {
		b.Status = BookingCancelled
		b.UpdatedAt = time.Now().UTC()
	}
}
