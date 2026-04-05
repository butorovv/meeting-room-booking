package usecase

import (
	"context"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

type SlotRepositoryInterface interface {
	GetByID(ctx context.Context, id string) (*domain.Slot, error)
	BatchCreate(ctx context.Context, slots []*domain.Slot) error
	GetAvailableByRoomAndDate(ctx context.Context, roomID, date string) ([]*domain.Slot, error)
}
