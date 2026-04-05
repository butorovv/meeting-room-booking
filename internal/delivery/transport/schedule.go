package transport

import (
	"fmt"
	"strings"
	"time"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

type CreateScheduleRequest struct {
	RoomID     string `json:"roomId"`
	DaysOfWeek []int  `json:"daysOfWeek"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

func (r *CreateScheduleRequest) Validate() error {
	if strings.TrimSpace(r.RoomID) == "" {
		return fmt.Errorf("roomId is required")
	}
	if len(r.DaysOfWeek) == 0 {
		return fmt.Errorf("daysOfWeek is required")
	}
	for _, d := range r.DaysOfWeek {
		if d < 1 || d > 7 {
			return fmt.Errorf("daysOfWeek must be between 1 and 7")
		}
	}
	if strings.TrimSpace(r.StartTime) == "" {
		return fmt.Errorf("startTime is required")
	}
	if strings.TrimSpace(r.EndTime) == "" {
		return fmt.Errorf("endTime is required")
	}

	if _, err := time.Parse("15:04", r.StartTime); err != nil {
		return fmt.Errorf("startTime must be in format HH:MM")
	}
	if _, err := time.Parse("15:04", r.EndTime); err != nil {
		return fmt.Errorf("endTime must be in format HH:MM")
	}

	if r.StartTime >= r.EndTime {
		return fmt.Errorf("startTime must be less than endTime")
	}

	return nil
}

type ScheduleResponse struct {
	ID         string    `json:"id"`
	RoomID     string    `json:"roomId"`
	DaysOfWeek []int     `json:"daysOfWeek"`
	StartTime  string    `json:"startTime"`
	EndTime    string    `json:"endTime"`
	CreatedAt  time.Time `json:"createdAt"`
}

func ToScheduleResponse(s *domain.Schedule) *ScheduleResponse {
	if s == nil {
		return nil
	}
	return &ScheduleResponse{
		ID:         s.ID,
		RoomID:     s.RoomID,
		DaysOfWeek: domain.MaskToDays(s.DaysMask),
		StartTime:  s.StartTime,
		EndTime:    s.EndTime,
		CreatedAt:  s.CreatedAt,
	}
}
