package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

type ScheduleRepositoryInterface interface {
	Create(ctx context.Context, schedule *domain.Schedule) error
	GetByRoomID(ctx context.Context, roomID string) (*domain.Schedule, error)
	ExistsByRoomID(ctx context.Context, roomID string) (bool, error)
}

type ScheduleUseCase struct {
	scheduleRepo ScheduleRepositoryInterface
}

func NewScheduleUseCase(scheduleRepo ScheduleRepositoryInterface) *ScheduleUseCase {
	return &ScheduleUseCase{
		scheduleRepo: scheduleRepo,
	}
}

func (uc *ScheduleUseCase) CreateSchedule(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
	exists, err := uc.scheduleRepo.ExistsByRoomID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrScheduleExists
	}

	schedule := &domain.Schedule{
		ID:        uuid.NewString(),
		RoomID:    roomID,
		DaysMask:  domain.DaysToMask(daysOfWeek),
		StartTime: startTime,
		EndTime:   endTime,
		CreatedAt: time.Now().UTC(),
	}
	err = uc.scheduleRepo.Create(ctx, schedule)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}
