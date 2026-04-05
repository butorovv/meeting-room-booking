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

type ScheduleUseCaseInterface interface {
	CreateSchedule(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error)
}

type ScheduleHandler struct {
	uc ScheduleUseCaseInterface
	rs *http_response.ResponseSender
}

func NewScheduleHandler(uc ScheduleUseCaseInterface) *ScheduleHandler {
	return &ScheduleHandler{
		uc: uc,
		rs: http_response.NewResponseSender(logger.Global()),
	}
}

// CreateSchedule godoc
// @Summary Создать расписание
// @Description Создаёт расписание для переговорки (только admin, один раз)
// @Tags Schedules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param roomId path string true "ID переговорки"
// @Param request body transport.CreateScheduleRequest true "Данные расписания"
// @Success 201 {object} map[string]interface{} "schedule"
// @Failure 400 {object} transport.ErrorResponse
// @Failure 401 {object} transport.ErrorResponse
// @Failure 403 {object} transport.ErrorResponse
// @Failure 404 {object} transport.ErrorResponse
// @Failure 409 {object} transport.ErrorResponse
// @Router /rooms/{roomId}/schedule/create [post]
func (h *ScheduleHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.rs.Error(ctx, w, http.StatusMethodNotAllowed, "CreateSchedule", transport.ErrInvalidRequest, nil)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		h.rs.Error(ctx, w, http.StatusUnsupportedMediaType, "CreateSchedule", transport.ErrInvalidRequest, nil)
		return
	}

	roomID := r.PathValue("roomId")

	var req transport.CreateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.rs.Error(ctx, w, http.StatusBadRequest, "CreateSchedule", transport.ErrInvalidRequest, err)
		return
	}
	req.RoomID = roomID

	if err := req.Validate(); err != nil {
		h.rs.Error(ctx, w, http.StatusBadRequest, "CreateSchedule", transport.ErrInvalidRequest, err)
		return
	}

	schedule, err := h.uc.CreateSchedule(ctx, req.RoomID, req.DaysOfWeek, req.StartTime, req.EndTime)
	if err != nil {
		if err == domain.ErrScheduleExists {
			h.rs.Error(ctx, w, http.StatusConflict, "CreateSchedule", transport.ErrScheduleExists, err)
			return
		}
		if err == domain.ErrRoomNotFound {
			h.rs.Error(ctx, w, http.StatusNotFound, "CreateSchedule", transport.ErrRoomNotFound, err)
			return
		}
		h.rs.Error(ctx, w, http.StatusInternalServerError, "CreateSchedule", transport.ErrInternalError, err)
		return
	}

	h.rs.Send(ctx, w, http.StatusCreated, map[string]interface{}{
		"schedule": transport.ToScheduleResponse(schedule),
	})
}
