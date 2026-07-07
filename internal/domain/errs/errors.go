package errs

import (
	"errors"
	"fmt"
)

var (
	ErrValidation            = errors.New("DomainValidationException")
	ErrIncompatibleComponent = errors.New("IncompatibleComponentException")
	ErrEntityNotFound        = errors.New("EntityNotFoundException")
	ErrInvalidTransition     = errors.New("invalid status transition")
	ErrInvalidScheduledTime  = errors.New("invalid scheduled time")
	ErrDuplicateSelection    = errors.New("duplicate selection")
	ErrMissingRequiredNode   = errors.New("missing required node")

	InvalidTransition    = ErrInvalidTransition
	InvalidScheduledTime = ErrInvalidScheduledTime
)

func Validation(message string) error {
	return fmt.Errorf("%w: %s", ErrValidation, message)
}

func IncompatibleComponent(modelName, detailType, detailName string) error {
	return fmt.Errorf("%w: %s %s is not available for model %s", ErrIncompatibleComponent, detailType, detailName, modelName)
}

func EntityNotFound(entity, id string) error {
	return fmt.Errorf("%w: %s %s", ErrEntityNotFound, entity, id)
}

func DuplicateSelection(detailType string) error {
	return fmt.Errorf("%w: %s", ErrDuplicateSelection, detailType)
}

func MissingRequiredNode(detailType string) error {
	return fmt.Errorf("%w: %s", ErrMissingRequiredNode, detailType)
}
