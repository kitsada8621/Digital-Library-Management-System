package controllers

import (
	"dlms/dtos"
	"dlms/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountController struct {
	accountService services.IAccountService
	userService    services.IUserService
	jwtService     services.IJwtService
}

func NewAccountController() AccountController {
	return AccountController{
		accountService: services.AccountService(),
		userService:    services.UserService(),
		jwtService:     services.JwtService(),
	}
}

func (s *AccountController) Login(ctx *gin.Context) {

	var data dtos.LoginDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			`success`: false,
			`message`: "Unprocessed Entities",
			`data`:    err.Error(),
		})
		return
	}

	user, err := s.accountService.Login(data)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			`success`: false,
			`message`: err.Error(),
			`data`:    nil,
		})
		return
	}

	token, err := s.jwtService.GenerateToken(user.ID.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			`success`: false,
			`message`: err.Error(),
			`data`:    nil,
		})
		return
	}

	refreshToken, err := s.jwtService.RefreshToken(user.ID.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			`success`: false,
			`message`: err.Error(),
			`data`:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		`success`:       true,
		`message`:       "Successfully logged in",
		`data`:          user,
		`access_token`:  token,
		`refresh_token`: refreshToken,
	})

}

func (s *AccountController) RefreshToken(ctx *gin.Context) {

	id := ctx.GetString("userId")
	user, err := s.accountService.ProfileInfo(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{
				`success`: false,
				`message`: "User not found",
				`data`:    nil,
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				`success`: false,
				`message`: err.Error(),
				`data`:    nil,
			})
			return
		}
	}

	token, err := s.jwtService.GenerateToken(user.ID.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			`success`: false,
			`message`: err.Error(),
			`data`:    nil,
		})
		return
	}

	refreshToken, err := s.jwtService.RefreshToken(user.ID.Hex())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			`success`: false,
			`message`: err.Error(),
			`data`:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		`success`:       true,
		`message`:       "Successfully refresh token",
		`data`:          user,
		`access_token`:  token,
		`refresh_token`: refreshToken,
	})

}

func (s *AccountController) Register(ctx *gin.Context) {

	var data dtos.RegisterDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			`success`: false,
			`message`: "Unprocessed Entities",
			`data`:    err.Error(),
		})
		return
	}

	if ok := s.accountService.EmailAlreadyExist(data.Email); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			`success`: false,
			`message`: "Email already exists",
			`data`:    nil,
		})
		return
	}

	if ok := s.accountService.UsernameAlreadyExist(data.Username); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			`success`: false,
			`message`: "Username already exists",
			`data`:    nil,
		})
		return
	}

	if err := s.accountService.Register(data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			`success`: false,
			`message`: err.Error(),
			`data`:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		`success`: true,
		`message`: "Registration successful!",
		`data`:    nil,
	})
}

func (s *AccountController) Profile(ctx *gin.Context) {

	id := ctx.GetString("userId")
	user, err := s.accountService.ProfileInfo(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{
				`success`: false,
				`message`: "User not found",
				`data`:    nil,
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				`success`: false,
				`message`: err.Error(),
				`data`:    nil,
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		`success`: true,
		`message`: "Successfully logged in",
		`data`:    user,
	})
}
