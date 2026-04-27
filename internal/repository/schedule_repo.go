package repository

import (
	"context"
	_ "embed"
	"time"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

//go:embed sql/schedule/create.sql
var createScheduleSQL string

//go:embed sql/schedule/get_by_room.sql
var getScheduleByRoomSQL string

//go:embed sql/schedule/exists_by_room.sql
var existsScheduleByRoomSQL string

type ScheduleRepository struct {
	db PgxIface
}

func NewScheduleRepository(db PgxIface) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(ctx context.Context, schedule *domain.Schedule) error {
	_, err := r.db.Exec(ctx, createScheduleSQL,
		schedule.ID, schedule.RoomID, schedule.DaysMask, schedule.StartTime, schedule.EndTime, schedule.CreatedAt,
	)
	return err
}

func (r *ScheduleRepository) GetByRoomID(ctx context.Context, roomID string) (*domain.Schedule, error) {
	var s domain.Schedule
	var startTime, endTime time.Time
	err := r.db.QueryRow(ctx, getScheduleByRoomSQL, roomID).Scan(
		&s.ID, &s.RoomID, &s.DaysMask, &startTime, &endTime, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	s.StartTime = startTime.Format("15:04")
	s.EndTime = endTime.Format("15:04")
	return &s, nil
}

func (r *ScheduleRepository) ExistsByRoomID(ctx context.Context, roomID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, existsScheduleByRoomSQL, roomID).Scan(&exists)
	return exists, err
}
