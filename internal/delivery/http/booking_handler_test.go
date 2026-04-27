package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock_usecase "github.com/butorovv/meeting-room-booking/internal/delivery/http/mock"
	"github.com/butorovv/meeting-room-booking/internal/delivery/middlewares"
	"github.com/butorovv/meeting-room-booking/internal/delivery/transport"
	"github.com/butorovv/meeting-room-booking/internal/domain"
)

func TestCreateBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockBookingUseCaseInterface(ctrl)
	handler := NewBookingHandler(mockUC)

	startTime := time.Date(2030, 4, 8, 9, 0, 0, 0, time.UTC)
	endTime := startTime.Add(30 * time.Minute)
	reqBody, _ := json.Marshal(transport.CreateBookingRequest{
		RoomID:    "room1",
		StartTime: startTime.Format(time.RFC3339),
		EndTime:   endTime.Format(time.RFC3339),
	})
	req := httptest.NewRequest(http.MethodPost, "/bookings/create", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), middlewares.UserIDKey, "user1")
	ctx = context.WithValue(ctx, middlewares.RoleKey, "user")
	req = req.WithContext(ctx)

	expectedBooking := &domain.Booking{ID: "book1", RoomID: "room1", StartTime: startTime, EndTime: endTime, Status: domain.BookingActive}
	mockUC.EXPECT().
		CreateBooking(gomock.Any(), "user1", "room1", startTime, endTime, nil).
		Return(expectedBooking, nil)

	handler.CreateBooking(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateBooking_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockBookingUseCaseInterface(ctrl)
	handler := NewBookingHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/bookings/create", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), middlewares.UserIDKey, "user1")
	ctx = context.WithValue(ctx, middlewares.RoleKey, "user")
	req = req.WithContext(ctx)

	handler.CreateBooking(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateBooking_RoomNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockBookingUseCaseInterface(ctrl)
	handler := NewBookingHandler(mockUC)

	startTime := time.Date(2030, 4, 8, 9, 0, 0, 0, time.UTC)
	endTime := startTime.Add(30 * time.Minute)
	reqBody, _ := json.Marshal(transport.CreateBookingRequest{
		RoomID:    "invalid",
		StartTime: startTime.Format(time.RFC3339),
		EndTime:   endTime.Format(time.RFC3339),
	})
	req := httptest.NewRequest(http.MethodPost, "/bookings/create", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), middlewares.UserIDKey, "user1")
	ctx = context.WithValue(ctx, middlewares.RoleKey, "user")
	req = req.WithContext(ctx)

	mockUC.EXPECT().
		CreateBooking(gomock.Any(), "user1", "invalid", startTime, endTime, nil).
		Return(nil, domain.ErrRoomNotFound)

	handler.CreateBooking(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateBooking_SlotAlreadyBooked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockBookingUseCaseInterface(ctrl)
	handler := NewBookingHandler(mockUC)

	startTime := time.Date(2030, 4, 8, 9, 0, 0, 0, time.UTC)
	endTime := startTime.Add(30 * time.Minute)
	reqBody, _ := json.Marshal(transport.CreateBookingRequest{
		RoomID:    "room1",
		StartTime: startTime.Format(time.RFC3339),
		EndTime:   endTime.Format(time.RFC3339),
	})
	req := httptest.NewRequest(http.MethodPost, "/bookings/create", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), middlewares.UserIDKey, "user1")
	req = req.WithContext(ctx)

	mockUC.EXPECT().
		CreateBooking(gomock.Any(), "user1", "room1", startTime, endTime, nil).
		Return(nil, domain.ErrSlotAlreadyBooked)

	handler.CreateBooking(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestMyBookings_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockBookingUseCaseInterface(ctrl)
	handler := NewBookingHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/bookings/my", nil)
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), middlewares.UserIDKey, "user1")
	ctx = context.WithValue(ctx, middlewares.RoleKey, "user")
	req = req.WithContext(ctx)

	expectedBookings := []*domain.Booking{{ID: "book1", UserID: "user1"}}
	mockUC.EXPECT().
		GetUserBookings(gomock.Any(), "user1").
		Return(expectedBookings, nil)

	handler.MyBookings(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCancelBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockBookingUseCaseInterface(ctrl)
	handler := NewBookingHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/bookings/book1/cancel", nil)
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), middlewares.UserIDKey, "user1")
	req = req.WithContext(ctx)
	req.SetPathValue("bookingId", "book1")

	expectedBooking := &domain.Booking{ID: "book1", Status: domain.BookingCancelled}
	mockUC.EXPECT().
		CancelBooking(gomock.Any(), "book1", "user1").
		Return(expectedBooking, nil)

	handler.CancelBooking(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListBookings_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockBookingUseCaseInterface(ctrl)
	handler := NewBookingHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/bookings/list?page=1&pageSize=10", nil)
	w := httptest.NewRecorder()

	mockUC.EXPECT().
		GetAllBookings(gomock.Any(), 1, 10).
		Return([]*domain.Booking{}, 0, nil)

	handler.ListBookings(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListBookings_DefaultPagination(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockBookingUseCaseInterface(ctrl)
	handler := NewBookingHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/bookings/list", nil)
	w := httptest.NewRecorder()

	mockUC.EXPECT().
		GetAllBookings(gomock.Any(), 1, 20).
		Return([]*domain.Booking{}, 0, nil)

	handler.ListBookings(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
