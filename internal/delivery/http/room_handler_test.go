package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock_usecase "github.com/butorovv/meeting-room-booking/internal/delivery/http/mock"
	"github.com/butorovv/meeting-room-booking/internal/delivery/transport"
	"github.com/butorovv/meeting-room-booking/internal/domain"
)

func TestRoomHandler_CreateRoom_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockRoomUseCaseInterface(ctrl)
	handler := NewRoomHandler(mockUC)

	reqBody, _ := json.Marshal(transport.CreateRoomRequest{Name: "Test Room"})
	req := httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	expectedRoom := &domain.Room{ID: "room1", Name: "Test Room"}
	mockUC.EXPECT().
		CreateRoom(gomock.Any(), "Test Room", nil, nil).
		Return(expectedRoom, nil)

	handler.CreateRoom(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRoomHandler_ListRooms_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockRoomUseCaseInterface(ctrl)
	handler := NewRoomHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/rooms/list", nil)
	w := httptest.NewRecorder()

	expectedRooms := []*domain.Room{{ID: "room1", Name: "Room 1"}}
	mockUC.EXPECT().
		ListRooms(gomock.Any()).
		Return(expectedRooms, nil)

	handler.ListRooms(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateRoom_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockRoomUseCaseInterface(ctrl)
	handler := NewRoomHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/rooms/create", nil)
	w := httptest.NewRecorder()

	handler.CreateRoom(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListRooms_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockRoomUseCaseInterface(ctrl)
	handler := NewRoomHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/rooms/list", nil)
	w := httptest.NewRecorder()

	mockUC.EXPECT().
		ListRooms(gomock.Any()).
		Return([]*domain.Room{}, nil)

	handler.ListRooms(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
