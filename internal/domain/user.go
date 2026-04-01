package domain

import (
	"time"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	ID        string
	Email     string
	Role      UserRole
	CreatedAt time.Time
}
