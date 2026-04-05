package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/butorovv/meeting-room-booking/internal/delivery/transport"
	"github.com/butorovv/meeting-room-booking/pkg/http_response"
	"github.com/butorovv/meeting-room-booking/pkg/logger"
)

type AuthUseCaseInterface interface {
	DummyLogin(ctx context.Context, role string) (string, error)
}

type AuthHandler struct {
	uc AuthUseCaseInterface
	rs *http_response.ResponseSender
}

func NewAuthHandler(uc AuthUseCaseInterface) *AuthHandler {
	return &AuthHandler{
		uc: uc,
		rs: http_response.NewResponseSender(logger.Global()),
	}
}

// DummyLogin godoc
// @Summary Получить тестовый JWT токен
// @Description Возвращает JWT токен для указанной роли (admin/user)
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body transport.DummyLoginRequest true "Роль (admin/user)"
// @Success 200 {object} transport.TokenResponse
// @Failure 400 {object} transport.ErrorResponse
// @Failure 415 {object} transport.ErrorResponse
// @Failure 405 {object} transport.ErrorResponse
// @Router /dummyLogin [post]
func (h *AuthHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.rs.Error(ctx, w, http.StatusMethodNotAllowed, "DummyLogin", transport.ErrInvalidRequest, nil)
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		h.rs.Error(ctx, w, http.StatusUnsupportedMediaType, "DummyLogin", transport.ErrInvalidRequest, nil)
		return
	}

	var req transport.DummyLoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.rs.Error(ctx, w, http.StatusBadRequest, "DummyLogin", transport.ErrInvalidRequest, err)
		return
	}

	if req.Role != "admin" && req.Role != "user" {
		h.rs.Error(ctx, w, http.StatusBadRequest, "DummyLogin", transport.ErrInvalidRequest, nil)
		return
	}

	token, err := h.uc.DummyLogin(ctx, req.Role)
	if err != nil {
		h.rs.Error(ctx, w, http.StatusInternalServerError, "DummyLogin", transport.ErrInternalError, err)
		return
	}

	h.rs.Send(ctx, w, http.StatusOK, transport.TokenResponse{Token: token})
}
