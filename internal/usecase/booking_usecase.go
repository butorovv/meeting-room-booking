package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/butorovv/meeting-room-booking/internal/domain"
	"github.com/butorovv/meeting-room-booking/pkg/cache"
	"github.com/butorovv/meeting-room-booking/pkg/logger"
	"github.com/butorovv/meeting-room-booking/pkg/worker"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type BookingRepositoryInterface interface {
	Create(ctx context.Context, booking *domain.Booking) error
	CreateWithTx(ctx context.Context, tx pgx.Tx, booking *domain.Booking) error
	GetBySlotID(ctx context.Context, slotID string) (*domain.Booking, error)
	GetActiveBookedIntervals(ctx context.Context, roomID string, startDate, endDate time.Time) ([]domain.BookedInterval, error)
	GetActiveBookedIntervalsWithTx(ctx context.Context, tx pgx.Tx, roomID string, startDate, endDate time.Time) ([]domain.BookedInterval, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Booking, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Booking, int, error)
	UpdateStatus(ctx context.Context, id, status string) error
	UpdateStatusWithTx(ctx context.Context, tx pgx.Tx, id, status string) error
	GetByID(ctx context.Context, id string) (*domain.Booking, error)
}

type BookingUseCase struct {
	db           PgxPoolIface
	bookingRepo  BookingRepositoryInterface
	roomRepo     RoomRepositoryInterface
	scheduleRepo ScheduleRepositoryInterface
	cache        cache.Cache
	workerPool   *worker.WorkerPool
}

func NewBookingUseCase(
	db PgxPoolIface,
	bookingRepo BookingRepositoryInterface,
	roomRepo RoomRepositoryInterface,
	scheduleRepo ScheduleRepositoryInterface,
	cache cache.Cache,
	workerPool *worker.WorkerPool,
) *BookingUseCase {
	return &BookingUseCase{
		db:           db,
		bookingRepo:  bookingRepo,
		roomRepo:     roomRepo,
		scheduleRepo: scheduleRepo,
		cache:        cache,
		workerPool:   workerPool,
	}
}

func (uc *BookingUseCase) CreateBooking(ctx context.Context, userID, roomID string, startTime, endTime time.Time, conferenceLink *string) (*domain.Booking, error) {
	startTime = startTime.UTC()
	endTime = endTime.UTC()

	// 1. валидация интервала
	if err := validateBookingInterval(startTime, endTime); err != nil {
		return nil, err
	}

	// 2. проверка существования комнаты
	roomExists, err := uc.roomRepo.ExistsByID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if !roomExists {
		return nil, domain.ErrRoomNotFound
	}

	// 3. проверка расписания
	schedule, err := uc.scheduleRepo.GetByRoomID(ctx, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrScheduleNotFound
		}
		return nil, err
	}
	if !isInsideSchedule(schedule, startTime, endTime) {
		return nil, domain.ErrInvalidBookingTime
	}

	// 4. начало транзакции
	tx, err := uc.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 5. проверка на пересечение с существующими бронями (внутри транзакции)
	dayStart := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.AddDate(0, 0, 1)

	bookedIntervals, err := uc.bookingRepo.GetActiveBookedIntervalsWithTx(ctx, tx, roomID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}

	for _, interval := range bookedIntervals {
		if intervalsOverlap(startTime, endTime, interval.StartTime.UTC(), interval.EndTime.UTC()) {
			return nil, domain.ErrSlotAlreadyBooked
		}
	}

	// 6. создаем бронь
	booking := &domain.Booking{
		ID:             uuid.NewString(),
		RoomID:         roomID,
		StartTime:      startTime,
		EndTime:        endTime,
		UserID:         userID,
		Status:         domain.BookingActive,
		ConferenceLink: conferenceLink,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	// 7. вставка через транзакцию
	err = uc.bookingRepo.CreateWithTx(ctx, tx, booking)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrSlotAlreadyBooked
		}
		return nil, err
	}

	// 8. коммит транзакции
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	uc.invalidateSlotsCache(ctx, booking.RoomID, booking.StartTime)

	// отправляем уведомление в воркер пул
	if uc.workerPool != nil {
		job := worker.NotificationJob{
			UserID:    userID,
			UserEmail: userID + "@example.com",
			StartTime: booking.StartTime,
			BookingID: booking.ID,
		}
		if !uc.workerPool.Add(job) {
			logger.FromContext(ctx).Warn("failed to add notification job to worker pool", "booking_id", booking.ID)
		}
	}

	return booking, nil
}

