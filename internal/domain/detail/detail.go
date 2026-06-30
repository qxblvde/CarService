package detail

import (
	"time"

	"github.com/qxblvde/CarService/internal/domain/vo"
)

type Detail struct {
	ID        string
	Type      DetailType
	Price     vo.Money
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
