package controllers

import (
	"dlms/dtos"
	"dlms/models"
	"dlms/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleService       services.IRoleService
	userService       services.IUserService
	validationService services.IValidationService
}

func NewRoleController() RoleController {
	return RoleController{
		roleService:       services.RoleService(),
		userService:       services.UserService(),
		validationService: services.ValidationService(),
	}
}

func (s *RoleController) GetRoles(ctx *gin.Context) {
	roles, err := s.roleService.GetRoles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "List of all roles",
		Data:    roles,
	})
}

func (s *RoleController) AssignRole(ctx *gin.Context) {

	var data dtos.AssignRoleDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	if httpStatus, err := s.userService.AssignRoles(data); err != nil {
		ctx.JSON(httpStatus, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "User role setup successfully",
	})
}
