package entities

type UserRole string // @Name UserRole

const (
	UserRoleAdmin      UserRole = "Admin"
	UserRoleDriver     UserRole = "Driver"
	UserRoleEngineer   UserRole = "Engineer"
	UserRoleHelper     UserRole = "Helper"
	UserRoleManagement UserRole = "Management"
	UserRoleStudent    UserRole = "Student"
	UserRoleUnknown    UserRole = "unknown"
)

type RoleResponse struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
} // @Name Role

type RoleListResponse struct {
	Data       []RoleResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
} // @Name RoleList
