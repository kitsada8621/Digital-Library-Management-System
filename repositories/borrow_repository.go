package repositories

import (
	"context"
	"dlms/configs"
	"dlms/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type IBorrowRepository interface {
	FindOne(filter interface{}) (*models.Borrow, error)
	Count(filter interface{}) (int64, error)
	Save(document interface{}) (*mongo.InsertOneResult, error)
	Update(document interface{}, filter interface{}) (*mongo.UpdateResult, error)
	Delete(filter interface{}) (*mongo.DeleteResult, error)
	Aggregate(pipeline interface{}) ([]models.Borrow, error)
}

type BorrowRepositoryImpl struct {
	borrowCollection *mongo.Collection
}

func BorrowRepository() IBorrowRepository {
	return &BorrowRepositoryImpl{
		borrowCollection: configs.GetCollection(configs.Client, "borrows"),
	}
}

func (s *BorrowRepositoryImpl) FindOne(filter interface{}) (*models.Borrow, error) {
	var borrow models.Borrow
	if err := s.borrowCollection.FindOne(context.Background(), filter).Decode(&borrow); err != nil {
		return nil, err
	}
	return &borrow, nil
}

func (s *BorrowRepositoryImpl) Count(filter interface{}) (int64, error) {
	count, err := s.borrowCollection.CountDocuments(context.Background(), filter)
	return count, err
}

func (s *BorrowRepositoryImpl) Save(document interface{}) (*mongo.InsertOneResult, error) {
	r, err := s.borrowCollection.InsertOne(context.Background(), document)
	return r, err
}

func (s *BorrowRepositoryImpl) Update(document interface{}, filter interface{}) (*mongo.UpdateResult, error) {
	r, err := s.borrowCollection.UpdateOne(context.Background(), filter, document)
	return r, err
}

func (s *BorrowRepositoryImpl) Delete(filter interface{}) (*mongo.DeleteResult, error) {
	r, err := s.borrowCollection.DeleteOne(context.Background(), filter)
	return r, err
}

func (s *BorrowRepositoryImpl) Aggregate(pipeline interface{}) ([]models.Borrow, error) {
	borrows := []models.Borrow{}
	ctx := context.Background()
	cur, err := s.borrowCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return borrows, err
	}

	if err := cur.All(ctx, &borrows); err != nil {
		return borrows, err
	}

	return borrows, nil

}
