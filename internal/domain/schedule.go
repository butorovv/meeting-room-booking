package domain

import (
	"time"

	"github.com/google/uuid"
)

type DaysMask int16

const (
	Monday    DaysMask = 1 << iota // 1
	Tuesday                        // 2
	Wednesday                      // 4
	Thursday                       // 8
	Friday                         // 16
	Saturday                       // 32
	Sunday                         // 64
)

// DaysToMask конвертирует массив дней (1..7) в битовую маску
func DaysToMask(days []int) DaysMask {
	var mask DaysMask
	for _, d := range days {
		switch d {
		case 1:
			mask |= Monday
		case 2:
			mask |= Tuesday
		case 3:
			mask |= Wednesday
		case 4:
			mask |= Thursday
		case 5:
			mask |= Friday
		case 6:
			mask |= Saturday
		case 7:
			mask |= Sunday
		}
	}

	return mask
}

// MaskToDays конвертирует битовую маску в массив дней (1..7)
func MaskToDays(mask DaysMask) []int {
	var days []int
	if mask&Monday != 0 {
		days = append(days, 1)
	}
	if mask&Tuesday != 0 {
		days = append(days, 2)
	}
	if mask&Wednesday != 0 {
		days = append(days, 3)
	}
	if mask&Thursday != 0 {
		days = append(days, 4)
	}
	if mask&Friday != 0 {
		days = append(days, 5)
	}
	if mask&Saturday != 0 {
		days = append(days, 6)
	}
	if mask&Sunday != 0 {
		days = append(days, 7)
	}

	return days
}

type Schedule struct {
	ID        string
	RoomID    string
	DaysMask  DaysMask
	StartTime string // "09:00"
	EndTime   string // "18:00"
	CreatedAt time.Time
}

// WeekdayToMask конвертирует time.Weekday в DaysMask
func WeekdayToMask(weekday time.Weekday) DaysMask {
	switch weekday {
	case time.Monday:
		return Monday
	case time.Tuesday:
		return Tuesday
	case time.Wednesday:
		return Wednesday
	case time.Thursday:
		return Thursday
	case time.Friday:
		return Friday
	case time.Saturday:
		return Saturday
	case time.Sunday:
		return Sunday
	default:
		return 0
	}
}

func (s *Schedule) GenerateSlots(startDate, endDate time.Time) []*Slot {
	var slots []*Slot
	current := startDate
	for current.Before(endDate) {
		// проверяем, что день недели подходит
		dayMask := WeekdayToMask(current.Weekday())
		if s.DaysMask&dayMask != 0 {
			startTime, _ := time.Parse("15:04", s.StartTime)
			endTime, _ := time.Parse("15:04", s.EndTime)

			slotStart := time.Date(current.Year(), current.Month(), current.Day(),
				startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
			slotEnd := time.Date(current.Year(), current.Month(), current.Day(),
				endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)

			for t := slotStart; t.Before(slotEnd); t = t.Add(SlotDuration) {
				slots = append(slots, &Slot{
					ID:        uuid.NewString(),
					RoomID:    s.RoomID,
					StartTime: t,
					EndTime:   t.Add(SlotDuration),
					CreatedAt: time.Now().UTC(),
				})
			}
		}
		current = current.AddDate(0, 0, 1)
	}
	return slots
}
