package transport

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

func TestGetAvailableSlotsParams_Validate(t *testing.T) {
	p := GetAvailableSlotsParams{RoomID: "room1", Date: "2024-04-07"}
	assert.NoError(t, p.Validate())

	p.RoomID = ""
	assert.Error(t, p.Validate())

	p.Date = ""
	assert.Error(t, p.Validate())
}

func TestToSlotResponse(t *testing.T) {
	slot := &domain.Slot{
		ID:        "slot1",
		RoomID:    "room1",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(30 * time.Minute),
	}
	resp := ToSlotResponse(slot)
	assert.Equal(t, "slot1", resp.ID)
	assert.Equal(t, "room1", resp.RoomID)
}

func TestToSlotResponseList(t *testing.T) {
	slots := []*domain.Slot{{ID: "1"}, {ID: "2"}}
	list := ToSlotResponseList(slots)
	assert.Len(t, list, 2)
}

func TestToSlotResponse_Nil(t *testing.T) {
	resp := ToSlotResponse(nil)
	assert.Nil(t, resp)
}

func TestToSlotResponseList_Nil(t *testing.T) {
	list := ToSlotResponseList(nil)
	assert.Empty(t, list)
}
