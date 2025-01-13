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

func (r *Role) SetName(roleName string) {
	switch roleName {
	case string(UserRoleTbz):
		r.Name = UserRoleTbz
	case string(UserRoleGreenEcolution):
		r.Name = UserRoleGreenEcolution
	case string(UserRoleSmarteGrenzregion):
		r.Name = UserRoleSmarteGrenzregion
	default:
		r.Name = UserRoleUnknown
	}
}
