package transport

import (
	"fmt"
	"time"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

type CreateBookingRequest struct {
	SlotID               string `json:"slotId"`
	CreateConferenceLink bool   `json:"createConferenceLink"`
}

func (r *CreateBookingRequest) Validate() error {
	if r.SlotID == "" {
		return fmt.Errorf("slotId is required")
	}
	return nil
}

type BookingResponse struct {
	ID             string    `json:"id"`
	SlotID         string    `json:"slotId"`
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
