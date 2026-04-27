package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/butorovv/meeting-room-booking/internal/domain"
	mock_usecase "github.com/butorovv/meeting-room-booking/internal/usecase/mock"
)

func TestCreateBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockRoomRepo := mock_usecase.NewMockRoomRepositoryInterface(ctrl)
	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	mockDB := mock_usecase.NewMockPgxPoolIface(ctrl)
	mockTx := mock_usecase.NewMockTx(ctrl)
	uc := NewBookingUseCase(mockDB, mockBookingRepo, mockRoomRepo, mockScheduleRepo)

	startTime := nextBookingStart()
	endTime := startTime.Add(domain.SlotDuration)
	dayStart := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.AddDate(0, 0, 1)

	mockDB.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	mockRoomRepo.EXPECT().
		ExistsByID(gomock.Any(), "room1").
		Return(true, nil)

	mockScheduleRepo.EXPECT().
		GetByRoomID(gomock.Any(), "room1").
		Return(scheduleFor(startTime), nil)

	mockBookingRepo.EXPECT().
		GetActiveBookedIntervals(gomock.Any(), "room1", dayStart, dayEnd).
		Return([]domain.BookedInterval{}, nil)

	mockBookingRepo.EXPECT().
		CreateWithTx(gomock.Any(), mockTx, gomock.Any()).
		Return(nil)

	booking, err := uc.CreateBooking(context.Background(), "user1", "room1", startTime, endTime, nil)

	assert.NoError(t, err)
	assert.Nil(t, booking.SlotID)
	assert.Equal(t, "room1", booking.RoomID)
	assert.Equal(t, "user1", booking.UserID)
	assert.Equal(t, domain.BookingActive, booking.Status)
	assert.Equal(t, startTime, booking.StartTime)
	assert.Equal(t, endTime, booking.EndTime)
}

func TestCreateBooking_SlotInPast(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockRoomRepo := mock_usecase.NewMockRoomRepositoryInterface(ctrl)
	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	mockDB := mock_usecase.NewMockPgxPoolIface(ctrl)
	uc := NewBookingUseCase(mockDB, mockBookingRepo, mockRoomRepo, mockScheduleRepo)

	startTime := time.Now().UTC().Add(-1 * time.Hour)
	endTime := startTime.Add(domain.SlotDuration)

	_, err := uc.CreateBooking(context.Background(), "user1", "room1", startTime, endTime, nil)

	assert.Equal(t, domain.ErrSlotInPast, err)
}

func TestCreateBooking_SlotAlreadyBooked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockRoomRepo := mock_usecase.NewMockRoomRepositoryInterface(ctrl)
	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	mockDB := mock_usecase.NewMockPgxPoolIface(ctrl)
	mockTx := mock_usecase.NewMockTx(ctrl)
	uc := NewBookingUseCase(mockDB, mockBookingRepo, mockRoomRepo, mockScheduleRepo)

	startTime := nextBookingStart()
	endTime := startTime.Add(domain.SlotDuration)
	dayStart := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.AddDate(0, 0, 1)

	mockDB.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	mockRoomRepo.EXPECT().
		ExistsByID(gomock.Any(), "room1").
		Return(true, nil)

	mockScheduleRepo.EXPECT().
		GetByRoomID(gomock.Any(), "room1").
		Return(scheduleFor(startTime), nil)

	mockBookingRepo.EXPECT().
		GetActiveBookedIntervals(gomock.Any(), "room1", dayStart, dayEnd).
		Return([]domain.BookedInterval{{StartTime: startTime, EndTime: endTime}}, nil)

	_, err := uc.CreateBooking(context.Background(), "user1", "room1", startTime, endTime, nil)

	assert.Equal(t, domain.ErrSlotAlreadyBooked, err)
}

