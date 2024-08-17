package services

import (
	"dlms/repositories"
)

type IUserService interface {
}
type UserServiceImpl struct {
	userRepository repositories.IUserRepository
}

func UserService() IUserService {
	return &UserServiceImpl{
		userRepository: repositories.UserRepository(),
	}
}
