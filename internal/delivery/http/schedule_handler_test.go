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

func TestScheduleHandler_CreateSchedule_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockScheduleUseCaseInterface(ctrl)
	handler := NewScheduleHandler(mockUC)

	reqBody, _ := json.Marshal(transport.CreateScheduleRequest{
		DaysOfWeek: []int{1, 2, 3},
		StartTime:  "09:00",
		EndTime:    "18:00",
	})
	req := httptest.NewRequest(http.MethodPost, "/rooms/room1/schedule/create", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	expectedSchedule := &domain.Schedule{ID: "sch1", RoomID: "room1"}
	mockUC.EXPECT().
		CreateSchedule(gomock.Any(), "room1", []int{1, 2, 3}, "09:00", "18:00").
		Return(expectedSchedule, nil)

	req.SetPathValue("roomId", "room1")

	handler.CreateSchedule(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp struct {
		Schedule transport.ScheduleResponse `json:"schedule"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "sch1", resp.Schedule.ID)
	assert.Equal(t, "room1", resp.Schedule.RoomID)
}

func TestScheduleHandler_CreateSchedule_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockScheduleUseCaseInterface(ctrl)
	handler := NewScheduleHandler(mockUC)

	reqBody, _ := json.Marshal(transport.CreateScheduleRequest{
		DaysOfWeek: []int{1, 2, 3},
		StartTime:  "09:00",
		EndTime:    "18:00",
	})
	req := httptest.NewRequest(http.MethodPost, "/rooms/room1/schedule/create", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUC.EXPECT().
		CreateSchedule(gomock.Any(), "room1", []int{1, 2, 3}, "09:00", "18:00").
		Return(nil, domain.ErrScheduleExists)

	req.SetPathValue("roomId", "room1")

	handler.CreateSchedule(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var errResp transport.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Equal(t, "SCHEDULE_EXISTS", errResp.Error.Code)
}

func TestScheduleHandler_CreateSchedule_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockScheduleUseCaseInterface(ctrl)
	handler := NewScheduleHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/rooms/room1/schedule/create", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	req.SetPathValue("roomId", "room1")

	handler.CreateSchedule(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestScheduleHandler_CreateSchedule_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockScheduleUseCaseInterface(ctrl)
	handler := NewScheduleHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/rooms/room1/schedule/create", bytes.NewReader([]byte(`{"daysOfWeek":`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	req.SetPathValue("roomId", "room1")

	handler.CreateSchedule(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestScheduleHandler_CreateSchedule_InvalidDaysOfWeek(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockScheduleUseCaseInterface(ctrl)
	handler := NewScheduleHandler(mockUC)

	reqBody, _ := json.Marshal(transport.CreateScheduleRequest{
		DaysOfWeek: []int{0, 8}, // невалидные дни
		StartTime:  "09:00",
		EndTime:    "18:00",
	})
	req := httptest.NewRequest(http.MethodPost, "/rooms/room1/schedule/create", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	req.SetPathValue("roomId", "room1")

	handler.CreateSchedule(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestScheduleHandler_CreateSchedule_InvalidTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockScheduleUseCaseInterface(ctrl)
	handler := NewScheduleHandler(mockUC)

	reqBody, _ := json.Marshal(transport.CreateScheduleRequest{
		DaysOfWeek: []int{1, 2, 3},
		StartTime:  "25:00", // невалидное время
		EndTime:    "18:00",
	})
	req := httptest.NewRequest(http.MethodPost, "/rooms/room1/schedule/create", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	req.SetPathValue("roomId", "room1")

	handler.CreateSchedule(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestScheduleHandler_CreateSchedule_EmptyRoomID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockScheduleUseCaseInterface(ctrl)
	handler := NewScheduleHandler(mockUC)

	reqBody, _ := json.Marshal(transport.CreateScheduleRequest{
		DaysOfWeek: []int{1, 2, 3},
		StartTime:  "09:00",
		EndTime:    "18:00",
	})
	req := httptest.NewRequest(http.MethodPost, "/rooms//schedule/create", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	req.SetPathValue("roomId", "")

	handler.CreateSchedule(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestScheduleHandler_CreateSchedule_RoomNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockScheduleUseCaseInterface(ctrl)
	handler := NewScheduleHandler(mockUC)

	reqBody, _ := json.Marshal(transport.CreateScheduleRequest{
		DaysOfWeek: []int{1, 2, 3},
		StartTime:  "09:00",
		EndTime:    "18:00",
	})
	req := httptest.NewRequest(http.MethodPost, "/rooms/room404/schedule/create", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	req.SetPathValue("roomId", "room404")

	mockUC.EXPECT().
		CreateSchedule(gomock.Any(), "room404", []int{1, 2, 3}, "09:00", "18:00").
		Return(nil, domain.ErrRoomNotFound)

	handler.CreateSchedule(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestScheduleHandler_CreateSchedule_WrongMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockScheduleUseCaseInterface(ctrl)
	handler := NewScheduleHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/rooms/room1/schedule/create", nil)
	w := httptest.NewRecorder()
	req.SetPathValue("roomId", "room1")

	handler.CreateSchedule(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestScheduleHandler_CreateSchedule_WrongContentType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockScheduleUseCaseInterface(ctrl)
	handler := NewScheduleHandler(mockUC)

	reqBody, _ := json.Marshal(transport.CreateScheduleRequest{
		DaysOfWeek: []int{1, 2, 3},
		StartTime:  "09:00",
		EndTime:    "18:00",
	})
	req := httptest.NewRequest(http.MethodPost, "/rooms/room1/schedule/create", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()
	req.SetPathValue("roomId", "room1")

	handler.CreateSchedule(w, req)

	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
}

func TestScheduleHandler_CreateSchedule_UseCaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockScheduleUseCaseInterface(ctrl)
	handler := NewScheduleHandler(mockUC)

	reqBody, _ := json.Marshal(transport.CreateScheduleRequest{
		DaysOfWeek: []int{1, 2, 3},
		StartTime:  "09:00",
		EndTime:    "18:00",
	})
	req := httptest.NewRequest(http.MethodPost, "/rooms/room1/schedule/create", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	req.SetPathValue("roomId", "room1")

	mockUC.EXPECT().
		CreateSchedule(gomock.Any(), "room1", []int{1, 2, 3}, "09:00", "18:00").
		Return(nil, assert.AnError)

	handler.CreateSchedule(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
