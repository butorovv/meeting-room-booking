package domain

import "time"

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
	EndTime string // "18:00"
	CreatedAt time.Time
}