package services

import (
	"dlms/dtos"
	"dlms/models"
	"dlms/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ICategoryService interface {
	GetCategories() ([]models.Category, error)
	FindCategoryById(id string) (*models.Category, error)
	FindCategoryByName(name string) (*models.Category, error)
	CheckDuplicateByName(categoryName string) bool
	CreateCategory(body dtos.CategoryDto) (*models.Category, error)
	UpdateCategory(body dtos.CategoryDto, id primitive.ObjectID) (*mongo.UpdateResult, error)
	DeleteCategory(id primitive.ObjectID) error
}
type CategoryServiceImpl struct {
	categoryRepository repositories.ICategoryRepository
}

func CategoryService() ICategoryService {
	return &CategoryServiceImpl{
		categoryRepository: repositories.CategoryRepository(),
	}
}

func (s *CategoryServiceImpl) FindCategoryById(id string) (*models.Category, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	category, err := s.categoryRepository.FindOne(bson.M{"_id": _id})
	return category, err
}

func (s *CategoryServiceImpl) FindCategoryByName(name string) (*models.Category, error) {
	category, err := s.categoryRepository.FindOne(bson.M{"categoryName": name})
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryServiceImpl) GetCategories() ([]models.Category, error) {
	findOption := options.Find().SetSort(bson.M{"createdAt": 1})
	categories, err := s.categoryRepository.FindAll(bson.M{}, findOption)
	return *categories, err
}

func (s *CategoryServiceImpl) CreateCategory(body dtos.CategoryDto) (*models.Category, error) {

	category := models.Category{
		ID:           primitive.NewObjectID(),
		CategoryName: body.CategoryName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := s.categoryRepository.Save(category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *CategoryServiceImpl) CheckDuplicateByName(categoryName string) bool {
	if count, _ := s.categoryRepository.Count(bson.M{"categoryName": categoryName}); count > 0 {
		return true
	}
	return false
}

func (s *CategoryServiceImpl) UpdateCategory(body dtos.CategoryDto, id primitive.ObjectID) (*mongo.UpdateResult, error) {
	result, err := s.categoryRepository.Update(bson.M{
		"$set": bson.M{"categoryName": body.CategoryName},
	}, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CategoryServiceImpl) DeleteCategory(id primitive.ObjectID) error {
	_, err := s.categoryRepository.Delete(bson.M{"_id": id})
	return err
}
