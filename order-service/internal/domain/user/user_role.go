package user

type UserRole string

const (
	UserRoleUser           UserRole = "user"
	UserRoleAdmin          UserRole = "admin"
	UserRoleWarehouseAdmin UserRole = "warehouse_admin"
	UserRoleManager        UserRole = "manager"
)
