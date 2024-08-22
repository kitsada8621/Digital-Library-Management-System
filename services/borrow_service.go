package services

import (
	"dlms/dtos"
	"dlms/models"
	"dlms/pkg/utils"
	"dlms/repositories"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IBorrowService interface {
	BorrowingHistory(filter dtos.BorrowedHistoryDto, id string) (*[]models.Borrow, int64, error)
	GetAllBorrowed(filter dtos.GetAllBorrowedDto) ([]models.Borrow, int64, error)
	FindById(id string) (*models.Borrow, error)
	Details(id string) (*models.Borrow, error)
	CreateBorrow(data dtos.BorrowDto, book models.Book, userId primitive.ObjectID) (*models.Borrow, error)
	UpdateBorrow(data dtos.BorrowDto, borrow models.Borrow) (int, error)
	DeleteBorrow(borrow models.Borrow) (int, error)
	CancelBorrowBook(borrow models.Borrow) (int, error)
	ApproveBookBorrowing(borrow models.Borrow) (int, error)
	ReturnBook(borrow models.Borrow) (int, error)
	ApproveBookReturn(borrow models.Borrow) (int, error)
}

type BorrowServiceImpl struct {
	borrowRepository repositories.IBorrowRepository
	bookRepository   repositories.IBookRepository
}

func BorrowService() IBorrowService {
	return &BorrowServiceImpl{
		borrowRepository: repositories.BorrowRepository(),
		bookRepository:   repositories.BookRepository(),
	}
}

func (s *BorrowServiceImpl) BorrowingHistory(filter dtos.BorrowedHistoryDto, id string) (*[]models.Borrow, int64, error) {

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, 0, err
	}

	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"userId", userId}}}},
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "userId"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		bson.D{{"$unwind", "$user"}},
		bson.D{{"$lookup", bson.D{
			{"from", "books"},
			{"localField", "bookId"},
			{"foreignField", "_id"},
			{"as", "book"},
		}}},
		bson.D{{"$unwind", "$book"}},
		bson.D{{"$skip", filter.Skip}},
		bson.D{{"$limit", filter.Limit}},
		bson.D{{"$sort", bson.D{
			{"createdAt", -1},
			{"borrowStatus", 1},
		}}},
	}

	borrowedList, err := s.borrowRepository.Aggregate(pipeline)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.borrowRepository.Count(bson.M{"userId": userId})
	if err != nil {
		return nil, 0, err
	}

	return &borrowedList, count, err
}
func (s *BorrowServiceImpl) GetAllBorrowed(filter dtos.GetAllBorrowedDto) ([]models.Borrow, int64, error) {

	pipeline := mongo.Pipeline{
		bson.D{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "userId"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		bson.D{{"$unwind", "$user"}},
		bson.D{{"$lookup", bson.D{
			{"from", "books"},
			{"localField", "bookId"},
			{"foreignField", "_id"},
			{"as", "book"},
		}}},
		bson.D{{"$unwind", "$book"}},
		bson.D{{"$skip", filter.Skip}},
		bson.D{{"$limit", filter.Limit}},
		bson.D{{"$sort", bson.D{
			{"createdAt", -1},
			{"borrowStatus", 1},
		}}},
	}
	borrows, err := s.borrowRepository.Aggregate(pipeline)
	if err != nil {
		return borrows, 0, err
	}

	count, err := s.borrowRepository.Count(bson.M{})
	if err != nil {
		return borrows, 0, err
	}

	return borrows, count, nil
}

func (s *BorrowServiceImpl) FindById(id string) (*models.Borrow, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	borrow, err := s.borrowRepository.FindOne(bson.M{"_id": _id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("Borrow not found")
		} else {
			return nil, err
		}
	}

	return borrow, nil
}
func (s *BorrowServiceImpl) Details(id string) (*models.Borrow, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	fmt.Println("_id: ", _id)

	borrows, err := s.borrowRepository.Aggregate(mongo.Pipeline{
		bson.D{{"$match", bson.D{{"_id", _id}}}},
		{{"$lookup", bson.D{
			{"from", "users"},
			{"localField", "userId"},
			{"foreignField", "_id"},
			{"as", "user"},
		}}},
		bson.D{{"$unwind", "$user"}},
		{{"$lookup", bson.D{
			{"from", "books"},
			{"localField", "bookId"},
			{"foreignField", "_id"},
			{"as", "book"},
		}}},
		bson.D{{"$unwind", "$book"}},
		bson.D{{"$limit", 1}},
	})
	if err != nil {
		return nil, err
	}

	if len(borrows) == 0 {
		return nil, fmt.Errorf("Borrow not found")
	}

	return &borrows[0], nil
}

func (s *BorrowServiceImpl) CreateBorrow(data dtos.BorrowDto, book models.Book, userId primitive.ObjectID) (*models.Borrow, error) {

	returnDate, err := utils.StringToDate(data.ReturnDate)
	if err != nil {
		return nil, err
	}

	borrow := models.Borrow{
		ID:                 primitive.NewObjectID(),
		UserId:             userId,
		BookId:             book.ID,
		BorrowDate:         time.Now(),
		ScheduleReturnDate: *returnDate,
		BorrowStatus:       1,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if _, err := s.borrowRepository.Save(borrow); err != nil {
		return nil, err
	}

	return &borrow, nil
}

func (s *BorrowServiceImpl) UpdateBorrow(data dtos.BorrowDto, borrow models.Borrow) (int, error) {

	if borrow.BorrowStatus != 1 {
		return http.StatusOK, fmt.Errorf("Unable to edit")
	}

	bookId, err := primitive.ObjectIDFromHex(data.BookId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	book, err := s.bookRepository.FindOne(bson.M{"_id": bookId})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if book.BookStatus != 0 {
		return http.StatusOK, fmt.Errorf("This book has been borrowed")
	}

	returnDate, err := utils.StringToDate(data.ReturnDate)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if _, err := s.borrowRepository.Update(bson.M{
		"$set": bson.M{
			"bookId":     book.ID,
			"returnDate": returnDate,
		},
	}, bson.M{"_id": borrow.ID}); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *BorrowServiceImpl) DeleteBorrow(borrow models.Borrow) (int, error) {

	if borrow.BorrowStatus != 1 && borrow.BorrowStatus != 5 {
		return http.StatusOK, fmt.Errorf("Unable to delete data")
	}

	if _, err := s.borrowRepository.Delete(bson.M{"_id": borrow.ID}); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *BorrowServiceImpl) CancelBorrowBook(borrow models.Borrow) (int, error) {

	if borrow.BorrowStatus == 2 || borrow.BorrowStatus == 4 {
		return http.StatusOK, fmt.Errorf("Cannot be canceled")
	}

	if _, err := s.borrowRepository.Update(bson.M{
		"$set": bson.M{
			"borrowStatus": 5,
			"updatedAt":    time.Now(),
		},
	}, bson.M{"_id": borrow.ID}); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *BorrowServiceImpl) ApproveBookBorrowing(borrow models.Borrow) (int, error) {

	if borrow.BorrowStatus != 1 {
		return http.StatusOK, fmt.Errorf("The book borrowing status is incorrect")
	}

	book, err := s.bookRepository.FindOne(bson.M{"_id": borrow.BookId})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return http.StatusNotFound, fmt.Errorf("Book not found")
		} else {
			return http.StatusInternalServerError, err
		}
	}

	if book.BookStatus == 1 {
		return http.StatusOK, fmt.Errorf("The book has been borrowed.")
	}

	if _, err := s.borrowRepository.Update(bson.M{
		"$set": bson.M{
			"borrowStatus": 2,
			"updatedAt":    time.Now(),
		},
	}, bson.M{"_id": borrow.ID}); err != nil {
		return http.StatusInternalServerError, err
	}

	if _, err := s.bookRepository.Update(bson.M{
		"$set": bson.M{
			"bookStatus":  1,
			"borrowCount": (book.BorrowCount + 1),
		},
	}, bson.M{"_id": book.ID}); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *BorrowServiceImpl) ReturnBook(borrow models.Borrow) (int, error) {

	if borrow.BorrowStatus != 2 {
		return http.StatusOK, fmt.Errorf("The book borrowing status is incorrect")
	}

	if _, err := s.borrowRepository.Update(bson.M{
		"$set": bson.M{
			"borrowStatus": 3,
			"updatedAt":    time.Now(),
		}}, bson.M{"_id": borrow.ID}); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (s *BorrowServiceImpl) ApproveBookReturn(borrow models.Borrow) (int, error) {

	if borrow.BorrowStatus != 3 {
		return http.StatusOK, fmt.Errorf("The book borrowing status is incorrect")
	}

	if _, err := s.borrowRepository.Update(bson.M{
		"$set": bson.M{
			"borrowStatus": 4,
			"returnDate":   time.Now(),
			"updatedAt":    time.Now(),
		}}, bson.M{"_id": borrow.ID}); err != nil {
		return http.StatusInternalServerError, err
	}

	if _, err := s.bookRepository.Update(bson.M{
		"$set": bson.M{
			"bookStatus": 0,
			"updatedAt":  time.Now(),
		},
	}, bson.M{"_id": borrow.BookId}); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
