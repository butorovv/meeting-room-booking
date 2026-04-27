package repository

import (
	"context"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

type SlotRepository struct {
	db PgxIface
}

func NewSlotRepository(db PgxIface) *SlotRepository {
	return &SlotRepository{db: db}
}

func (r *SlotRepository) BatchCreate(ctx context.Context, slots []*domain.Slot) error {
	return nil
}

func (r *SlotRepository) GetAvailableByRoomAndDate(ctx context.Context, roomID string, date string) ([]*domain.Slot, error) {
	return []*domain.Slot{}, nil
}

func (r *SlotRepository) GetByID(ctx context.Context, id string) (*domain.Slot, error) {
	return nil, domain.ErrSlotNotFound
}
