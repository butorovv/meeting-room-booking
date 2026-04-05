package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/butorovv/meeting-room-booking/internal/domain"
	mock_usecase "github.com/butorovv/meeting-room-booking/internal/usecase/mock"
)

func TestAuthDummyLogin_Admin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_usecase.NewMockAuthRepository(ctrl)
	uc := NewAuthUseCase(mockRepo, "test_secret")

	expectedUser := &domain.User{
		ID:    "00000000-0000-0000-0000-000000000001",
		Email: "admin@dummy.com",
		Role:  domain.RoleAdmin,
	}
	mockRepo.EXPECT().
		GetFixedDummy(gomock.Any(), domain.RoleAdmin).
		Return(expectedUser, nil)

	token, err := uc.DummyLogin(context.Background(), "admin")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthDummyLogin_User(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_usecase.NewMockAuthRepository(ctrl)
	uc := NewAuthUseCase(mockRepo, "test_secret")

	expectedUser := &domain.User{
		ID:    "00000000-0000-0000-0000-000000000002",
		Email: "user@dummy.com",
		Role:  domain.RoleUser,
	}
	mockRepo.EXPECT().
		GetFixedDummy(gomock.Any(), domain.RoleUser).
		Return(expectedUser, nil)

	token, err := uc.DummyLogin(context.Background(), "user")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
