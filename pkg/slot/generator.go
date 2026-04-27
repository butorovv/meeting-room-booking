package slot

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

var namespace = uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

func GenerateSlots(schedule *domain.Schedule, date time.Time) []*domain.Slot {
	if schedule == nil {
		return []*domain.Slot{}
	}

	day := utcDate(date)
	if schedule.DaysMask&domain.WeekdayToMask(day.Weekday()) == 0 {
		return []*domain.Slot{}
	}

	startClock, ok := parseClock(schedule.StartTime)
	if !ok {
		return []*domain.Slot{}
	}
	endClock, ok := parseClock(schedule.EndTime)
	if !ok {
		return []*domain.Slot{}
	}

	windowStart := time.Date(day.Year(), day.Month(), day.Day(), startClock.Hour(), startClock.Minute(), 0, 0, time.UTC)
	windowEnd := time.Date(day.Year(), day.Month(), day.Day(), endClock.Hour(), endClock.Minute(), 0, 0, time.UTC)
	if !windowStart.Before(windowEnd) {
		return []*domain.Slot{}
	}

	slots := make([]*domain.Slot, 0, int(windowEnd.Sub(windowStart)/domain.SlotDuration))
	for start := windowStart; !start.Add(domain.SlotDuration).After(windowEnd); start = start.Add(domain.SlotDuration) {
		end := start.Add(domain.SlotDuration)
		uniqueKey := fmt.Sprintf("%s:%s:%s", schedule.RoomID, start.Format(time.RFC3339), end.Format(time.RFC3339))
		slotID := uuid.NewSHA1(namespace, []byte(uniqueKey))

		slots = append(slots, &domain.Slot{
			ID:        slotID.String(),
			RoomID:    schedule.RoomID,
			StartTime: start,
			EndTime:   end,
			IsBooked:  false,
			CreatedAt: time.Now().UTC(),
		})
	}

	return slots
}

func utcDate(date time.Time) time.Time {
	utc := date.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}

func parseClock(value string) (time.Time, bool) {
	for _, layout := range []string{"15:04", "15:04:05"} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, true
		}
	}

	return time.Time{}, false
}
