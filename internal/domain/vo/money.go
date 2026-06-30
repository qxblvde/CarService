package vo

import "errors"

type Money struct {
	Amount int64
}

func NewMoney(amount int64) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount must be non-negative")
	}
	return Money{Amount: amount}, nil
}

func (m Money) Add(amount int64) (Money, error) {
	if amount < 0 {
		return m, errors.New("amount must be positive")
	}
	return Money{Amount: m.Amount + amount}, nil
}
