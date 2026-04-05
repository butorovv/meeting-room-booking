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
