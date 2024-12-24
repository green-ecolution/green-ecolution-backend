package entities

type UserRole string

const (
	UserRoleAdmin      UserRole = "Admin"
	UserRoleDriver     UserRole = "Driver"
	UserRoleEngineer   UserRole = "Engineer"
	UserRoleHelper     UserRole = "Helper"
	UserRoleManagement UserRole = "Management"
	UserRoleStudent    UserRole = "Student"
	UserRoleUnknown    UserRole = "unknown"
)

type Role struct {
	ID          int32
	Name        UserRole
	Description string
}

func ParseUserRole(role string) UserRole {
	switch role {
	case string(UserRoleAdmin):
		return UserRoleAdmin
	case string(UserRoleDriver):
		return UserRoleDriver
	case string(UserRoleEngineer):
		return UserRoleEngineer
	case string(UserRoleHelper):
		return UserRoleHelper
	case string(UserRoleManagement):
		return UserRoleManagement
	case string(UserRoleStudent):
		return UserRoleStudent
	default:
		return UserRoleUnknown
	}
}
