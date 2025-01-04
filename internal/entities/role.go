package entities

type UserRole string

const (
	UserRoleTbz               UserRole = "tbz"
	UserRoleGreenEcolution    UserRole = "green-ecolution"
	UserRoleSmarteGrenzregion UserRole = "smarte-grenzregion"
	UserRoleUnknown           UserRole = "unknown"
)

type Role struct {
	ID   int32
	Name UserRole
}

func ParseUserRole(role string) UserRole {
	switch role {
	case string(UserRoleTbz):
		return UserRoleTbz
	case string(UserRoleGreenEcolution):
		return UserRoleGreenEcolution
	case string(UserRoleSmarteGrenzregion):
		return UserRoleSmarteGrenzregion
	default:
		return UserRoleUnknown
	}
}
