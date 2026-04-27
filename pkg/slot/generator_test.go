package slot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

func TestGenerateSlots_WorkingDay(t *testing.T) {
	schedule := &domain.Schedule{
		RoomID:    "11111111-1111-1111-1111-111111111111",
		DaysMask:  domain.Monday,
		StartTime: "09:00",
		EndTime:   "10:00",
	}
	date := time.Date(2024, 4, 8, 12, 15, 0, 0, time.FixedZone("MSK", 3*60*60))

	slots := GenerateSlots(schedule, date)
	sameSlots := GenerateSlots(schedule, date)

	assert.Len(t, slots, 2)
	assert.Equal(t, slots[0].ID, sameSlots[0].ID)
	assert.Equal(t, time.Date(2024, 4, 8, 9, 0, 0, 0, time.UTC), slots[0].StartTime)
	assert.Equal(t, time.Date(2024, 4, 8, 9, 30, 0, 0, time.UTC), slots[0].EndTime)
	assert.Equal(t, "11111111-1111-1111-1111-111111111111", slots[0].RoomID)
}

func TestGenerateSlots_NonWorkingDay(t *testing.T) {
	schedule := &domain.Schedule{
		RoomID:    "11111111-1111-1111-1111-111111111111",
		DaysMask:  domain.Monday,
		StartTime: "09:00",
		EndTime:   "10:00",
	}
	date := time.Date(2024, 4, 9, 0, 0, 0, 0, time.UTC)

	slots := GenerateSlots(schedule, date)

	assert.Empty(t, slots)
}

func TestGenerateSlots_WithSecondsClock(t *testing.T) {
	schedule := &domain.Schedule{
		RoomID:    "11111111-1111-1111-1111-111111111111",
		DaysMask:  domain.Monday,
		StartTime: "09:00:00",
		EndTime:   "09:30:00",
	}
	date := time.Date(2024, 4, 8, 0, 0, 0, 0, time.UTC)

	slots := GenerateSlots(schedule, date)

	assert.Len(t, slots, 1)
	assert.Equal(t, 30*time.Minute, slots[0].EndTime.Sub(slots[0].StartTime))
}
