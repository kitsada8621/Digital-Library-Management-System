package repositories

import (
	"context"
	"dlms/configs"
	"dlms/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IRoleRepository interface {
	FindAll(filter interface{}, opts ...*options.FindOptions) ([]models.Role, error)
	FindOne(filter interface{}, opts ...*options.FindOneOptions) (models.Role, error)
	CreateMany(documents []interface{}) (*mongo.InsertManyResult, error)
	Create(document interface{}) (*mongo.InsertOneResult, error)
	Update(document interface{}, filter interface{}) (*mongo.UpdateResult, error)
	Delete(filter interface{}) (*mongo.DeleteResult, error)
}
type RoleRepositoryImpl struct {
	roleCollection *mongo.Collection
}

func RoleRepository() IRoleRepository {
	return &RoleRepositoryImpl{
		roleCollection: configs.GetCollection(configs.Client, "roles"),
	}
}

func (s *RoleRepositoryImpl) FindAll(filter interface{}, opts ...*options.FindOptions) ([]models.Role, error) {

	roles := []models.Role{}
	ctx := context.Background()
	cur, err := s.roleCollection.Find(ctx, filter, opts...)
	if err != nil {
		return roles, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &roles); err != nil {
		return roles, err
	}

	return roles, nil
}

func (s *RoleRepositoryImpl) FindOne(filter interface{}, opts ...*options.FindOneOptions) (models.Role, error) {
	role := models.Role{}
	if err := s.roleCollection.FindOne(context.Background(), filter, opts...).Decode(&role); err != nil {
		return role, err
	}
	return role, nil
}

func (s *RoleRepositoryImpl) CreateMany(documents []interface{}) (*mongo.InsertManyResult, error) {
	r, err := s.roleCollection.InsertMany(context.Background(), documents)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *RoleRepositoryImpl) Create(document interface{}) (*mongo.InsertOneResult, error) {
	r, err := s.roleCollection.InsertOne(context.Background(), document)
	return r, err
}
func (s *RoleRepositoryImpl) Update(document interface{}, filter interface{}) (*mongo.UpdateResult, error) {
	r, err := s.roleCollection.UpdateOne(context.Background(), filter, document)
	return r, err
}
func (s *RoleRepositoryImpl) Delete(filter interface{}) (*mongo.DeleteResult, error) {
	r, err := s.roleCollection.DeleteOne(context.Background(), filter)
	return r, err
}
