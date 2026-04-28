package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrWeakPassword      = errors.New("weak password")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidToken      = errors.New("invalid token")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternalServer    = errors.New("internal server error")
	ErrHTTPMethod        = errors.New("method not allowed")
	ErrRequestParams     = errors.New("invalid request parameters")

	ErrRoomNotFound       = errors.New("room not found")
	ErrScheduleExists     = errors.New("schedule already exists")
	ErrScheduleNotFound   = errors.New("schedule not found")
	ErrInvalidBookingTime = errors.New("invalid booking time")
	ErrSlotNotFound       = errors.New("slot not found")
	ErrSlotInPast         = errors.New("slot is in the past")
	ErrSlotAlreadyBooked  = errors.New("slot is already booked")
	ErrBookingNotFound    = errors.New("booking not found")
)
