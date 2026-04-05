package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := middleware.UserIDFromContext(ctx)

	var req transport.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.rs.Error(ctx, w, http.StatusBadRequest, "CreateBooking", transport.ErrInvalidRequest, err)
		return
	}

	var conferenceLink *string
	if req.CreateConferenceLink {
		link := "https://meet.example.com/" + req.SlotID
		conferenceLink = &link
	}

	booking, err := h.uc.CreateBooking(ctx, userID, req.SlotID, conferenceLink)
	if err != nil {
		if err == domain.ErrSlotInPast {
			h.rs.Error(ctx, w, http.StatusBadRequest, "CreateBooking", transport.ErrInvalidRequest, err)
			return
		}
		if err == domain.ErrSlotAlreadyBooked {
			h.rs.Error(ctx, w, http.StatusConflict, "CreateBooking", transport.ErrSlotAlreadyBooked, err)
			return
		}
		h.rs.Error(ctx, w, http.StatusInternalServerError, "CreateBooking", transport.ErrInternalError, err)
		return
	}

	h.rs.Send(ctx, w, http.StatusCreated, transport.ToBookingResponse(booking))
}

func (h *BookingHandler) MyBookings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := middleware.UserIDFromContext(ctx)

	bookings, err := h.uc.GetUserBookings(ctx, userID)
	if err != nil {
		h.rs.Error(ctx, w, http.StatusInternalServerError, "MyBookings", transport.ErrInternalError, err)
		return
	}

	h.rs.Send(ctx, w, http.StatusOK, map[string]interface{}{"bookings": transport.ToBookingResponseList(bookings)})
}

func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bookingID := r.PathValue("bookingId")
	userID, _ := middleware.UserIDFromContext(ctx)

	booking, err := h.uc.CancelBooking(ctx, bookingID, userID)
	if err != nil {
		h.rs.Error(ctx, w, http.StatusInternalServerError, "CancelBooking", transport.ErrInternalError, err)
		return
	}

	h.rs.Send(ctx, w, http.StatusOK, transport.ToBookingResponse(booking))
}

func (h *BookingHandler) ListBookings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := 1
	pageSize := 20

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := r.URL.Query().Get("pageSize"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
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