func TestCreateBooking_RoomNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockRoomRepo := mock_usecase.NewMockRoomRepositoryInterface(ctrl)
	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	mockDB := mock_usecase.NewMockPgxPoolIface(ctrl)
	uc := NewBookingUseCase(mockDB, mockBookingRepo, mockRoomRepo, mockScheduleRepo)

	startTime := nextBookingStart()
	endTime := startTime.Add(domain.SlotDuration)
	dayStart := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.AddDate(0, 0, 1)

	mockDB.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	mockRoomRepo.EXPECT().
		ExistsByID(gomock.Any(), "room1").
		Return(false, nil)

	mockScheduleRepo.EXPECT().
		GetByRoomID(gomock.Any(), "room1").
		Return(nil, nil)

	mockBookingRepo.EXPECT().
		GetActiveBookedIntervals(gomock.Any(), "room1", dayStart, dayEnd).
		Return([]domain.BookedInterval{}, nil)

	_, err := uc.CreateBooking(context.Background(), "user1", "room1", startTime, endTime, nil)

	assert.Equal(t, domain.ErrRoomNotFound, err)
}

func TestCreateBooking_ScheduleNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockRoomRepo := mock_usecase.NewMockRoomRepositoryInterface(ctrl)
	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	mockDB := mock_usecase.NewMockPgxPoolIface(ctrl)
	uc := NewBookingUseCase(mockDB, mockBookingRepo, mockRoomRepo, mockScheduleRepo)

	startTime := nextBookingStart()
	endTime := startTime.Add(domain.SlotDuration)
	dayStart := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.AddDate(0, 0, 1)

	mockDB.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	mockRoomRepo.EXPECT().
		ExistsByID(gomock.Any(), "room1").
		Return(true, nil)

	mockScheduleRepo.EXPECT().
		GetByRoomID(gomock.Any(), "room1").
		Return(nil, nil)

	mockBookingRepo.EXPECT().
		GetActiveBookedIntervals(gomock.Any(), "room1", dayStart, dayEnd).
		Return([]domain.BookedInterval{}, nil)

	_, err := uc.CreateBooking(context.Background(), "user1", "room1", startTime, endTime, nil)

	assert.Equal(t, domain.ErrScheduleNotFound, err)
}

func TestCreateBooking_NotInSchedule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockRoomRepo := mock_usecase.NewMockRoomRepositoryInterface(ctrl)
	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	mockDB := mock_usecase.NewMockPgxPoolIface(ctrl)
	uc := NewBookingUseCase(mockDB, mockBookingRepo, mockRoomRepo, mockScheduleRepo)

	startTime := nextBookingStart()
	endTime := startTime.Add(domain.SlotDuration)
	dayStart := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.AddDate(0, 0, 1)

	mockDB.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	mockRoomRepo.EXPECT().
		ExistsByID(gomock.Any(), "room1").
		Return(true, nil)

	mockScheduleRepo.EXPECT().
		GetByRoomID(gomock.Any(), "room1").
		Return(scheduleFor(startTime), nil)

	mockBookingRepo.EXPECT().
		GetActiveBookedIntervals(gomock.Any(), "room1", dayStart, dayEnd).
		Return([]domain.BookedInterval{{StartTime: startTime, EndTime: endTime}}, nil)

	_, err := uc.CreateBooking(context.Background(), "user1", "room1", startTime, endTime, nil)

	assert.Equal(t, domain.ErrNotInSchedule, err)
}

func TestCancelBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockRoomRepo := mock_usecase.NewMockRoomRepositoryInterface(ctrl)
	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	mockDB := mock_usecase.NewMockPgxPoolIface(ctrl)
	mockTx := mock_usecase.NewMockTx(ctrl)
	uc := NewBookingUseCase(mockDB, mockBookingRepo, mockRoomRepo, mockScheduleRepo)

	booking := &domain.Booking{
		ID:     "booking1",
		UserID: "user1",
		Status: domain.BookingActive,
	}

	mockDB.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	mockBookingRepo.EXPECT().
		GetByID(gomock.Any(), "booking1").
		Return(booking, nil)

	mockBookingRepo.EXPECT().
		UpdateStatusWithTx(gomock.Any(), mockTx, "booking1", string(domain.BookingCancelled)).
		Return(nil)

	updatedBooking, err := uc.CancelBooking(context.Background(), "booking1", "user1")

	assert.NoError(t, err)
	assert.Equal(t, domain.BookingCancelled, updatedBooking.Status)
}

