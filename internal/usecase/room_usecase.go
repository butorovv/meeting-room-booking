package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

type RoomRepositoryInterface interface {
	Create(ctx context.Context, room *domain.Room) error
	List(ctx context.Context) ([]*domain.Room, error)
}

type RoomUseCase struct {
	repo RoomRepositoryInterface
}

func NewRoomUseCase(repo RoomRepositoryInterface) *RoomUseCase {
	return &RoomUseCase{repo: repo}
}

func (uc *RoomUseCase) CreateRoom(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error) {
	room := &domain.Room{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		Capacity:    capacity,
		CreatedAt:   time.Now().UTC(),
	}
	err := uc.repo.Create(ctx, room)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (uc *RoomUseCase) ListRooms(ctx context.Context) ([]*domain.Room, error) {
	return uc.repo.List(ctx)
}
