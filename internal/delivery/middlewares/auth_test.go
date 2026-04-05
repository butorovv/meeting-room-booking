package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/butorovv/meeting-room-booking/pkg/jwt"
)

func TestWithUserID(t *testing.T) {
	ctx := WithUserID(context.Background(), "user123")
	id, ok := UserIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, "user123", id)
}

func TestWithRole(t *testing.T) {
	ctx := WithRole(context.Background(), "admin")
	role, ok := RoleFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, "admin", role)
}

func TestUserIDFromContext_NotFound(t *testing.T) {
	_, ok := UserIDFromContext(context.Background())
	assert.False(t, ok)
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	handler := AuthMiddleware("secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	handler := AuthMiddleware("secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	token, _ := jwt.GenerateToken("user1", "admin", "test_secret")

	handler := AuthMiddleware("test_secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := UserIDFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, "user1", userID)
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminOnly_WithoutRole(t *testing.T) {
	handler := AdminOnly(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAdminOnly_WithUserRole(t *testing.T) {
	handler := AdminOnly(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := WithRole(req.Context(), "user")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAdminOnly_WithAdminRole(t *testing.T) {
	handler := AdminOnly(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := WithRole(req.Context(), "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
