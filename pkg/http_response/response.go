package http_response

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/butorovv/meeting-room-booking/internal/delivery/transport"
	"github.com/butorovv/meeting-room-booking/pkg/logger"
)

type ResponseSender struct {
	log logger.Logger
}

func NewResponseSender(log logger.Logger) *ResponseSender {
	return &ResponseSender{log: log}
}

func (rs *ResponseSender) Send(ctx context.Context, w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		rs.log.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}

func (rs *ResponseSender) Error(ctx context.Context, w http.ResponseWriter, status int, operation string, errDetail transport.ErrorDetail, err error) {
	rs.log.ErrorContext(ctx, operation, "error", err)
	rs.Send(ctx, w, status, transport.ErrorResponse{Error: errDetail})
}
