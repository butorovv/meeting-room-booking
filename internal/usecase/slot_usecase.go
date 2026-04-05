package usecase

import (
	"context"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

type SlotUseCase struct {
	repo SlotRepositoryInterface
}

func NewSlotUseCase(repo SlotRepositoryInterface) *SlotUseCase {
	return &SlotUseCase{repo: repo}
}

func (uc *SlotUseCase) GetAvailableSlots(ctx context.Context, roomID, date string) ([]*domain.Slot, error) {
	return uc.repo.GetAvailableByRoomAndDate(ctx, roomID, date)
}
