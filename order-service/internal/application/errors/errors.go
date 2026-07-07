package errors

import "errors"

var (
	ErrNotFound             = errors.New("not found")
	ErrAlreadyExists        = errors.New("already exists")
	ErrInvalidInput         = errors.New("invalid input")
	ErrInvalidPrice         = errors.New("invalid price")
	ErrInvalidStatus        = errors.New("invalid status")
	ErrInvalidTransition    = errors.New("invalid status transition")
	ErrInvalidScheduledTime = errors.New("invalid scheduled time")
	ErrCarNotFound          = errors.New("car not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrOrderNotFound        = errors.New("order not found")
	ErrTestDriveNotFound    = errors.New("test drive request not found")
	ErrPartNotCompatible    = errors.New("part not compatible")
)
