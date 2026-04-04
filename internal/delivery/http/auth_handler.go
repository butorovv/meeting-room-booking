package http

import (
	"context"
	"encoding/json"
	"net/http"

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

func (h *AuthHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req transport.DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
