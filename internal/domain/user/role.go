package user

type Role string

const (
	User           Role = "user"
	Admin          Role = "admin"
	WarehouseAdmin Role = "warehouse_admin"
	Manager        Role = "manager"
)
