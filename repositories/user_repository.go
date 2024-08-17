package repositories

import (
	"context"
	"dlms/configs"
	"dlms/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type IUserRepository interface {
	FindOne(filter interface{}) (*models.User, error)
	CountOne(filter interface{}) (int64, error)
	Save(document interface{}) (*mongo.InsertOneResult, error)
}
type UserRepositoryImpl struct {
	userCollection *mongo.Collection
}

func UserRepository() IUserRepository {
	return &UserRepositoryImpl{
		userCollection: configs.GetCollection(configs.Client, "users"),
	}
}

func (s *UserRepositoryImpl) FindOne(filter interface{}) (*models.User, error) {

	var user models.User
	err := s.userCollection.FindOne(context.Background(), filter).Decode(&user)
	return &user, err

}

func (s *UserRepositoryImpl) CountOne(filter interface{}) (int64, error) {
	count, err := s.userCollection.CountDocuments(context.Background(), filter)
	return count, err

}

func (s *UserRepositoryImpl) Save(document interface{}) (*mongo.InsertOneResult, error) {
	r, err := s.userCollection.InsertOne(context.Background(), document)
	return r, err
}
