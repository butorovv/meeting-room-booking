package transport

import (
	"fmt"
	"time"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

type GetAvailableSlotsParams struct {
	RoomID string // from path
	Date   string // from query
}

func (p *GetAvailableSlotsParams) Validate() error {
	if p.RoomID == "" {
		return fmt.Errorf("roomId is required")
	}
	if p.Date == "" {
		return fmt.Errorf("date is required")
	}
	return nil
}

type SlotResponse struct {
	ID     string    `json:"id"`
	RoomID string    `json:"roomId"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}

func ToSlotResponse(s *domain.Slot) *SlotResponse {
	if s == nil {
		return nil
	}
	return &SlotResponse{
		ID:     s.ID,
		RoomID: s.RoomID,
		Start:  s.StartTime,
		End:    s.EndTime,
	}
}

func ToSlotResponseList(slots []*domain.Slot) []SlotResponse {
	if slots == nil {
		return []SlotResponse{}
	}
	result := make([]SlotResponse, 0, len(slots))
	for _, s := range slots {
		result = append(result, *ToSlotResponse(s))
	}
	return result
}
