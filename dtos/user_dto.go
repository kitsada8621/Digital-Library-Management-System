package dtos

type GetUsersDto struct {
	Skip   *int    `form:"skip" binding:"required"`
	Limit  *int    `form:"limit" binding:"required"`
	Search *string `form:"search" binding:"required"`
}
