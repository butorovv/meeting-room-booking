package repository

import (
	"context"
	_ "embed"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

//go:embed sql/booking/create.sql
var createBookingSQL string

//go:embed sql/booking/get_by_slot.sql
var getBookingBySlotSQL string

//go:embed sql/booking/get_by_user.sql
var getBookingsByUserSQL string

//go:embed sql/booking/get_all.sql
var getAllBookingsSQL string

//go:embed sql/booking/count_all.sql
var countAllBookingsSQL string

//go:embed sql/booking/update_status.sql
var updateBookingStatusSQL string

//go:embed sql/booking/get_by_id.sql
var getBookingByIDSQL string

const getActiveBookedIntervalsSQL = `
SELECT id, start_time, end_time
FROM bookings
WHERE room_id = $1
  AND start_time >= $2
  AND start_time < $3
  AND status = 'active'
ORDER BY start_time
`

type BookingRepository struct {
	db PgxIface
}

func NewBookingRepository(db PgxIface) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	_, err := r.db.Exec(ctx, createBookingSQL,
		booking.ID, booking.SlotID, booking.RoomID, booking.StartTime, booking.EndTime,
		booking.UserID, booking.Status, booking.ConferenceLink, booking.CreatedAt, booking.UpdatedAt,
	)
	return err
}

func (r *BookingRepository) GetBySlotID(ctx context.Context, slotID string) (*domain.Booking, error) {
	var b domain.Booking
	var nullableSlotID pgtype.UUID
	err := r.db.QueryRow(ctx, getBookingBySlotSQL, slotID).Scan(
		&b.ID, &nullableSlotID, &b.RoomID, &b.StartTime, &b.EndTime,
		&b.UserID, &b.Status, &b.ConferenceLink, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	setBookingSlotID(&b, nullableSlotID)
	return &b, nil
}

func (r *BookingRepository) GetActiveBookedIntervals(ctx context.Context, roomID string, startDate, endDate time.Time) ([]domain.BookedInterval, error) {
	rows, err := r.db.Query(ctx, getActiveBookedIntervalsSQL, roomID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	intervals := make([]domain.BookedInterval, 0)
	for rows.Next() {
		var interval domain.BookedInterval
		if err := rows.Scan(&interval.ID, &interval.StartTime, &interval.EndTime); err != nil {
			return nil, err
		}
		intervals = append(intervals, interval)
	}

	return intervals, rows.Err()
}

func (r *BookingRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Booking, error) {
	rows, err := r.db.Query(ctx, getBookingsByUserSQL, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var b domain.Booking
		var nullableSlotID pgtype.UUID
		err := rows.Scan(
			&b.ID, &nullableSlotID, &b.RoomID, &b.StartTime, &b.EndTime,
			&b.UserID, &b.Status, &b.ConferenceLink, &b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		setBookingSlotID(&b, nullableSlotID)
		bookings = append(bookings, &b)
	}
	return bookings, rows.Err()
}

func (r *BookingRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Booking, int, error) {
	rows, err := r.db.Query(ctx, getAllBookingsSQL, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var b domain.Booking
		var nullableSlotID pgtype.UUID
		err := rows.Scan(
			&b.ID, &nullableSlotID, &b.RoomID, &b.StartTime, &b.EndTime,
			&b.UserID, &b.Status, &b.ConferenceLink, &b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		setBookingSlotID(&b, nullableSlotID)
		bookings = append(bookings, &b)
	}

	var total int
	err = r.db.QueryRow(ctx, countAllBookingsSQL).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	return bookings, total, nil
}

func (r *BookingRepository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(ctx, updateBookingStatusSQL, status, id)
	return err
}

func (r *BookingRepository) GetByID(ctx context.Context, id string) (*domain.Booking, error) {
	var b domain.Booking
	var nullableSlotID pgtype.UUID
	err := r.db.QueryRow(ctx, getBookingByIDSQL, id).Scan(
		&b.ID, &nullableSlotID, &b.RoomID, &b.StartTime, &b.EndTime,
		&b.UserID, &b.Status, &b.ConferenceLink, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	setBookingSlotID(&b, nullableSlotID)
	return &b, nil
}

func setBookingSlotID(booking *domain.Booking, slotID pgtype.UUID) {
	if !slotID.Valid {
		return
	}

	value := uuid.UUID(slotID.Bytes).String()
	booking.SlotID = &value
}

func (r *BookingRepository) CreateWithTx(ctx context.Context, tx pgx.Tx, booking *domain.Booking) error {
	query := `
        INSERT INTO bookings (id, room_id, start_time, end_time, user_id, status, conference_link, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err := tx.Exec(ctx, query,
		booking.ID,
		booking.RoomID,
		booking.StartTime,
		booking.EndTime,
		booking.UserID,
		booking.Status,
		booking.ConferenceLink,
		booking.CreatedAt,
		booking.UpdatedAt,
	)
	return err
}

func (r *BookingRepository) UpdateStatusWithTx(ctx context.Context, tx pgx.Tx, id, status string) error {
	query := `UPDATE bookings SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := tx.Exec(ctx, query, status, id)
	return err
}
