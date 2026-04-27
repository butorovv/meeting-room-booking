package usecase

import (
	"context"

	"github.com/butorovv/meeting-room-booking/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PgxPoolIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Close()
}

type SlotRepositoryInterface interface {
	GetByID(ctx context.Context, id string) (*domain.Slot, error)
	BatchCreate(ctx context.Context, slots []*domain.Slot) error
	GetAvailableByRoomAndDate(ctx context.Context, roomID, date string) ([]*domain.Slot, error)
}
