package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/butorovv/meeting-room-booking/internal/domain"
	"github.com/butorovv/meeting-room-booking/pkg/cache"
	"github.com/butorovv/meeting-room-booking/pkg/logger"
	slotgenerator "github.com/butorovv/meeting-room-booking/pkg/slot"
)

const slotsCacheTTL = 5 * time.Minute

type SlotUseCase struct {
	scheduleRepo ScheduleRepositoryInterface
	bookingRepo  BookingRepositoryInterface
	cache        cache.Cache
}

func NewSlotUseCase(scheduleRepo ScheduleRepositoryInterface, bookingRepo BookingRepositoryInterface, cache cache.Cache) *SlotUseCase {
	return &SlotUseCase{
		scheduleRepo: scheduleRepo,
		bookingRepo:  bookingRepo,
		cache:        cache,
	}
}

func (uc *SlotUseCase) GetAvailableSlots(ctx context.Context, roomID, date string) ([]*domain.Slot, error) {
	return uc.GetSlots(ctx, roomID, date)
}

func (uc *SlotUseCase) GetSlots(ctx context.Context, roomID, date string) ([]*domain.Slot, error) {
	day, err := time.ParseInLocation("2006-01-02", date, time.UTC)
	if err != nil {
		return nil, err
	}
	day = day.UTC()

	key := slotsCacheKey(roomID, day)
	if slots, ok := uc.getCachedSlots(ctx, key); ok {
		return slots, nil
	}

	schedule, err := uc.scheduleRepo.GetByRoomID(ctx, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slots := []*domain.Slot{}
			uc.setCachedSlots(ctx, key, slots)
			return slots, nil
		}
		return nil, err
	}

	if schedule.DaysMask&domain.WeekdayToMask(day.Weekday()) == 0 {
		slots := []*domain.Slot{}
		uc.setCachedSlots(ctx, key, slots)
		return slots, nil
	}

	slots := slotgenerator.GenerateSlots(schedule, day)

	dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.AddDate(0, 0, 1)

	bookedIntervals, err := uc.bookingRepo.GetActiveBookedIntervals(ctx, roomID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}

	for i := range slots {
		for _, booked := range bookedIntervals {
			if slots[i].StartTime.Equal(booked.StartTime) && slots[i].EndTime.Equal(booked.EndTime) {
				slots[i].IsBooked = true
				slots[i].BookingID = booked.ID
				break
			}
		}
	}

	uc.setCachedSlots(ctx, key, slots)

	return slots, nil
}

func (uc *SlotUseCase) getCachedSlots(ctx context.Context, key string) ([]*domain.Slot, bool) {
	if uc.cache == nil {
		return nil, false
	}

	value, err := uc.cache.Get(ctx, key)
	if err != nil {
		logger.Global().WarnContext(ctx, "redis get failed", "key", key, "err", err)
		return nil, false
	}
	if len(value) == 0 {
		return nil, false
	}

	var slots []*domain.Slot
	if err := json.Unmarshal(value, &slots); err != nil {
		logger.Global().WarnContext(ctx, "redis invalid slots json", "key", key, "err", err)
		if err := uc.cache.Del(ctx, key); err != nil {
			logger.Global().WarnContext(ctx, "redis del failed", "key", key, "err", err)
		}
		return nil, false
	}

	return slots, true
}

func (uc *SlotUseCase) setCachedSlots(ctx context.Context, key string, slots []*domain.Slot) {
	if uc.cache == nil {
		logger.Global().Warn("cache is nil, skipping set", "key", key)
		return
	}

	logger.Global().Info("saving slots to cache", "key", key)

	value, err := json.Marshal(slots)
	if err != nil {
		logger.Global().WarnContext(ctx, "slots cache marshal failed", "key", key, "err", err)
		return
	}

	if err := uc.cache.Set(ctx, key, value, slotsCacheTTL); err != nil {
		logger.Global().WarnContext(ctx, "redis set failed", "key", key, "err", err)
	}
}

func slotsCacheKey(roomID string, date time.Time) string {
	day := date.UTC().Format("2006-01-02")
	return fmt.Sprintf("slots:v1:%s:%s", roomID, day)
}
