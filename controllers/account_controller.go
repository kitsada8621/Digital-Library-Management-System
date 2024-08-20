package controllers

import (
	"dlms/dtos"
	"dlms/models"
	"dlms/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountController struct {
	accountService    services.IAccountService
	userService       services.IUserService
	jwtService        services.IJwtService
	validationService services.IValidationService
}

func NewAccountController() AccountController {
	return AccountController{
		accountService:    services.AccountService(),
		userService:       services.UserService(),
		jwtService:        services.JwtService(),
		validationService: services.ValidationService(),
	}
}

func (s *AccountController) Login(ctx *gin.Context) {

	var data dtos.LoginDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	user, err := s.accountService.Login(data)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ResponseJson{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	token, err := s.jwtService.GenerateToken(*user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	refreshToken, err := s.jwtService.RefreshToken(*user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success:      true,
		Message:      "Successfully logged in",
		Data:         user,
		AccessToken:  token,
		RefreshToken: refreshToken,
	})

}

func (s *AccountController) RefreshToken(ctx *gin.Context) {

	id := ctx.GetString("userId")
	user, err := s.accountService.ProfileInfo(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, models.ResponseJson{
				Success: false,
				Message: "User not found",
				Data:    nil,
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
	}

	token, err := s.jwtService.GenerateToken(*user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	refreshToken, err := s.jwtService.RefreshToken(*user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success:      true,
		Message:      "Successfully refresh token",
		Data:         user,
		AccessToken:  token,
		RefreshToken: refreshToken,
	})

}

func (s *AccountController) Register(ctx *gin.Context) {

	var data dtos.RegisterDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	if ok := s.accountService.EmailAlreadyExist(data.Email); !ok {
		ctx.JSON(http.StatusBadRequest, models.ResponseJson{
			Success: false,
			Message: "Email already exists",
			Data:    nil,
		})
		return
	}

	if ok := s.accountService.UsernameAlreadyExist(data.Username); !ok {
		ctx.JSON(http.StatusBadRequest, models.ResponseJson{
			Success: false,
			Message: "Username already exists",
			Data:    nil,
		})
		return
	}

	if err := s.accountService.Register(data); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseJson{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Registration successful!",
		Data:    nil,
	})
}

func (s *AccountController) Profile(ctx *gin.Context) {

	id := ctx.GetString("userId")
	user, err := s.accountService.ProfileInfo(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, models.ResponseJson{
				Success: false,
				Message: "User not found",
				Data:    nil,
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
				Success: false,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Successfully logged in",
		Data:    user,
	})
}
