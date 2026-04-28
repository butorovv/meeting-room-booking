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
	mockPool := mock_usecase.NewMockPgxPoolIface(ctrl)
	mockTx := mock_usecase.NewMockTx(ctrl)

	uc := NewBookingUseCase(mockPool, mockBookingRepo, mockRoomRepo, mockScheduleRepo, nil, nil)

	startTime := time.Now().UTC().Add(1 * time.Hour).Truncate(30 * time.Minute)
	endTime := startTime.Add(domain.SlotDuration)

	mockRoomRepo.EXPECT().ExistsByID(gomock.Any(), "room1").Return(true, nil)
	mockScheduleRepo.EXPECT().GetByRoomID(gomock.Any(), "room1").Return(&domain.Schedule{
		DaysMask:  domain.WeekdayToMask(startTime.Weekday()),
		StartTime: "00:00",
		EndTime:   "23:59",
	}, nil)

	mockPool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)
	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	mockBookingRepo.EXPECT().GetActiveBookedIntervalsWithTx(gomock.Any(), mockTx, "room1", gomock.Any(), gomock.Any()).Return(nil, nil)
	mockBookingRepo.EXPECT().CreateWithTx(gomock.Any(), mockTx, gomock.Any()).Return(nil)

	booking, err := uc.CreateBooking(context.Background(), "user1", "room1", startTime, endTime, nil)

	assert.NoError(t, err)
	assert.NotNil(t, booking)
	assert.Equal(t, "user1", booking.UserID)
}

func TestCancelBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	uc := NewBookingUseCase(nil, mockBookingRepo, nil, nil, nil, nil)

	booking := &domain.Booking{
		ID:        "booking1",
		UserID:    "user1",
		Status:    domain.BookingActive,
		StartTime: time.Now().UTC().Add(1 * time.Hour),
	}

	mockBookingRepo.EXPECT().GetByID(gomock.Any(), "booking1").Return(booking, nil)
	mockBookingRepo.EXPECT().UpdateStatus(gomock.Any(), "booking1", domain.BookingCancelled).Return(nil)

	err := uc.CancelBooking(context.Background(), "booking1", "user1", "user")
	assert.NoError(t, err)
}

func TestCancelBooking_AlreadyCancelled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	uc := NewBookingUseCase(nil, mockBookingRepo, nil, nil, nil, nil)

	booking := &domain.Booking{
		ID:     "book1",
		UserID: "user1",
		Status: domain.BookingCancelled,
	}
	mockBookingRepo.EXPECT().
		GetByID(gomock.Any(), "book1").
		Return(booking, nil)

	err := uc.CancelBooking(context.Background(), "book1", "user1", "user")
	assert.Equal(t, domain.ErrBookingAlreadyCancelled, err)
}


	updated, err := uc.CancelBooking(context.Background(), "book1", "user1")

	assert.NoError(t, err)
	assert.Equal(t, domain.BookingCancelled, updated.Status)
}

func TestCancelBooking_WrongUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockSlotRepo := mock_usecase.NewMockSlotRepositoryInterface(ctrl)

	uc := NewBookingUseCase(mockBookingRepo, mockSlotRepo)

	booking := &domain.Booking{
		ID:     "book1",
		UserID: "user2",
		Status: domain.BookingActive,
	}
	mockBookingRepo.EXPECT().
		GetByID(gomock.Any(), "book1").
		Return(booking, nil)

	_, err := uc.CancelBooking(context.Background(), "book1", "user1")

	assert.Equal(t, domain.ErrForbidden, err)
}

func TestGetUserBookings_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockSlotRepo := mock_usecase.NewMockSlotRepositoryInterface(ctrl)
	uc := NewBookingUseCase(mockBookingRepo, mockSlotRepo)

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
	mockSlotRepo := mock_usecase.NewMockSlotRepositoryInterface(ctrl)
	uc := NewBookingUseCase(mockBookingRepo, mockSlotRepo)

	expected := []*domain.Booking{{ID: "book1"}, {ID: "book2"}}
	mockBookingRepo.EXPECT().
		GetAll(gomock.Any(), 10, 0).
		Return(expected, 2, nil)

	bookings, total, err := uc.GetAllBookings(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Len(t, bookings, 2)
	assert.Equal(t, 2, total)
}
