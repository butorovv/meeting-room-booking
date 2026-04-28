package http

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Info(w http.ResponseWriter, r *http.Request) {
	h.write(w, "alive")
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.write(w, "ok")
}

func (h *HealthHandler) write(w http.ResponseWriter, status string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":    status,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
