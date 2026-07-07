package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qxblvde/CarService/internal/domain/errs"
	apperrs "github.com/qxblvde/CarService/order-service/internal/application/errors"
)

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrs.ErrNotFound),
		errors.Is(err, apperrs.ErrCarNotFound),
		errors.Is(err, apperrs.ErrUserNotFound),
		errors.Is(err, apperrs.ErrOrderNotFound),
		errors.Is(err, apperrs.ErrTestDriveNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	case errors.Is(err, errs.ErrValidation),
		errors.Is(err, errs.ErrDuplicateSelection),
		errors.Is(err, errs.ErrMissingRequiredNode):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	case errors.Is(err, errs.ErrInvalidTransition),
		errors.Is(err, errs.ErrInvalidScheduledTime),
		errors.Is(err, errs.ErrIncompatibleComponent),
		errors.Is(err, apperrs.ErrPartNotCompatible):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})

	default:
		log.Printf("unexpected error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
