package repository

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"

	"github.com/butorovv/meeting-room-booking/internal/domain"
	"github.com/butorovv/meeting-room-booking/pkg/logger"
)

//go:embed sql/user/get_fixed_dummy.sql
var getFixedDummySQL string

//go:embed sql/user/get_by_id.sql
var getUserByIDSQL string

type UserRepository struct {
	db PgxIface // ← интерфейс, не конкретный тип
}

func NewUserRepository(db PgxIface) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetFixedDummy(ctx context.Context, role domain.UserRole) (*domain.User, error) {
	log := logger.FromContext(ctx)

	var id string
	if role == domain.RoleAdmin {
		id = "00000000-0000-0000-0000-000000000001"
	} else {
		id = "00000000-0000-0000-0000-000000000002"
	}
	email := string(role) + "@dummy.com"

	var u domain.User
	err := r.db.QueryRow(ctx, getFixedDummySQL, id, email, role).Scan(
		&u.ID, &u.Email, &u.Role, &u.CreatedAt,
	)
	if err != nil {
		log.ErrorContext(ctx, "repo GetFixedDummy failed", "err", err)
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	log := logger.FromContext(ctx)

	var u domain.User
	err := r.db.QueryRow(ctx, getUserByIDSQL, id).Scan(
		&u.ID, &u.Email, &u.Role, &u.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		log.ErrorContext(ctx, "repo GetByID failed", "id", id, "err", err)
		return nil, err
	}
	return &u, nil
}
