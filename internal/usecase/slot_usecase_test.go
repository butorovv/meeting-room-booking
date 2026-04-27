package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"

	"github.com/butorovv/meeting-room-booking/internal/domain"
	mock_usecase "github.com/butorovv/meeting-room-booking/internal/usecase/mock"
)

func TestGetAvailableSlots_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	uc := NewSlotUseCase(mockScheduleRepo, mockBookingRepo)

	schedule := &domain.Schedule{
		RoomID:    "room1",
		DaysMask:  domain.Monday,
		StartTime: "09:00",
		EndTime:   "10:00",
	}
	dayStart := time.Date(2024, 4, 8, 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.AddDate(0, 0, 1)
	bookedStart := time.Date(2024, 4, 8, 9, 30, 0, 0, time.UTC)
	bookedEnd := bookedStart.Add(domain.SlotDuration)

	mockScheduleRepo.EXPECT().
		GetByRoomID(gomock.Any(), "room1").
		Return(schedule, nil)

	mockBookingRepo.EXPECT().
		GetActiveBookedIntervals(gomock.Any(), "room1", dayStart, dayEnd).
		Return([]domain.BookedInterval{{StartTime: bookedStart, EndTime: bookedEnd}}, nil)

	slots, err := uc.GetAvailableSlots(context.Background(), "room1", "2024-04-08")

	assert.NoError(t, err)
	assert.Len(t, slots, 2)
	assert.False(t, slots[0].IsBooked)
	assert.True(t, slots[1].IsBooked)
}

func TestGetAvailableSlots_NoSchedule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	uc := NewSlotUseCase(mockScheduleRepo, mockBookingRepo)

	mockScheduleRepo.EXPECT().
		GetByRoomID(gomock.Any(), "room1").
		Return(nil, pgx.ErrNoRows)

	slots, err := uc.GetAvailableSlots(context.Background(), "room1", "2024-04-08")

	assert.NoError(t, err)
	assert.Empty(t, slots)
}
