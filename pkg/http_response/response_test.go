package http_response

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/butorovv/meeting-room-booking/internal/delivery/transport"
	"github.com/butorovv/meeting-room-booking/pkg/logger"
)

func TestResponseSender_Send(t *testing.T) {
	rs := NewResponseSender(logger.Global())
	w := httptest.NewRecorder()
	rs.Send(context.Background(), w, 200, map[string]string{"status": "ok"})
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestResponseSender_Error(t *testing.T) {
	rs := NewResponseSender(logger.Global())
	w := httptest.NewRecorder()
	rs.Error(context.Background(), w, 400, "TestError", transport.ErrInvalidRequest, nil)
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
}

func TestResponseSender_ErrorWithNilError(t *testing.T) {
	rs := NewResponseSender(logger.Global())
	w := httptest.NewRecorder()
	rs.Error(context.Background(), w, 500, "TestError", transport.ErrInternalError, nil)
	assert.Equal(t, 500, w.Code)
}
