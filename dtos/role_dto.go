package dtos

type AssignRoleDto struct {
	UserId string   `json:"userId" binding:"required"`
	Roles  []string `json:"roles" binding:"required"`
}
