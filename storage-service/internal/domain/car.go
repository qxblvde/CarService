package domain

import "errors"

var ErrCarNotFound = errors.New("car not found")

type Car struct {
	ID        string
	Brand     string
	Model     string
	Year      int
	Price     int64
	Available bool
}
