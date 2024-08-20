package repositories

import (
	"context"
	"dlms/configs"
	"dlms/models"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IBookRepository interface {
	FindOne(filter interface{}) (*models.Book, error)
	Aggregate(pipeline interface{}) ([]models.Book, error)
	FindAll(filter interface{}, opts ...*options.FindOptions) (*[]models.Book, error)
	Count(filter interface{}) (int64, error)
	Save(document interface{}) (*mongo.InsertOneResult, error)
	Update(document interface{}, filter interface{}) (*mongo.UpdateResult, error)
	Delete(id primitive.ObjectID) (*mongo.DeleteResult, error)
}

type BookRepositoryImpl struct {
	bookCollection *mongo.Collection
}

func BookRepository() IBookRepository {
	return &BookRepositoryImpl{
		bookCollection: configs.GetCollection(configs.Client, "books"),
	}
}

func (s *BookRepositoryImpl) FindOne(filter interface{}) (*models.Book, error) {
	var book models.Book
	err := s.bookCollection.FindOne(context.Background(), filter).Decode(&book)
	return &book, err
}
func (s *BookRepositoryImpl) Aggregate(pipeline interface{}) ([]models.Book, error) {

	books := []models.Book{}
	ctx := context.Background()
	cur, err := s.bookCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &books); err != nil {
		return nil, err
	}

	fmt.Println("\n\nbooks", books)
	return books, nil
}

func (s *BookRepositoryImpl) FindAll(filter interface{}, opts ...*options.FindOptions) (*[]models.Book, error) {

	ctx := context.Background()
	cur, err := s.bookCollection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	books := []models.Book{}
	if err := cur.All(ctx, &books); err != nil {
		return nil, err
	}

	return &books, nil
}

func (s *BookRepositoryImpl) Count(filter interface{}) (int64, error) {
	count, err := s.bookCollection.CountDocuments(context.Background(), filter)
	return count, err
}

func (s *BookRepositoryImpl) Save(document interface{}) (*mongo.InsertOneResult, error) {
	r, err := s.bookCollection.InsertOne(context.Background(), document)
	return r, err
}

func (s *BookRepositoryImpl) Update(document interface{}, filter interface{}) (*mongo.UpdateResult, error) {
	r, err := s.bookCollection.UpdateOne(context.Background(), filter, document)
	return r, err
}

func (s *BookRepositoryImpl) Delete(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	r, err := s.bookCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return r, err
}
