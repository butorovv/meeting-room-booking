package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/butorovv/meeting-room-booking/internal/domain"
	mock_usecase "github.com/butorovv/meeting-room-booking/internal/usecase/mock"
)

func TestGetAvailableSlots_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSlotRepo := mock_usecase.NewMockSlotRepositoryInterface(ctrl)
	uc := NewSlotUseCase(mockSlotRepo)

	mockSlotRepo.EXPECT().
		GetAvailableByRoomAndDate(gomock.Any(), "room1", "2024-04-07").
		Return([]*domain.Slot{{ID: "slot1"}, {ID: "slot2"}}, nil)

	slots, err := uc.GetAvailableSlots(context.Background(), "room1", "2024-04-07")

	assert.NoError(t, err)
	assert.Len(t, slots, 2)
}

func TestGetAvailableSlots_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSlotRepo := mock_usecase.NewMockSlotRepositoryInterface(ctrl)
	uc := NewSlotUseCase(mockSlotRepo)

	mockSlotRepo.EXPECT().
		GetAvailableByRoomAndDate(gomock.Any(), "room1", "2024-04-07").
		Return([]*domain.Slot{}, nil)

	slots, err := uc.GetAvailableSlots(context.Background(), "room1", "2024-04-07")

	assert.NoError(t, err)
	assert.Empty(t, slots)
}
