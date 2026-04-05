package repository

import (
	"context"
	_ "embed"

	"github.com/butorovv/meeting-room-booking/internal/domain"

	"github.com/jackc/pgx/v5"
)

//go:embed sql/slot/batch_create.sql
var batchCreateSlotsSQL string

//go:embed sql/slot/get_available_by_room_date.sql
var getAvailableSlotsSQL string

type SlotRepository struct {
	db PgxIface
}

func NewSlotRepository(db PgxIface) *SlotRepository {
	return &SlotRepository{db: db}
}

func (r *SlotRepository) BatchCreate(ctx context.Context, slots []*domain.Slot) error {
	batch := &pgx.Batch{}
	for _, slot := range slots {
		batch.Queue(batchCreateSlotsSQL, slot.ID, slot.RoomID, slot.StartTime, slot.EndTime, slot.CreatedAt)
	}
	return r.db.SendBatch(ctx, batch).Close()
}

func (r *SlotRepository) GetAvailableByRoomAndDate(ctx context.Context, roomID string, date string) ([]*domain.Slot, error) {
	rows, err := r.db.Query(ctx, getAvailableSlotsSQL, roomID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*domain.Slot
	for rows.Next() {
		var s domain.Slot
		var isBooked bool
		err := rows.Scan(&s.ID, &s.RoomID, &s.StartTime, &s.EndTime, &s.CreatedAt, &isBooked)
		if err != nil {
			return nil, err
		}
		if !isBooked {
			slots = append(slots, &s)
		}
	}
	return slots, rows.Err()
}

func (r *SlotRepository) GetByID(ctx context.Context, id string) (*domain.Slot, error) {
	var s domain.Slot
	err := r.db.QueryRow(ctx, "SELECT id, room_id, start_time, end_time, created_at FROM slots WHERE id = $1", id).
		Scan(&s.ID, &s.RoomID, &s.StartTime, &s.EndTime, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
