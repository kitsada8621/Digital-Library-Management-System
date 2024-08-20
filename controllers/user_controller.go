package controllers

import (
	"dlms/dtos"
	"dlms/models"
	"dlms/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService       services.IUserService
	validationService services.IValidationService
}

func NewUserController() UserController {
	return UserController{
		userService:       services.UserService(),
		validationService: services.ValidationService(),
	}
}

func (s *UserController) GetUsers(ctx *gin.Context) {

	var data dtos.GetUsersDto
	if err := ctx.ShouldBindQuery(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	rows, count, err := s.userService.GetUsers(data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "User data list",
		Data:    rows,
		Total:   &count,
	})
}
