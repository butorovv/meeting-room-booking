package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/butorovv/meeting-room-booking/internal/domain"
	mock_usecase "github.com/butorovv/meeting-room-booking/internal/usecase/mock"
)

func TestCreateSchedule_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	uc := NewScheduleUseCase(mockScheduleRepo)

	mockScheduleRepo.EXPECT().
		ExistsByRoomID(gomock.Any(), "room1").
		Return(false, nil)

	mockScheduleRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	schedule, err := uc.CreateSchedule(context.Background(), "room1", []int{1, 2, 3}, "09:00", "18:00")

	assert.NoError(t, err)
	assert.Equal(t, "room1", schedule.RoomID)
}

func TestCreateSchedule_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	uc := NewScheduleUseCase(mockScheduleRepo)

	mockScheduleRepo.EXPECT().
		ExistsByRoomID(gomock.Any(), "room1").
		Return(true, nil)

	_, err := uc.CreateSchedule(context.Background(), "room1", []int{1, 2, 3}, "09:00", "18:00")

	assert.Equal(t, domain.ErrScheduleExists, err)
}

func TestCreateSchedule_ExistsCheckError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	uc := NewScheduleUseCase(mockScheduleRepo)

	mockScheduleRepo.EXPECT().
		ExistsByRoomID(gomock.Any(), "room1").
		Return(false, assert.AnError)

	_, err := uc.CreateSchedule(context.Background(), "room1", []int{1, 2, 3}, "09:00", "18:00")

	assert.Error(t, err)
}

func TestCreateSchedule_CreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	uc := NewScheduleUseCase(mockScheduleRepo)

	mockScheduleRepo.EXPECT().
		ExistsByRoomID(gomock.Any(), "room1").
		Return(false, nil)

	mockScheduleRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(assert.AnError)

	_, err := uc.CreateSchedule(context.Background(), "room1", []int{1, 2, 3}, "09:00", "18:00")

	assert.Error(t, err)
}

func TestCreateSchedule_EmptyDaysOfWeek(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockScheduleRepo := mock_usecase.NewMockScheduleRepositoryInterface(ctrl)
	uc := NewScheduleUseCase(mockScheduleRepo)

	mockScheduleRepo.EXPECT().
		ExistsByRoomID(gomock.Any(), "room1").
		Return(false, nil)

	mockScheduleRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	schedule, err := uc.CreateSchedule(context.Background(), "room1", []int{}, "09:00", "18:00")

	assert.NoError(t, err)
	assert.Equal(t, "room1", schedule.RoomID)
}
