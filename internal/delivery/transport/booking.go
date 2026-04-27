package transport

import (
	"fmt"
	"time"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

type CreateBookingRequest struct {
	RoomID               string `json:"roomId"`
	StartTime            string `json:"startTime"`
	EndTime              string `json:"endTime"`
	CreateConferenceLink bool   `json:"createConferenceLink"`
}

func (r *CreateBookingRequest) Validate() error {
	if r.RoomID == "" {
		return fmt.Errorf("roomId is required")
	}
	if r.StartTime == "" {
		return fmt.Errorf("startTime is required")
	}
	if r.EndTime == "" {
		return fmt.Errorf("endTime is required")
	}
	startTime, err := time.Parse(time.RFC3339, r.StartTime)
	if err != nil {
		return fmt.Errorf("startTime must be RFC3339")
	}
	endTime, err := time.Parse(time.RFC3339, r.EndTime)
	if err != nil {
		return fmt.Errorf("endTime must be RFC3339")
	}
	if endTime.Sub(startTime) != domain.SlotDuration {
		return fmt.Errorf("duration must be 30 minutes")
	}
	return nil
}

type BookingResponse struct {
	ID             string    `json:"id"`
	SlotID         *string   `json:"slotId,omitempty"`
	RoomID         string    `json:"roomId"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
	UserID         string    `json:"userId"`
	Status         string    `json:"status"`
	ConferenceLink *string   `json:"conferenceLink,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

type BookingsListResponse struct {
	Bookings   []BookingResponse  `json:"bookings"`
	Pagination PaginationResponse `json:"pagination"`
}

type PaginationResponse struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

func ToBookingResponse(b *domain.Booking) *BookingResponse {
	if b == nil {
		return nil
	}
	return &BookingResponse{
		ID:             b.ID,
		SlotID:         b.SlotID,
		RoomID:         b.RoomID,
		StartTime:      b.StartTime,
		EndTime:        b.EndTime,
		UserID:         b.UserID,
		Status:         string(b.Status),
		ConferenceLink: b.ConferenceLink,
		CreatedAt:      b.CreatedAt,
	}
}

func ToBookingResponseList(bookings []*domain.Booking) []BookingResponse {
	if bookings == nil {
		return []BookingResponse{}
	}
	result := make([]BookingResponse, 0, len(bookings))
	for _, b := range bookings {
		result = append(result, *ToBookingResponse(b))
	}
	return result
}
