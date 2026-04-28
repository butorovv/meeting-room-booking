package transport

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

func TestCreateBookingRequest_Validate(t *testing.T) {
	startTime := time.Date(2030, 4, 8, 9, 0, 0, 0, time.UTC)
	req := CreateBookingRequest{
		RoomID:    "room1",
		StartTime: startTime.Format(time.RFC3339),
		EndTime:   startTime.Add(30 * time.Minute).Format(time.RFC3339),
	}
	assert.NoError(t, req.Validate())

	req.RoomID = ""
	assert.Error(t, req.Validate())
}

func TestToBookingResponse(t *testing.T) {
	slotID := "slot1"
	booking := &domain.Booking{
		ID:     "book1",
		SlotID: &slotID,
		UserID: "user1",
		Status: domain.BookingActive,
	}
	resp := ToBookingResponse(booking)
	assert.Equal(t, "book1", resp.ID)
	assert.Equal(t, "slot1", *resp.SlotID)
	assert.Equal(t, "active", resp.Status)
}

func TestToBookingResponseList(t *testing.T) {
	bookings := []*domain.Booking{{ID: "1"}, {ID: "2"}}
	list := ToBookingResponseList(bookings)
	assert.Len(t, list, 2)
}

func TestToBookingResponse_Nil(t *testing.T) {
	resp := ToBookingResponse(nil)
	assert.Nil(t, resp)
}

func TestToBookingResponseList_Nil(t *testing.T) {
	list := ToBookingResponseList(nil)
	assert.Empty(t, list)
}
