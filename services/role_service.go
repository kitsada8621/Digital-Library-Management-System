package services

import (
	"dlms/models"
	"dlms/repositories"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IRoleService interface {
	GetRoles() ([]models.Role, error)
	FindById(id string) (models.Role, error)
	FindByName(name string) (models.Role, error)
	CreateRoleMany(roles []models.Role) error
	// AssingRole(user models.User, roles []string) error
}
type RoleServiceImpl struct {
	roleRepository repositories.IRoleRepository
	userService    IUserService
}

func RoleService() IRoleService {
	return &RoleServiceImpl{
		roleRepository: repositories.RoleRepository(),
		userService:    UserService(),
	}
}

func (s *RoleServiceImpl) GetRoles() ([]models.Role, error) {
	findOption := options.Find().SetSort(bson.M{"createdAt": 1})
	roles, err := s.roleRepository.FindAll(bson.D{}, findOption)
	return roles, err
}

func (s *RoleServiceImpl) FindById(id string) (models.Role, error) {
	roleId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Role{}, err
	}
	role, err := s.roleRepository.FindOne(bson.M{"_id": roleId})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return role, fmt.Errorf("Role not found")
		} else {
			return role, err
		}
	}

	return role, nil
}
func (s *RoleServiceImpl) FindByName(name string) (models.Role, error) {
	role, err := s.roleRepository.FindOne(bson.M{"roleName": name})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return role, fmt.Errorf("Role not found")
		} else {
			return role, err
		}
	}
	return role, nil
}
func (s *RoleServiceImpl) CreateRoleMany(roles []models.Role) error {

	var documents []interface{}
	for _, role := range roles {
		documents = append(documents, &role)
	}

	if _, err := s.roleRepository.CreateMany(documents); err != nil {
		return err
	}

	return nil
}

// func (s *RoleServiceImpl) AssingRole(user models.User, roles []string, delRoles []string) error {
// 	// roleIds := []primitive.ObjectID{}
// 	return nil
// }
