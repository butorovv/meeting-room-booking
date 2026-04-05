package usecase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/butorovv/meeting-room-booking/internal/domain"
	mock_repo "github.com/butorovv/meeting-room-booking/internal/usecase/mock"
)

func TestCreateRoom_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockRoomRepositoryInterface(ctrl)
	uc := NewRoomUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	room, err := uc.CreateRoom(context.Background(), "Test Room", nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, "Test Room", room.Name)
}

func TestCreateRoom_WithDescriptionAndCapacity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockRoomRepositoryInterface(ctrl)
	uc := NewRoomUseCase(mockRepo)

	description := "Test Description"
	capacity := 10

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	room, err := uc.CreateRoom(context.Background(), "Test Room", &description, &capacity)

	assert.NoError(t, err)
	assert.Equal(t, "Test Room", room.Name)
	assert.Equal(t, &description, room.Description)
	assert.Equal(t, &capacity, room.Capacity)
}

func TestCreateRoom_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockRoomRepositoryInterface(ctrl)
	uc := NewRoomUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	room, err := uc.CreateRoom(context.Background(), "", nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, "", room.Name)
}

func TestCreateRoom_CreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockRoomRepositoryInterface(ctrl)
	uc := NewRoomUseCase(mockRepo)

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(assert.AnError)

	_, err := uc.CreateRoom(context.Background(), "Test Room", nil, nil)

	assert.Error(t, err)
}

func TestListRooms_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockRoomRepositoryInterface(ctrl)
	uc := NewRoomUseCase(mockRepo)

	expectedRooms := []*domain.Room{
		{ID: "1", Name: "Room 1"},
		{ID: "2", Name: "Room 2"},
	}
	mockRepo.EXPECT().
		List(gomock.Any()).
		Return(expectedRooms, nil)

	rooms, err := uc.ListRooms(context.Background())

	assert.NoError(t, err)
	assert.Len(t, rooms, 2)
	assert.Equal(t, "Room 1", rooms[0].Name)
}

func TestListRooms_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockRoomRepositoryInterface(ctrl)
	uc := NewRoomUseCase(mockRepo)

	mockRepo.EXPECT().
		List(gomock.Any()).
		Return([]*domain.Room{}, nil)

	rooms, err := uc.ListRooms(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, rooms)
}

func TestListRooms_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockRoomRepositoryInterface(ctrl)
	uc := NewRoomUseCase(mockRepo)

	mockRepo.EXPECT().
		List(gomock.Any()).
		Return(nil, assert.AnError)

	_, err := uc.ListRooms(context.Background())

	assert.Error(t, err)
}
