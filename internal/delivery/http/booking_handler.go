package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	middleware "github.com/butorovv/meeting-room-booking/internal/delivery/middlewares"
	"github.com/butorovv/meeting-room-booking/internal/delivery/transport"
	"github.com/butorovv/meeting-room-booking/internal/domain"
	"github.com/butorovv/meeting-room-booking/pkg/http_response"
	"github.com/butorovv/meeting-room-booking/pkg/logger"
)

type BookingUseCaseInterface interface {
	CreateBooking(ctx context.Context, userID, slotID string, conferenceLink *string) (*domain.Booking, error)
	CancelBooking(ctx context.Context, bookingID, userID string) (*domain.Booking, error)
	GetUserBookings(ctx context.Context, userID string) ([]*domain.Booking, error)
	GetAllBookings(ctx context.Context, page, pageSize int) ([]*domain.Booking, int, error)
}

type BookingHandler struct {
	uc BookingUseCaseInterface
	rs *http_response.ResponseSender
}

func NewBookingHandler(uc BookingUseCaseInterface) *BookingHandler {
	return &BookingHandler{
		uc: uc,
		rs: http_response.NewResponseSender(logger.Global()),
	}
}

// CreateBooking godoc
// @Summary Создать бронь
// @Description Создаёт бронь на слот (только user)
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body transport.CreateBookingRequest true "Данные брони"
// @Success 201 {object} map[string]interface{} "booking"
// @Failure 400 {object} transport.ErrorResponse
// @Failure 401 {object} transport.ErrorResponse
// @Failure 403 {object} transport.ErrorResponse
// @Failure 404 {object} transport.ErrorResponse
// @Failure 409 {object} transport.ErrorResponse
// @Router /bookings/create [post]
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		h.rs.Error(ctx, w, http.StatusUnauthorized, "CreateBooking", transport.ErrUnauthorized, nil)
		return
	}

	if r.Method != http.MethodPost {
		h.rs.Error(ctx, w, http.StatusMethodNotAllowed, "CreateBooking", transport.ErrInvalidRequest, nil)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		h.rs.Error(ctx, w, http.StatusUnsupportedMediaType, "CreateBooking", transport.ErrInvalidRequest, nil)
		return
	}

	var req transport.CreateBookingRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.rs.Error(ctx, w, http.StatusBadRequest, "CreateBooking", transport.ErrInvalidRequest, err)
		return
	}

	if strings.TrimSpace(req.SlotID) == "" {
		h.rs.Error(ctx, w, http.StatusBadRequest, "CreateBooking", transport.ErrInvalidRequest, nil)
		return
	}

	var conferenceLink *string
	if req.CreateConferenceLink {
		link := "https://meet.example.com/" + req.SlotID
		conferenceLink = &link
	}

	booking, err := h.uc.CreateBooking(ctx, userID, req.SlotID, conferenceLink)
	if err != nil {
		switch err {
		case domain.ErrSlotNotFound:
			h.rs.Error(ctx, w, http.StatusNotFound, "CreateBooking", transport.ErrSlotNotFound, err)
			return
		case domain.ErrSlotInPast:
			h.rs.Error(ctx, w, http.StatusBadRequest, "CreateBooking", transport.ErrInvalidRequest, err)
			return
		case domain.ErrSlotAlreadyBooked:
			h.rs.Error(ctx, w, http.StatusConflict, "CreateBooking", transport.ErrSlotAlreadyBooked, err)
			return
		default:
			h.rs.Error(ctx, w, http.StatusInternalServerError, "CreateBooking", transport.ErrInternalError, err)
			return
		}
	}

	h.rs.Send(ctx, w, http.StatusCreated, map[string]interface{}{
		"booking": transport.ToBookingResponse(booking),
	})
}

// MyBookings godoc
// @Summary Мои брони
// @Description Возвращает список броней текущего пользователя (только user)
// @Tags Bookings
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "bookings"
// @Failure 401 {object} transport.ErrorResponse
// @Failure 403 {object} transport.ErrorResponse
// @Router /bookings/my [get]
func (h *BookingHandler) MyBookings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	role, ok := middleware.RoleFromContext(ctx)
	if !ok || role != "user" {
		h.rs.Error(ctx, w, http.StatusForbidden, "MyBookings", transport.ErrForbidden, nil)
		return
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		h.rs.Error(ctx, w, http.StatusUnauthorized, "MyBookings", transport.ErrUnauthorized, nil)
		return
	}

	if r.Method != http.MethodGet {
		h.rs.Error(ctx, w, http.StatusMethodNotAllowed, "MyBookings", transport.ErrInvalidRequest, nil)
		return
	}

	bookings, err := h.uc.GetUserBookings(ctx, userID)
	if err != nil {
		h.rs.Error(ctx, w, http.StatusInternalServerError, "MyBookings", transport.ErrInternalError, err)
		return
	}

	h.rs.Send(ctx, w, http.StatusOK, map[string]interface{}{
		"bookings": transport.ToBookingResponseList(bookings),
	})
}

