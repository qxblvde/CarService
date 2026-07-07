package detail

import (
	"time"

	"github.com/qxblvde/CarService/internal/domain/errs"
	"github.com/qxblvde/CarService/order-service/internal/domain/vo"
)

type Detail struct {
	ID                    string
	Name                  string
	Type                  DetailType
	Price                 vo.Money
	CompatibleCarModelIDs []string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
}

func NewDetail(id, name string, detailType DetailType, price vo.Money, compatibleCarModelIDs ...string) (Detail, error) {
	if id == "" {
		return Detail{}, errs.Validation("detail id is required")
	}
	if name == "" {
		return Detail{}, errs.Validation("detail name is required")
	}
	if !detailType.IsValid() {
		return Detail{}, errs.Validation("invalid detail type")
	}

	now := time.Now().UTC()
	modelIDs := uniqueNonEmptyStrings(compatibleCarModelIDs)
	return Detail{
		ID:                    id,
		Name:                  name,
		Type:                  detailType,
		Price:                 price,
		CompatibleCarModelIDs: modelIDs,
		CreatedAt:             now,
		UpdatedAt:             now,
	}, nil
}

func (d Detail) CompatibleWith(modelID string) bool {
	for _, id := range d.CompatibleCarModelIDs {
		if id == modelID {
			return true
		}
	}
	return false
}

func uniqueNonEmptyStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
