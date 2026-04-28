package domain

import (
	"time"
)

const SlotDuration = 30 * time.Minute

type Slot struct {
	ID        string
	RoomID    string
	StartTime time.Time
	EndTime   time.Time `json:"endTime"`
	IsBooked  bool      `json:"isBooked"`
	BookingID string    `json:"bookingId,omitempty"`
	CreatedAt time.Time `json:"-"`
} // создается при генерации слота

func (s *Slot) IsInPast() bool {
	return time.Now().After(s.StartTime)
}
