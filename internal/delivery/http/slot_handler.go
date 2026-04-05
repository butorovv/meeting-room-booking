package http

import (
	"context"
	"net/http"
	"time"

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

// GetAvailableSlots godoc
// @Summary Доступные слоты
// @Description Возвращает список доступных слотов по переговорке и дате
// @Tags Slots
// @Produce json
// @Security BearerAuth
// @Param roomId path string true "ID переговорки"
// @Param date query string true "Дата (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{} "slots"
// @Failure 400 {object} transport.ErrorResponse
// @Failure 401 {object} transport.ErrorResponse
// @Failure 404 {object} transport.ErrorResponse
// @Router /rooms/{roomId}/slots/list [get]
func (h *SlotHandler) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomID := r.PathValue("roomId")
	date := r.URL.Query().Get("date")

	if date == "" {
		h.rs.Error(ctx, w, http.StatusBadRequest, "GetAvailableSlots", transport.ErrInvalidRequest, nil)
		return
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		h.rs.Error(ctx, w, http.StatusBadRequest, "GetAvailableSlots", transport.ErrInvalidRequest, nil)
		return
	}

	slots, err := h.uc.GetAvailableSlots(ctx, roomID, date)
	if err != nil {
		if err == domain.ErrRoomNotFound {
			h.rs.Error(ctx, w, http.StatusNotFound, "GetAvailableSlots", transport.ErrRoomNotFound, err)
			return
		}
		h.rs.Error(ctx, w, http.StatusInternalServerError, "GetAvailableSlots", transport.ErrInternalError, err)
		return
	}

	h.rs.Send(ctx, w, http.StatusOK, map[string]interface{}{"slots": transport.ToSlotResponseList(slots)})
}