func TestCancelBooking_WrongUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockRoomRepo := mock_usecase.NewMockRoomRepositoryInterface(ctrl)
	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	mockDB := mock_usecase.NewMockPgxPoolIface(ctrl)
	mockTx := mock_usecase.NewMockTx(ctrl)
	uc := NewBookingUseCase(mockDB, mockBookingRepo, mockRoomRepo, mockScheduleRepo)

	booking := &domain.Booking{
		ID:     "booking1",
		UserID: "user2",
		Status: domain.BookingActive,
	}

	mockDB.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	mockBookingRepo.EXPECT().
		GetByID(gomock.Any(), "booking1").
		Return(booking, nil)

	_, err := uc.CancelBooking(context.Background(), "booking1", "user1")

	assert.Equal(t, domain.ErrForbidden, err)
}

func TestCancelBooking_AlreadyCancelled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockRoomRepo := mock_usecase.NewMockRoomRepositoryInterface(ctrl)
	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	mockDB := mock_usecase.NewMockPgxPoolIface(ctrl)
	mockTx := mock_usecase.NewMockTx(ctrl)
	uc := NewBookingUseCase(mockDB, mockBookingRepo, mockRoomRepo, mockScheduleRepo)

	booking := &domain.Booking{
		ID:     "booking1",
		UserID: "user1",
		Status: domain.BookingCancelled,
	}

	mockDB.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	mockBookingRepo.EXPECT().
		GetByID(gomock.Any(), "booking1").
		Return(booking, nil)

	mockBookingRepo.EXPECT().
		UpdateStatusWithTx(gomock.Any(), mockTx, "booking1", string(domain.BookingCancelled)).
		Return(nil)

	updatedBooking, err := uc.CancelBooking(context.Background(), "booking1", "user1")

	assert.NoError(t, err)
	assert.Equal(t, domain.BookingCancelled, updatedBooking.Status)
}

func TestGetUserBookings_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockRoomRepo := mock_usecase.NewMockRoomRepositoryInterface(ctrl)
	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	uc := NewBookingUseCase(mockBookingRepo, mockRoomRepo, mockScheduleRepo)

	expected := []*domain.Booking{{ID: "book1", UserID: "user1"}}
	mockBookingRepo.EXPECT().
		GetByUserID(gomock.Any(), "user1").
		Return(expected, nil)

	bookings, err := uc.GetUserBookings(context.Background(), "user1")
	assert.NoError(t, err)
	assert.Len(t, bookings, 1)
}

func TestGetAllBookings_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockRoomRepo := mock_usecase.NewMockRoomRepositoryInterface(ctrl)
	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	uc := NewBookingUseCase(mockBookingRepo, mockRoomRepo, mockScheduleRepo)

	expected := []*domain.Booking{{ID: "book1"}, {ID: "book2"}}
	mockBookingRepo.EXPECT().
		GetAll(gomock.Any(), 10, 0).
		Return(expected, 2, nil)

	bookings, total, err := uc.GetAllBookings(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Len(t, bookings, 2)
	assert.Equal(t, 2, total)
}

func nextBookingStart() time.Time {
	now := time.Now().UTC()
	candidate := now.AddDate(0, 0, 1)
	start := time.Date(candidate.Year(), candidate.Month(), candidate.Day(), 9, 0, 0, 0, time.UTC)
	if !start.After(now) {
		candidate = candidate.AddDate(0, 0, 1)
		start = time.Date(candidate.Year(), candidate.Month(), candidate.Day(), 9, 0, 0, 0, time.UTC)
	}
	return start
}

func scheduleFor(startTime time.Time) *domain.Schedule {
	return &domain.Schedule{
		RoomID:    "room1",
		DaysMask:  domain.WeekdayToMask(startTime.Weekday()),
		StartTime: "09:00",
		EndTime:   "18:00",
	}
}
