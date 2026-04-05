package transport

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

func TestCreateRoomRequest_Validate(t *testing.T) {
	req := CreateRoomRequest{Name: "Test Room"}
	assert.NoError(t, req.Validate())

	req.Name = ""
	assert.Error(t, req.Validate())

	cap := 0
	req.Capacity = &cap
	assert.Error(t, req.Validate())
}

func TestToRoomResponse(t *testing.T) {
	room := &domain.Room{
		ID:        "room1",
		Name:      "Conference Room",
		CreatedAt: time.Now(),
	}
	resp := ToRoomResponse(room)
	assert.Equal(t, "room1", resp.ID)
	assert.Equal(t, "Conference Room", resp.Name)
}

func TestToRoomResponseList(t *testing.T) {
	rooms := []*domain.Room{
		{ID: "1", Name: "Room 1"},
		{ID: "2", Name: "Room 2"},
	}
	list := ToRoomResponseList(rooms)
	assert.Len(t, list, 2)
	assert.Equal(t, "1", list[0].ID)
}

func TestToRoomResponse_Nil(t *testing.T) {
	resp := ToRoomResponse(nil)
	assert.Nil(t, resp)
}

func TestToRoomResponseList_Nil(t *testing.T) {
	list := ToRoomResponseList(nil)
	assert.Empty(t, list)
}
