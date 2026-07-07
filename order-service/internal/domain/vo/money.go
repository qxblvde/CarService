package vo

import "github.com/qxblvde/CarService/internal/domain/errs"

type Money struct {
	Amount int64
}

func NewMoney(amount int64) (Money, error) {
	if amount < 0 {
		return Money{}, errs.Validation("money amount must be non-negative")
	}
	return Money{Amount: amount}, nil
}

func MustMoney(amount int64) Money {
	return Money{Amount: amount}
}

func (m Money) Add(amount int64) Money {
	return Money{Amount: m.Amount + amount}
}

func (m Money) Sub(amount int64) Money {
	return Money{Amount: m.Amount - amount}
}

func (m Money) IsNegative() bool {
	return m.Amount < 0
}

func (m Money) IsZero() bool {
	return m.Amount == 0
}

func (m Money) Int64() int64 {
	return m.Amount
}
