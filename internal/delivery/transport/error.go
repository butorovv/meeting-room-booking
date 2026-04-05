package transport

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var (
	ErrInvalidRequest    = ErrorDetail{Code: "INVALID_REQUEST", Message: "invalid request"}
	ErrUnauthorized      = ErrorDetail{Code: "UNAUTHORIZED", Message: "unauthorized"}
	ErrNotFound          = ErrorDetail{Code: "NOT_FOUND", Message: "not found"}
	ErrForbidden         = ErrorDetail{Code: "FORBIDDEN", Message: "access denied"}
	ErrRoomNotFound      = ErrorDetail{Code: "ROOM_NOT_FOUND", Message: "room not found"}
	ErrSlotNotFound      = ErrorDetail{Code: "SLOT_NOT_FOUND", Message: "slot not found"}
	ErrSlotAlreadyBooked = ErrorDetail{Code: "SLOT_ALREADY_BOOKED", Message: "slot is already booked"}
	ErrScheduleExists    = ErrorDetail{Code: "SCHEDULE_EXISTS", Message: "schedule already exists"}
	ErrInternalError     = ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal server error"}
	ErrBookingNotFound   = ErrorDetail{Code: "BOOKING_NOT_FOUND", Message: "booking not found"}
)
