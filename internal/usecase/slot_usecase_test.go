package usecase

import (
	"context"
	"encoding/json"
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
	uc := NewSlotUseCase(mockScheduleRepo, mockBookingRepo, nil)

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
	uc := NewSlotUseCase(mockScheduleRepo, mockBookingRepo, nil)

	mockScheduleRepo.EXPECT().
		GetByRoomID(gomock.Any(), "room1").
		Return(nil, pgx.ErrNoRows)

	slots, err := uc.GetAvailableSlots(context.Background(), "room1", "2024-04-08")

	assert.NoError(t, err)
	assert.Empty(t, slots)
}

func TestGetAvailableSlots_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	cache := newTestCache()
	uc := NewSlotUseCase(mockScheduleRepo, mockBookingRepo, cache)

	expected := []*domain.Slot{{
		ID:        "slot1",
		RoomID:    "room1",
		StartTime: time.Date(2024, 4, 8, 9, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2024, 4, 8, 9, 30, 0, 0, time.UTC),
	}}
	value, err := json.Marshal(expected)
	assert.NoError(t, err)
	cache.values["slots:v1:room1:2024-04-08"] = value

	slots, err := uc.GetAvailableSlots(context.Background(), "room1", "2024-04-08")

	assert.NoError(t, err)
	assert.Equal(t, expected, slots)
}

func TestGetAvailableSlots_InvalidCacheFallsBack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	cache := newTestCache()
	uc := NewSlotUseCase(mockScheduleRepo, mockBookingRepo, cache)

	key := "slots:v1:room1:2024-04-08"
	cache.values[key] = []byte("{")
	schedule := &domain.Schedule{
		RoomID:    "room1",
		DaysMask:  domain.Monday,
		StartTime: "09:00",
		EndTime:   "09:30",
	}
	dayStart := time.Date(2024, 4, 8, 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.AddDate(0, 0, 1)

	mockScheduleRepo.EXPECT().
		GetByRoomID(gomock.Any(), "room1").
		Return(schedule, nil)

	mockBookingRepo.EXPECT().
		GetActiveBookedIntervals(gomock.Any(), "room1", dayStart, dayEnd).
		Return([]domain.BookedInterval{}, nil)

	slots, err := uc.GetAvailableSlots(context.Background(), "room1", "2024-04-08")

	assert.NoError(t, err)
	assert.Len(t, slots, 1)
	assert.Contains(t, cache.deleted, key)
	assert.Contains(t, cache.values, key)
}

type testCache struct {
	values  map[string][]byte
	deleted []string
}

func newTestCache() *testCache {
	return &testCache{values: make(map[string][]byte)}
}

func (c *testCache) Get(_ context.Context, key string) ([]byte, error) {
	return c.values[key], nil
}

func (c *testCache) Set(_ context.Context, key string, value []byte, _ time.Duration) error {
	c.values[key] = value
	return nil
}

func (c *testCache) Del(_ context.Context, key string) error {
	delete(c.values, key)
	c.deleted = append(c.deleted, key)
	return nil
}

func (c *testCache) Ping(_ context.Context) error {
	return nil
}

func (c *testCache) Close() error {
	return nil
}
