package transport

import (
	"fmt"
	"time"

	"github.com/butorovv/meeting-room-booking/internal/domain"
)

type CreateRoomRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Capacity    *int    `json:"capacity,omitempty"`
}

func (r *CreateRoomRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if r.Capacity != nil && *r.Capacity <= 0 {
		return fmt.Errorf("capacity must be positive")
	}
	return nil
}

type RoomResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Capacity    *int       `json:"capacity,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
}

func ToRoomResponse(r *domain.Room) *RoomResponse {
	if r == nil {
		return nil
	}
	return &RoomResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Capacity:    r.Capacity,
		CreatedAt:   &r.CreatedAt,
	}
}

func ToRoomResponseList(rooms []*domain.Room) []RoomResponse {
	if rooms == nil {
		return []RoomResponse{}
	}
	result := make([]RoomResponse, 0, len(rooms))
	for _, r := range rooms {
		result = append(result, *ToRoomResponse(r))
	}
	return result
}
