package repositories

import (
	"context"
	"dlms/configs"
	"dlms/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ICategoryRepository interface {
	FindAll(filter interface{}, opts ...*options.FindOptions) (*[]models.Category, error)
	FindOne(filter interface{}) (*models.Category, error)
	Save(document interface{}) (*mongo.InsertOneResult, error)
	Update(document interface{}, filter interface{}) (*mongo.UpdateResult, error)
	Delete(filter interface{}) (*mongo.DeleteResult, error)
	Count(filter interface{}) (int64, error)
}
type CategoryRepositoryImpl struct {
	categoryCollection *mongo.Collection
}

func CategoryRepository() ICategoryRepository {
	return &CategoryRepositoryImpl{
		categoryCollection: configs.GetCollection(configs.Client, "categories"),
	}
}

func (s *CategoryRepositoryImpl) FindAll(filter interface{}, opts ...*options.FindOptions) (*[]models.Category, error) {

	ctx := context.Background()
	cur, err := s.categoryCollection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	categories := []models.Category{}
	if err := cur.All(ctx, &categories); err != nil {
		return nil, err
	}

	return &categories, nil
}

func (s *CategoryRepositoryImpl) FindOne(filter interface{}) (*models.Category, error) {
	var category models.Category
	err := s.categoryCollection.FindOne(context.Background(), filter).Decode(&category)
	return &category, err
}

func (s *CategoryRepositoryImpl) Save(document interface{}) (*mongo.InsertOneResult, error) {
	r, err := s.categoryCollection.InsertOne(context.Background(), document)
	return r, err
}

func (s *CategoryRepositoryImpl) Update(document interface{}, filter interface{}) (*mongo.UpdateResult, error) {
	r, err := s.categoryCollection.UpdateOne(context.Background(), filter, document)
	return r, err
}

func (s *CategoryRepositoryImpl) Delete(filter interface{}) (*mongo.DeleteResult, error) {
	r, err := s.categoryCollection.DeleteOne(context.Background(), filter)
	return r, err
}

func (s *CategoryRepositoryImpl) Count(filter interface{}) (int64, error) {
	count, err := s.categoryCollection.CountDocuments(context.Background(), filter)
	return count, err
}