func (uc *BookingUseCase) CancelBooking(ctx context.Context, bookingID, userID string) (*domain.Booking, error) {
	// 1. получаем бронь
	booking, err := uc.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	// 2. проверяем владельца
	if booking.UserID != userID {
		return nil, domain.ErrForbidden
	}

	// 3. если уже отменена — идемпотентность
	if booking.Status == domain.BookingCancelled {
		uc.invalidateSlotsCache(ctx, booking.RoomID, booking.StartTime)
		return booking, nil
	}

	// 4. начало транзакции
	tx, err := uc.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 5. обновляем статус через транзакцию
	err = uc.bookingRepo.UpdateStatusWithTx(ctx, tx, bookingID, string(domain.BookingCancelled))
	if err != nil {
		return nil, err
	}

	// 6. коммит
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	booking.Status = domain.BookingCancelled
	uc.invalidateSlotsCache(ctx, booking.RoomID, booking.StartTime)
	return booking, nil
}

func (uc *BookingUseCase) invalidateSlotsCache(ctx context.Context, roomID string, startTime time.Time) {
	if uc.cache == nil || roomID == "" || startTime.IsZero() {
		return
	}

	key := slotsCacheKey(roomID, startTime.UTC())
	if err := uc.cache.Del(ctx, key); err != nil {
		logger.Global().WarnContext(ctx, "redis del failed", "key", key, "err", err)
	}
}

func (uc *BookingUseCase) enqueueNotification(ctx context.Context, booking *domain.Booking) {
	if uc.workerPool == nil || booking == nil {
		return
	}

	ok := uc.workerPool.Add(worker.NotificationJob{
		UserID:    booking.UserID,
		UserEmail: booking.UserID + "@example.com",
		StartTime: booking.StartTime,
		BookingID: booking.ID,
	})
	if !ok {
		logger.Global().WarnContext(ctx, "notification job was dropped", "bookingID", booking.ID, "userID", booking.UserID)
	}
}

func validateBookingInterval(startTime, endTime time.Time) error {
	if !startTime.After(time.Now().UTC()) {
		return domain.ErrSlotInPast
	}
	if endTime.Sub(startTime) != domain.SlotDuration {
		return domain.ErrInvalidBookingTime
	}
	if startTime.Minute()%30 != 0 || startTime.Second() != 0 || startTime.Nanosecond() != 0 {
		return domain.ErrInvalidBookingTime
	}
	if !endTime.Equal(startTime.Add(domain.SlotDuration)) {
		return domain.ErrInvalidBookingTime
	}

	return nil
}

func isInsideSchedule(schedule *domain.Schedule, startTime, endTime time.Time) bool {
	if schedule == nil || schedule.DaysMask&domain.WeekdayToMask(startTime.Weekday()) == 0 {
		return false
	}

	startClock, ok := parseBookingClock(schedule.StartTime)
	if !ok {
		return false
	}
	endClock, ok := parseBookingClock(schedule.EndTime)
	if !ok {
		return false
	}

	windowStart := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), startClock.Hour(), startClock.Minute(), 0, 0, time.UTC)
	windowEnd := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), endClock.Hour(), endClock.Minute(), 0, 0, time.UTC)

	return !startTime.Before(windowStart) && !endTime.After(windowEnd)
}

func parseBookingClock(value string) (time.Time, bool) {
	for _, layout := range []string{"15:04", "15:04:05"} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, true
		}
	}

	return time.Time{}, false
}

func intervalsOverlap(startTime, endTime, bookedStart, bookedEnd time.Time) bool {
	return startTime.Before(bookedEnd) && endTime.After(bookedStart)
}

func (uc *BookingUseCase) GetUserBookings(ctx context.Context, userID string) ([]*domain.Booking, error) {
	return uc.bookingRepo.GetByUserID(ctx, userID)
}

func (uc *BookingUseCase) GetAllBookings(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error) {
	offset := (page - 1) * pageSize
	return uc.bookingRepo.GetAll(ctx, pageSize, offset)
}
