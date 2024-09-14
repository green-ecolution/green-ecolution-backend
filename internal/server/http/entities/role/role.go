package role

import "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/pagination"

type RoleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
} // @Name Role

type RoleListResponse struct {
	Data       []RoleResponse        `json:"data"`
	Pagination pagination.Pagination `json:"pagination"`
} // @Name RoleList