// CancelBooking godoc
// @Summary Отменить бронь
// @Description Отменяет бронь по ID (только user, только свою)
// @Tags Bookings
// @Produce json
// @Security BearerAuth
// @Param bookingId path string true "ID брони"
// @Success 200 {object} map[string]interface{} "booking"
// @Failure 400 {object} transport.ErrorResponse
// @Failure 401 {object} transport.ErrorResponse
// @Failure 403 {object} transport.ErrorResponse
// @Failure 404 {object} transport.ErrorResponse
// @Router /bookings/{bookingId}/cancel [post]
func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bookingID := r.PathValue("bookingId")
	if strings.TrimSpace(bookingID) == "" {
		h.rs.Error(ctx, w, http.StatusBadRequest, "CancelBooking", transport.ErrInvalidRequest, nil)
		return
	}

	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		h.rs.Error(ctx, w, http.StatusUnauthorized, "CancelBooking", transport.ErrUnauthorized, nil)
		return
	}

	if r.Method != http.MethodPost {
		h.rs.Error(ctx, w, http.StatusMethodNotAllowed, "CancelBooking", transport.ErrInvalidRequest, nil)
		return
	}

	booking, err := h.uc.CancelBooking(ctx, bookingID, userID)
	if err != nil {
		switch err {
		case domain.ErrBookingNotFound:
			h.rs.Error(ctx, w, http.StatusNotFound, "CancelBooking", transport.ErrBookingNotFound, err)
			return
		case domain.ErrForbidden:
			h.rs.Error(ctx, w, http.StatusForbidden, "CancelBooking", transport.ErrForbidden, err)
			return
		default:
			h.rs.Error(ctx, w, http.StatusInternalServerError, "CancelBooking", transport.ErrInternalError, err)
			return
		}
	}

	h.rs.Send(ctx, w, http.StatusOK, map[string]interface{}{
		"booking": transport.ToBookingResponse(booking),
	})
}

// ListBookings godoc
// @Summary Список всех броней
// @Description Возвращает список всех броней с пагинацией (только admin)
// @Tags Bookings
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы" default(1)
// @Param pageSize query int false "Размер страницы" default(20) maximum(100)
// @Success 200 {object} transport.BookingsListResponse
// @Failure 400 {object} transport.ErrorResponse
// @Failure 401 {object} transport.ErrorResponse
// @Failure 403 {object} transport.ErrorResponse
// @Router /bookings/list [get]
func (h *BookingHandler) ListBookings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodGet {
		h.rs.Error(ctx, w, http.StatusMethodNotAllowed, "ListBookings", transport.ErrInvalidRequest, nil)
		return
	}

	page := 1
	pageSize := 20

	if p := r.URL.Query().Get("page"); p != "" {
		parsed, err := strconv.Atoi(p)
		if err != nil || parsed < 1 {
			h.rs.Error(ctx, w, http.StatusBadRequest, "ListBookings", transport.ErrInvalidRequest, nil)
			return
		}
		page = parsed
	}

	if ps := r.URL.Query().Get("pageSize"); ps != "" {
		parsed, err := strconv.Atoi(ps)
		if err != nil || parsed < 1 || parsed > 100 {
			h.rs.Error(ctx, w, http.StatusBadRequest, "ListBookings", transport.ErrInvalidRequest, nil)
			return
		}
		pageSize = parsed
	}

	bookings, total, err := h.uc.GetAllBookings(ctx, page, pageSize)
	if err != nil {
		h.rs.Error(ctx, w, http.StatusInternalServerError, "ListBookings", transport.ErrInternalError, err)
		return
	}

	h.rs.Send(ctx, w, http.StatusOK, transport.BookingsListResponse{
		Bookings: transport.ToBookingResponseList(bookings),
		Pagination: transport.PaginationResponse{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	})
}
