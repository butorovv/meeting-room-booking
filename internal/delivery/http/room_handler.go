package http

import (
	"context"
	"encoding/json"
	"net/http"

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

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req transport.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

	h.rs.Send(ctx, w, http.StatusCreated, transport.ToRoomResponse(room))
}

func (h *RoomHandler) ListRooms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rooms, err := h.uc.ListRooms(ctx)
	if err != nil {
		h.rs.Error(ctx, w, http.StatusInternalServerError, "ListRooms", transport.ErrInternalError, err)
		return
	}

	h.rs.Send(ctx, w, http.StatusOK, map[string]interface{}{"rooms": transport.ToRoomResponseList(rooms)})
}
