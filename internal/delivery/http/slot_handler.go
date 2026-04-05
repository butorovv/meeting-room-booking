package http

import (
	"context"
	"net/http"

	"github.com/butorovv/meeting-room-booking/internal/delivery/transport"
	"github.com/butorovv/meeting-room-booking/internal/domain"
	"github.com/butorovv/meeting-room-booking/pkg/http_response"
	"github.com/butorovv/meeting-room-booking/pkg/logger"
)

type SlotUseCaseInterface interface {
	GetAvailableSlots(ctx context.Context, roomID, date string) ([]*domain.Slot, error)
}

type SlotHandler struct {
	uc SlotUseCaseInterface
	rs *http_response.ResponseSender
}

func NewSlotHandler(uc SlotUseCaseInterface) *SlotHandler {
	return &SlotHandler{
		uc: uc,
		rs: http_response.NewResponseSender(logger.Global()),
	}
}

func (h *SlotHandler) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomID := r.PathValue("roomId")
	date := r.URL.Query().Get("date")

	if date == "" {
		h.rs.Error(ctx, w, http.StatusBadRequest, "GetAvailableSlots", transport.ErrInvalidRequest, nil)
		return
	}

	slots, err := h.uc.GetAvailableSlots(ctx, roomID, date)
	if err != nil {
		h.rs.Error(ctx, w, http.StatusInternalServerError, "GetAvailableSlots", transport.ErrInternalError, err)
		return
	}

	h.rs.Send(ctx, w, http.StatusOK, map[string]interface{}{"slots": transport.ToSlotResponseList(slots)})
}
