package user

import "time"

type User struct {
	ID        string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
