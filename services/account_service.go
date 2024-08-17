package services

import (
	"dlms/dtos"
	"dlms/models"
	"dlms/repositories"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type IAccountService interface {
	EmailAlreadyExist(email string) bool
	UsernameAlreadyExist(email string) bool
	Login(body dtos.LoginDto) (*models.User, error)
	Register(body dtos.RegisterDto) error
	ProfileInfo(id string) (*models.User, error)
}
type AccountServiceImpl struct {
	userRepository repositories.IUserRepository
}

func AccountService() IAccountService {
	return &AccountServiceImpl{
		userRepository: repositories.UserRepository(),
	}
}

func (s *AccountServiceImpl) EmailAlreadyExist(email string) bool {
	if count, _ := s.userRepository.CountOne(bson.M{"email": email}); count > 0 {
		return false
	}
	return true
}

func (s *AccountServiceImpl) UsernameAlreadyExist(username string) bool {
	if count, _ := s.userRepository.CountOne(bson.M{"username": username}); count > 0 {
		return false
	}
	return true
}

func (s *AccountServiceImpl) Login(body dtos.LoginDto) (*models.User, error) {

	user, err := s.userRepository.FindOne(bson.M{"username": body.Username})
	if err != nil {
		return nil, fmt.Errorf("invalid username")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Masking data
	user.Password = ""
	user.Token = nil

	return user, nil
}

func (s *AccountServiceImpl) Register(body dtos.RegisterDto) error {

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if _, err := s.userRepository.Save(&models.User{
		ID:        primitive.NewObjectID(),
		Name:      body.Name,
		Email:     body.Email,
		Phone:     body.Phone,
		Username:  body.Username,
		Password:  string(passwordHashed),
		CreatedAt: time.Now(),
		UpdatedAT: time.Now(),
	}); err != nil {
		return err
	}

	return nil
}

func (s *AccountServiceImpl) ProfileInfo(id string) (*models.User, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.FindOne(bson.M{"_id": _id})
	if err != nil {
		return nil, err
	}

	// Masking data
	user.Password = ""
	user.Token = nil

	return user, nil
}
