package entities

type RoleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
} // @Name Role

type RoleListResponse struct {
	Data       []RoleResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
} // @Name RoleList
