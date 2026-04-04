package usecase

import (
	"context"

	"github.com/butorovv/meeting-room-booking/internal/domain"
	"github.com/butorovv/meeting-room-booking/pkg/jwt"
)

type AuthRepository interface {
	GetFixedDummy(ctx context.Context, role domain.UserRole) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

type AuthUseCase struct {
	repo      AuthRepository
	jwtSecret string
}

func NewAuthUseCase(repo AuthRepository, jwtSecret string) *AuthUseCase {
	return &AuthUseCase{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (uc *AuthUseCase) DummyLogin(ctx context.Context, role string) (string, error) {
	var userRole domain.UserRole
	if role == "admin" {
		userRole = domain.RoleAdmin
	} else {
		userRole = domain.RoleUser
	}

	user, err := uc.repo.GetFixedDummy(ctx, userRole)
	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateToken(user.ID, string(user.Role), uc.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
