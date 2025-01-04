package entities

type UserRole string // @Name UserRole

const (
	UserRoleTbz               UserRole = "tbz"
	UserRoleGreenEcolution    UserRole = "green-ecolution"
	UserRoleSmarteGrenzregion UserRole = "smarte-grenzregion"
	UserRoleUnknown           UserRole = "unknown"
)

type RoleResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
} // @Name Role

type RoleListResponse struct {
	Data       []RoleResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
} // @Name RoleList
