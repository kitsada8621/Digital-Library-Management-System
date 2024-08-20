package services

import (
	"dlms/dtos"
	"dlms/models"
	"dlms/repositories"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserService interface {
	GetUsers(data dtos.GetUsersDto) ([]models.User, int, error)
	FindByUsername(username string) (*models.User, error)
	FindById(id string) (*models.User, error)
	CreateMany(user []models.User) error
	AssignRoles(data dtos.AssignRoleDto) (int, error)
}
type UserServiceImpl struct {
	userRepository repositories.IUserRepository
}

func UserService() IUserService {
	return &UserServiceImpl{
		userRepository: repositories.UserRepository(),
	}
}

func (s *UserServiceImpl) GetUsers(data dtos.GetUsersDto) ([]models.User, int, error) {

	filter := bson.M{"name": bson.D{{"$regex", data.Search}, {"$options", "i"}}}
	findOptions := options.Find().SetSort(bson.M{"updatedAt": -1}).SetSkip(int64(*data.Skip)).SetLimit(int64(*data.Limit))
	users, err := s.userRepository.FindAll(filter, findOptions)
	if err != nil {
		return users, 0, err
	}

	count, err := s.userRepository.CountOne(filter)
	if err != nil {
		return users, 0, err
	}

	return users, int(count), nil
}

func (s *UserServiceImpl) FindByUsername(username string) (*models.User, error) {

	user, err := s.userRepository.FindOne(bson.M{"username": username})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("User not found")
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (s *UserServiceImpl) FindById(id string) (*models.User, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.FindOne(bson.M{"_id": _id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("User not found")
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (s *UserServiceImpl) CreateMany(users []models.User) error {

	if len(users) == 0 {
		return fmt.Errorf("users is empty")
	}

	var documents []interface{}
	for _, row := range users {
		documents = append(documents, row)
	}

	if _, err := s.userRepository.SaveMany(documents); err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) AssignRoles(data dtos.AssignRoleDto) (int, error) {

	userId, err := primitive.ObjectIDFromHex(data.UserId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	user, err := s.userRepository.FindOne(bson.M{"_id": userId})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return http.StatusNotFound, err
		} else {
			return http.StatusInternalServerError, err
		}
	}

	roles := []primitive.ObjectID{}
	for _, row := range data.Roles {
		if roleId, err := primitive.ObjectIDFromHex(row); err == nil {
			roles = append(roles, roleId)
		}
	}

	if _, err := s.userRepository.Update(bson.M{
		"$set": bson.M{
			"roles": roles,
		},
	}, bson.M{"_id": user.ID}); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
