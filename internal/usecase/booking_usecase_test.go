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
	mockSlotRepo := mock_usecase.NewMockSlotRepositoryInterface(ctrl)

	uc := NewBookingUseCase(mockBookingRepo, mockSlotRepo)

	slot := &domain.Slot{
		ID:        "slot1",
		StartTime: time.Now().UTC().Add(1 * time.Hour),
	}
	mockSlotRepo.EXPECT().
		GetByID(gomock.Any(), "slot1").
		Return(slot, nil)

	mockBookingRepo.EXPECT().
		GetBySlotID(gomock.Any(), "slot1").
		Return(nil, nil)

	mockBookingRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	booking, err := uc.CreateBooking(context.Background(), "user1", "slot1", nil)

	assert.NoError(t, err)
	assert.Equal(t, "user1", booking.UserID)
	assert.Equal(t, domain.BookingActive, booking.Status)
}

func TestCreateBooking_SlotInPast(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockSlotRepo := mock_usecase.NewMockSlotRepositoryInterface(ctrl)

	uc := NewBookingUseCase(mockBookingRepo, mockSlotRepo)

	slot := &domain.Slot{
		ID:        "slot1",
		StartTime: time.Now().UTC().Add(-1 * time.Hour),
	}
	mockSlotRepo.EXPECT().
		GetByID(gomock.Any(), "slot1").
		Return(slot, nil)

	_, err := uc.CreateBooking(context.Background(), "user1", "slot1", nil)

	assert.Equal(t, domain.ErrSlotInPast, err)
}

func TestCreateBooking_SlotAlreadyBooked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockSlotRepo := mock_usecase.NewMockSlotRepositoryInterface(ctrl)

	uc := NewBookingUseCase(mockBookingRepo, mockSlotRepo)

	slot := &domain.Slot{
		ID:        "slot1",
		StartTime: time.Now().UTC().Add(1 * time.Hour),
	}
	existingBooking := &domain.Booking{Status: domain.BookingActive}

	mockSlotRepo.EXPECT().
		GetByID(gomock.Any(), "slot1").
		Return(slot, nil)

	mockBookingRepo.EXPECT().
		GetBySlotID(gomock.Any(), "slot1").
		Return(existingBooking, nil)

	_, err := uc.CreateBooking(context.Background(), "user1", "slot1", nil)

	assert.Equal(t, domain.ErrSlotAlreadyBooked, err)
}

func TestCancelBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockSlotRepo := mock_usecase.NewMockSlotRepositoryInterface(ctrl)

	uc := NewBookingUseCase(mockBookingRepo, mockSlotRepo)

	booking := &domain.Booking{
		ID:     "book1",
		UserID: "user1",
		Status: domain.BookingActive,
	}
	mockBookingRepo.EXPECT().
		GetByID(gomock.Any(), "book1").
		Return(booking, nil)

	mockBookingRepo.EXPECT().
		UpdateStatus(gomock.Any(), "book1", string(domain.BookingCancelled)).
		Return(nil)

	updated, err := uc.CancelBooking(context.Background(), "book1", "user1")

	assert.NoError(t, err)
	assert.Equal(t, domain.BookingCancelled, updated.Status)
}

func TestCancelBooking_AlreadyCancelled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock_usecase.NewMockBookingRepositoryInterface(ctrl)
	mockSlotRepo := mock_usecase.NewMockSlotRepositoryInterface(ctrl)

	uc := NewBookingUseCase(mockBookingRepo, mockSlotRepo)

	booking := &domain.Booking{
		ID:     "book1",
		UserID: "user1",
		Status: domain.BookingCancelled,
	}
	mockBookingRepo.EXPECT().
		GetByID(gomock.Any(), "book1").
		Return(booking, nil)

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
