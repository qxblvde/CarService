package domain

import (
	"fmt"
	"time"

	"github.com/qxblvde/CarService/internal/domain/errs"
)

type AssemblyStatus string

const (
	StatusCreated   AssemblyStatus = "CREATED"
	StatusAssembled AssemblyStatus = "ASSEMBLED"
	StatusFail      AssemblyStatus = "FAIL"
)

type SourceOrderType string

const (
	OrderTypeInStock SourceOrderType = "IN_STOCK"
	OrderTypeCustom  SourceOrderType = "CUSTOM"
)

type AssemblyOrder struct {
	ID            string
	SourceOrderID string
	OrderType     SourceOrderType
	Status        AssemblyStatus
	Removed       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewAssemblyOrder(id, sourceOrderID string, orderType SourceOrderType) (AssemblyOrder, error) {
	if id == "" {
		return AssemblyOrder{}, errs.Validation("id is required")
	}
	if sourceOrderID == "" {
		return AssemblyOrder{}, errs.Validation("sourceOrderId is required")
	}

	now := time.Now().UTC()
	return AssemblyOrder{
		ID:            id,
		SourceOrderID: sourceOrderID,
		OrderType:     orderType,
		Status:        StatusCreated,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (o *AssemblyOrder) Assemble() error {
	if o.Status != StatusCreated {
		return fmt.Errorf("%w: %s -> %s", errs.ErrInvalidTransition, o.Status, StatusAssembled)
	}
	o.Status = StatusAssembled
	o.UpdatedAt = time.Now().UTC()
	return nil
}

func (o *AssemblyOrder) Fail() error {
	if o.Status != StatusCreated {
		return fmt.Errorf("%w: %s -> %s", errs.ErrInvalidTransition, o.Status, StatusFail)
	}
	o.Status = StatusFail
	o.UpdatedAt = time.Now().UTC()
	return nil
}

func (o *AssemblyOrder) Remove() {
	o.Removed = true
	o.UpdatedAt = time.Now().UTC()
}
