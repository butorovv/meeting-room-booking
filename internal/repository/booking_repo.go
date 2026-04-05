package repository

import (
	"context"
	_ "embed"

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

type BookingRepository struct {
	db PgxIface
}

func NewBookingRepository(db PgxIface) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	_, err := r.db.Exec(ctx, createBookingSQL,
		booking.ID, booking.SlotID, booking.UserID, booking.Status, booking.ConferenceLink, booking.CreatedAt, booking.UpdatedAt,
	)
	return err
}

func (r *BookingRepository) GetBySlotID(ctx context.Context, slotID string) (*domain.Booking, error) {
	var b domain.Booking
	err := r.db.QueryRow(ctx, getBookingBySlotSQL, slotID).Scan(
		&b.ID, &b.SlotID, &b.UserID, &b.Status, &b.ConferenceLink, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
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
		err := rows.Scan(&b.ID, &b.SlotID, &b.UserID, &b.Status, &b.ConferenceLink, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
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
		err := rows.Scan(&b.ID, &b.SlotID, &b.UserID, &b.Status, &b.ConferenceLink, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
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
	err := r.db.QueryRow(ctx, getBookingByIDSQL, id).Scan(
		&b.ID, &b.SlotID, &b.UserID, &b.Status, &b.ConferenceLink, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
