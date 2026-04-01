package domain

import ("time")

const SlotDuration = 30 * time.Minute

type Slot struct {
	ID string
	RoomID string
	StartTime time.Time
	EndTime time.Time
	CreatedAt time.Time // создается при генерации слота
}

func (s *Slot) isInPast() bool {
	return s.StartTime.Before(time.Now().UTC())
}