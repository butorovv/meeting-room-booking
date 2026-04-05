package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock_usecase "github.com/butorovv/meeting-room-booking/internal/delivery/http/mock"
	"github.com/butorovv/meeting-room-booking/internal/delivery/transport"
)

func TestAuthHandler_DummyLogin_Success_Admin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockAuthUseCaseInterface(ctrl)
	handler := NewAuthHandler(mockUC)

	reqBody, _ := json.Marshal(transport.DummyLoginRequest{Role: "admin"})
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUC.EXPECT().
		DummyLogin(gomock.Any(), "admin").
		Return("test_token_admin", nil)

	handler.DummyLogin(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp transport.TokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "test_token_admin", resp.Token)
}

func TestAuthHandler_DummyLogin_Success_User(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockAuthUseCaseInterface(ctrl)
	handler := NewAuthHandler(mockUC)

	reqBody, _ := json.Marshal(transport.DummyLoginRequest{Role: "user"})
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUC.EXPECT().
		DummyLogin(gomock.Any(), "user").
		Return("test_token_user", nil)

	handler.DummyLogin(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp transport.TokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "test_token_user", resp.Token)
}

func TestAuthHandler_DummyLogin_InvalidRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockAuthUseCaseInterface(ctrl)
	handler := NewAuthHandler(mockUC)

	reqBody, _ := json.Marshal(transport.DummyLoginRequest{Role: "superadmin"})
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.DummyLogin(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errResp transport.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Equal(t, "INVALID_REQUEST", errResp.Error.Code)
}

func TestAuthHandler_DummyLogin_EmptyBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockAuthUseCaseInterface(ctrl)
	handler := NewAuthHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.DummyLogin(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_DummyLogin_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockAuthUseCaseInterface(ctrl)
	handler := NewAuthHandler(mockUC)

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader([]byte(`{"role":`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.DummyLogin(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_DummyLogin_EmptyRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockAuthUseCaseInterface(ctrl)
	handler := NewAuthHandler(mockUC)

	reqBody, _ := json.Marshal(transport.DummyLoginRequest{Role: ""})
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.DummyLogin(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_DummyLogin_UseCaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockAuthUseCaseInterface(ctrl)
	handler := NewAuthHandler(mockUC)

	reqBody, _ := json.Marshal(transport.DummyLoginRequest{Role: "admin"})
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUC.EXPECT().
		DummyLogin(gomock.Any(), "admin").
		Return("", assert.AnError)

	handler.DummyLogin(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAuthHandler_DummyLogin_WrongMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockAuthUseCaseInterface(ctrl)
	handler := NewAuthHandler(mockUC)

	req := httptest.NewRequest(http.MethodGet, "/dummyLogin", nil)
	w := httptest.NewRecorder()

	handler.DummyLogin(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestAuthHandler_DummyLogin_WrongContentType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mock_usecase.NewMockAuthUseCaseInterface(ctrl)
	handler := NewAuthHandler(mockUC)

	reqBody, _ := json.Marshal(transport.DummyLoginRequest{Role: "admin"})
	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	handler.DummyLogin(w, req)

	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
}
