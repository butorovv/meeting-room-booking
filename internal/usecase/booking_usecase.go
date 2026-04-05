package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

type BookingRepositoryInterface interface {
	Create(ctx context.Context, booking *domain.Booking) error
	GetBySlotID(ctx context.Context, slotID string) (*domain.Booking, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Booking, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Booking, int, error)
	UpdateStatus(ctx context.Context, id, status string) error
	GetByID(ctx context.Context, id string) (*domain.Booking, error)
}

type BookingUseCase struct {
	bookingRepo BookingRepositoryInterface
	slotRepo    SlotRepositoryInterface
}

func NewBookingUseCase(bookingRepo BookingRepositoryInterface, slotRepo SlotRepositoryInterface) *BookingUseCase {
	return &BookingUseCase{
		bookingRepo: bookingRepo,
		slotRepo:    slotRepo,
	}
}

func (uc *BookingUseCase) CreateBooking(ctx context.Context, userID, slotID string, conferenceLink *string) (*domain.Booking, error) {
	slot, err := uc.slotRepo.GetByID(ctx, slotID)
	if err != nil {
		return nil, err
	}
	if slot.IsInPast() {
		return nil, domain.ErrSlotInPast
	}

	existing, _ := uc.bookingRepo.GetBySlotID(ctx, slotID)
	if existing != nil && existing.Status == domain.BookingActive {
		return nil, domain.ErrSlotAlreadyBooked
	}

	booking := &domain.Booking{
		ID:             uuid.NewString(),
		SlotID:         slotID,
		UserID:         userID,
		Status:         domain.BookingActive,
		ConferenceLink: conferenceLink,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
	err = uc.bookingRepo.Create(ctx, booking)
	if err != nil {
		return nil, err
	}
	return booking, nil
}

func (uc *BookingUseCase) CancelBooking(ctx context.Context, bookingID, userID string) (*domain.Booking, error) {
	booking, err := uc.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if booking.UserID != userID {
		return nil, domain.ErrForbidden
	}

	// если уже отменена — возвращаем ее (идемпотентность)
	if booking.Status == domain.BookingCancelled {
		return booking, nil
	}

	// отменяем бронь
	err = uc.bookingRepo.UpdateStatus(ctx, bookingID, string(domain.BookingCancelled))
	if err != nil {
		return nil, err
	}

	// возвращаем обновленную бронь
	booking.Status = domain.BookingCancelled
	return booking, nil
}

func (uc *BookingUseCase) GetUserBookings(ctx context.Context, userID string) ([]*domain.Booking, error) {
	return uc.bookingRepo.GetByUserID(ctx, userID)
}

func (uc *BookingUseCase) GetAllBookings(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error) {
	offset := (page - 1) * pageSize
	return uc.bookingRepo.GetAll(ctx, pageSize, offset)
}
