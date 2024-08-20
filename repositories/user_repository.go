package repositories

import (
	"context"
	"dlms/configs"
	"dlms/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserRepository interface {
	FindAll(filter interface{}, opts ...*options.FindOptions) ([]models.User, error)
	FindOne(filter interface{}) (*models.User, error)
	CountOne(filter interface{}) (int64, error)
	Save(document interface{}) (*mongo.InsertOneResult, error)
	SaveMany(documents []interface{}) (*mongo.InsertManyResult, error)
	Update(document interface{}, filter interface{}) (*mongo.UpdateResult, error)
}
type UserRepositoryImpl struct {
	userCollection *mongo.Collection
}

func UserRepository() IUserRepository {
	return &UserRepositoryImpl{
		userCollection: configs.GetCollection(configs.Client, "users"),
	}
}

func (s *UserRepositoryImpl) FindAll(filter interface{}, opts ...*options.FindOptions) ([]models.User, error) {

	users := []models.User{}

	ctx := context.Background()
	cur, err := s.userCollection.Find(ctx, filter)
	if err != nil {
		return users, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &users); err != nil {
		return users, err
	}

	return users, nil

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

func (s *UserRepositoryImpl) SaveMany(documents []interface{}) (*mongo.InsertManyResult, error) {
	r, err := s.userCollection.InsertMany(context.Background(), documents)
	return r, err
}

func (s *UserRepositoryImpl) Update(document interface{}, filter interface{}) (*mongo.UpdateResult, error) {
	r, err := s.userCollection.UpdateOne(context.Background(), filter, document)
	return r, err
}
