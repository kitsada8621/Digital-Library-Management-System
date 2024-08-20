package services

import (
	"dlms/dtos"
	"dlms/models"
	"dlms/repositories"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IBookService interface {
	MostBorrowedBooks(filter dtos.GetBooksDto) ([]models.Book, int64, error)
	FindById(id string) (*models.Book, error)
	BookStatus(id string) (bool, error)
	FindByName(name string) (*models.Book, error)
	GetBooks(filter dtos.GetBooksDto) (*[]models.Book, int64, error)
	Details(id string) (*models.Book, error)
	CreateBook(body dtos.BookDto, userId string) (*models.Book, error)
	UpdateBook(body dtos.BookDto, book models.Book, userId string) (int, error)
	DeleteBook(id primitive.ObjectID) error
}
type BookServiceImpl struct {
	bookRepository repositories.IBookRepository
}

func BookService() IBookService {
	return &BookServiceImpl{
		bookRepository: repositories.BookRepository(),
	}
}

func (s *BookServiceImpl) MostBorrowedBooks(filter dtos.GetBooksDto) ([]models.Book, int64, error) {

	match := bson.D{
		{"$or", bson.A{
			bson.D{{"bookTitle", bson.D{{"$regex", filter.Search}, {"$options", "i"}}}},
			bson.D{{"author", bson.D{{"$regex", filter.Search}, {"$options", "i"}}}},
		}},
	}

	if len(filter.Categories) > 0 {
		categories := []primitive.ObjectID{}
		for _, categoryId := range filter.Categories {
			if id, err := primitive.ObjectIDFromHex(categoryId); err == nil {
				categories = append(categories, id)
			}
		}
		match = append(match, bson.E{"categoryId", bson.M{"$in": categories}})
	}

	pipeline := mongo.Pipeline{
		bson.D{{"$match", match}},
		bson.D{{"$lookup", bson.D{
			{"from", "categories"},
			{"localField", "categoryId"},
			{"foreignField", "_id"},
			{"as", "category"},
		}}},
		bson.D{{"$unwind", "$category"}},
		bson.D{{"$skip", filter.Skip}},
		bson.D{{"$limit", filter.Limit}},
		bson.D{{"$sort", bson.D{
			{"borrowCount", -1},
			{"updatedAt", -1},
		}}},
	}

	rows, err := s.bookRepository.Aggregate(pipeline)
	if err != nil {
		return rows, 0, err
	}

	count, err := s.bookRepository.Count(match)
	if err != nil {
		return rows, 0, err
	}

	return rows, count, err
}

func (s *BookServiceImpl) FindById(id string) (*models.Book, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	book, err := s.bookRepository.FindOne(bson.M{"_id": _id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("Book not found")
		} else {
			return nil, err
		}
	}
	return book, nil
}

func (s *BookServiceImpl) BookStatus(id string) (bool, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	count, err := s.bookRepository.Count(bson.M{"_id": _id})
	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, nil
	}

	return true, nil
}

func (s *BookServiceImpl) FindByName(name string) (*models.Book, error) {
	book, err := s.bookRepository.FindOne(bson.M{"bookTitle": name})
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *BookServiceImpl) GetBooks(filter dtos.GetBooksDto) (*[]models.Book, int64, error) {

	match := bson.D{
		{"$or", bson.A{
			bson.D{{"bookTitle", bson.D{{"$regex", filter.Search}, {"$options", "i"}}}},
			bson.D{{"author", bson.D{{"$regex", filter.Search}, {"$options", "i"}}}},
		}},
	}
	if len(filter.Categories) > 0 {
		categories := []primitive.ObjectID{}
		for _, categoryId := range filter.Categories {
			if id, err := primitive.ObjectIDFromHex(categoryId); err == nil {
				categories = append(categories, id)
			}
		}
		match = append(match, bson.E{"categoryId", bson.M{"$in": categories}})
	}

	pipeline := mongo.Pipeline{
		bson.D{{"$match", match}},
		bson.D{{"$lookup", bson.D{
			{"from", "categories"},
			{"localField", "categoryId"},
			{"foreignField", "_id"},
			{"as", "category"},
		}}},
		bson.D{{"$unwind", "$category"}},
		bson.D{{"$skip", filter.Skip}},
		bson.D{{"$limit", filter.Limit}},
		bson.D{{"$sort", bson.D{
			{"updatedAt", -1},
			{"createdAt", -1},
		}}},
	}

	books, err := s.bookRepository.Aggregate(pipeline)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.bookRepository.Count(match)
	if err != nil {
		return nil, 0, err
	}

	return &books, count, nil
}

func (s *BookServiceImpl) Details(id string) (*models.Book, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"_id", _id}}}},
		bson.D{{"$lookup", bson.D{
			{"from", "categories"},
			{"localField", "categoryId"},
			{"foreignField", "_id"},
			{"as", "category"},
		}}},
		bson.D{{"$unwind", "$category"}},
		bson.D{{"$limit", 1}},
	}
	books, err := s.bookRepository.Aggregate(pipeline)
	if err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return nil, fmt.Errorf("Book not found")
	}

	return &books[0], nil
}

func (s *BookServiceImpl) CreateBook(body dtos.BookDto, userId string) (*models.Book, error) {

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	categoryId, err := primitive.ObjectIDFromHex(body.CategoryId)
	if err != nil {
		return nil, err
	}

	book := models.Book{
		ID:              primitive.NewObjectID(),
		CreatedBy:       id,
		CategoryId:      categoryId,
		Author:          body.Author,
		BookTitle:       body.BookTitle,
		BookDescription: body.BookDescription,
		BorrowCount:     0,
		BookStatus:      0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if _, err := s.bookRepository.Save(book); err != nil {
		return nil, err
	}

	return &book, nil
}

func (s *BookServiceImpl) UpdateBook(body dtos.BookDto, book models.Book, userId string) (int, error) {

	categoryId, err := primitive.ObjectIDFromHex(body.CategoryId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// if book.CreatedBy.Hex() != userId {
	// 	return http.StatusOK, fmt.Errorf("Unable to edit data")
	// }

	if _, err := s.bookRepository.Update(bson.M{
		"$set": bson.M{
			"categoryId": categoryId,
			"author":     body.Author,
			"bookTitle":  body.BookTitle,
			"bookDesc":   body.BookDescription,
			"updatedAt":  time.Now(),
		},
	}, bson.M{"_id": book.ID}); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *BookServiceImpl) DeleteBook(id primitive.ObjectID) error {
	_, err := s.bookRepository.Delete(id)
	return err
}
