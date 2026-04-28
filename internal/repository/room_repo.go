package repository

import (
	"context"
	_ "embed"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

//go:embed sql/room/create.sql
var createRoomSQL string

//go:embed sql/room/list.sql
var listRoomsSQL string

const existsRoomByIDSQL = `SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)`

type RoomRepository struct {
	db PgxIface
}

func NewRoomRepository(db PgxIface) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(ctx context.Context, room *domain.Room) error {
	_, err := r.db.Exec(ctx, createRoomSQL,
		room.ID, room.Name, room.Description, room.Capacity, room.CreatedAt,
	)
	return err
}

func (r *RoomRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, existsRoomByIDSQL, id).Scan(&exists)
	return exists, err
}

func (r *RoomRepository) List(ctx context.Context) ([]*domain.Room, error) {
	rows, err := r.db.Query(ctx, listRoomsSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*domain.Room
	for rows.Next() {
		var room domain.Room
		err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.Capacity, &room.CreatedAt)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}
	return rooms, rows.Err()
}
