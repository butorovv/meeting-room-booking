package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/butorovv/meeting-room-booking/internal/delivery/transport"
	"github.com/butorovv/meeting-room-booking/internal/domain"
	"github.com/butorovv/meeting-room-booking/pkg/http_response"
	"github.com/butorovv/meeting-room-booking/pkg/logger"
)

type RoomUseCaseInterface interface {
	CreateRoom(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error)
	ListRooms(ctx context.Context) ([]*domain.Room, error)
}

type RoomHandler struct {
	uc RoomUseCaseInterface
	rs *http_response.ResponseSender
}

func NewRoomHandler(uc RoomUseCaseInterface) *RoomHandler {
	return &RoomHandler{
		uc: uc,
		rs: http_response.NewResponseSender(logger.Global()),
	}
}

// CreateRoom godoc
// @Summary Создать переговорку
// @Description Создаёт новую переговорку (только admin)
// @Tags Rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body transport.CreateRoomRequest true "Данные переговорки"
// @Success 201 {object} map[string]interface{} "room"
// @Failure 400 {object} transport.ErrorResponse
// @Failure 401 {object} transport.ErrorResponse
// @Failure 403 {object} transport.ErrorResponse
// @Router /rooms/create [post]
func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.rs.Error(ctx, w, http.StatusMethodNotAllowed, "CreateRoom", transport.ErrInvalidRequest, nil)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		h.rs.Error(ctx, w, http.StatusUnsupportedMediaType, "CreateRoom", transport.ErrInvalidRequest, nil)
		return
	}

	var req transport.CreateRoomRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.rs.Error(ctx, w, http.StatusBadRequest, "CreateRoom", transport.ErrInvalidRequest, err)
		return
	}

	if err := req.Validate(); err != nil {
		h.rs.Error(ctx, w, http.StatusBadRequest, "CreateRoom", transport.ErrInvalidRequest, err)
		return
	}

	room, err := h.uc.CreateRoom(ctx, req.Name, req.Description, req.Capacity)
	if err != nil {
		h.rs.Error(ctx, w, http.StatusInternalServerError, "CreateRoom", transport.ErrInternalError, err)
		return
	}

	h.rs.Send(ctx, w, http.StatusCreated, map[string]interface{}{
		"room": transport.ToRoomResponse(room),
	})
}

// ListRooms godoc
// @Summary Список переговорок
// @Description Возвращает список всех переговорок
// @Tags Rooms
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "rooms"
// @Failure 401 {object} transport.ErrorResponse
// @Router /rooms/list [get]
func (h *RoomHandler) ListRooms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rooms, err := h.uc.ListRooms(ctx)
	if err != nil {
		h.rs.Error(ctx, w, http.StatusInternalServerError, "ListRooms", transport.ErrInternalError, err)
		return
	}

	h.rs.Send(ctx, w, http.StatusOK, map[string]interface{}{
		"rooms": transport.ToRoomResponseList(rooms),
	})
}
