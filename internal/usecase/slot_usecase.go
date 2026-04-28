package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/butorovv/meeting-room-booking/internal/domain"
	slotgenerator "github.com/butorovv/meeting-room-booking/pkg/slot"
)

type SlotUseCase struct {
	scheduleRepo ScheduleRepositoryInterface
	bookingRepo  BookingRepositoryInterface
}

func NewSlotUseCase(scheduleRepo ScheduleRepositoryInterface, bookingRepo BookingRepositoryInterface) *SlotUseCase {
	return &SlotUseCase{
		scheduleRepo: scheduleRepo,
		bookingRepo:  bookingRepo,
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

	schedule, err := uc.scheduleRepo.GetByRoomID(ctx, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*domain.Slot{}, nil
		}
		return nil, err
	}

	if schedule.DaysMask&domain.WeekdayToMask(day.Weekday()) == 0 {
		return []*domain.Slot{}, nil
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

	return slots, nil
}
