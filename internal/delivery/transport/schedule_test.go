package transport

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

func TestCreateScheduleRequest_Validate(t *testing.T) {
	req := CreateScheduleRequest{
		RoomID:     "room1",
		DaysOfWeek: []int{1, 2, 3},
		StartTime:  "09:00",
		EndTime:    "18:00",
	}
	assert.NoError(t, req.Validate())

	req.RoomID = ""
	assert.Error(t, req.Validate())

	req.DaysOfWeek = []int{8}
	assert.Error(t, req.Validate())

	req.StartTime = ""
	assert.Error(t, req.Validate())
}

func TestToScheduleResponse(t *testing.T) {
	schedule := &domain.Schedule{
		ID:        "sch1",
		RoomID:    "room1",
		DaysMask:  domain.Monday | domain.Tuesday,
		StartTime: "09:00",
		EndTime:   "18:00",
	}
	resp := ToScheduleResponse(schedule)
	assert.Equal(t, "sch1", resp.ID)
	assert.Contains(t, resp.DaysOfWeek, 1)
	assert.Contains(t, resp.DaysOfWeek, 2)
}

func TestToScheduleResponse_Nil(t *testing.T) {
	resp := ToScheduleResponse(nil)
	assert.Nil(t, resp)
}
