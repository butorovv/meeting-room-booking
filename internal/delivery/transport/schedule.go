package transport

import (
	"fmt"
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
	if r.RoomID == "" {
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
	if r.StartTime == "" {
		return fmt.Errorf("startTime is required")
	}
	if r.EndTime == "" {
		return fmt.Errorf("endTime is required")
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
