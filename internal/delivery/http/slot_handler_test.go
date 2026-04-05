package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock_usecase "github.com/butorovv/meeting-room-booking/internal/delivery/http/mock"
	"github.com/butorovv/meeting-room-booking/internal/domain"
)

func TestSlotHandler_GetAvailableSlots_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockSlotUseCaseInterface(ctrl)
	handler := NewSlotHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/rooms/room1/slots/list?date=2024-04-07", nil)
	w := httptest.NewRecorder()

	expectedSlots := []*domain.Slot{{ID: "slot1"}, {ID: "slot2"}}
	mockUC.EXPECT().
		GetAvailableSlots(gomock.Any(), "room1", "2024-04-07").
		Return(expectedSlots, nil)

	req.SetPathValue("roomId", "room1")

	handler.GetAvailableSlots(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSlotHandler_GetAvailableSlots_MissingDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockSlotUseCaseInterface(ctrl)
	handler := NewSlotHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/rooms/room1/slots/list", nil)
	w := httptest.NewRecorder()

	req.SetPathValue("roomId", "room1")

	handler.GetAvailableSlots(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAvailableSlots_InvalidRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockSlotUseCaseInterface(ctrl)
	handler := NewSlotHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/rooms/invalid/slots/list?date=2024-04-07", nil)
	w := httptest.NewRecorder()
	req.SetPathValue("roomId", "invalid")

	mockUC.EXPECT().
		GetAvailableSlots(gomock.Any(), "invalid", "2024-04-07").
		Return(nil, domain.ErrRoomNotFound)

	handler.GetAvailableSlots(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
