package controllers

import (
	"dlms/dtos"
	"dlms/models"
	"dlms/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookController struct {
	bookService       services.IBookService
	validationService services.IValidationService
}

func NewBookController() BookController {
	return BookController{
		bookService:       services.BookService(),
		validationService: services.ValidationService(),
	}
}

func (s *BookController) MostBorrowedBooks(ctx *gin.Context) {
	var data dtos.GetBooksDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	rows, count, err := s.bookService.MostBorrowedBooks(data)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	total := int(count)
	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "The most borrowed books",
		Data:    rows,
		Total:   &total,
	})
}

func (s *BookController) GetBooks(ctx *gin.Context) {

	var data dtos.GetBooksDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	books, count, err := s.bookService.GetBooks(data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	total := int(count)
	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "List of all books",
		Data:    books,
		Total:   &total,
	})
}

func (s *BookController) BookDetails(ctx *gin.Context) {

	id := ctx.Param("id")
	book, err := s.bookService.Details(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, models.ResponseJson{
				Success: false,
				Message: "Book not found",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
				Success: false,
				Message: err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Successfully retrieved book data",
		Data:    book,
	})
}

func (s *BookController) CreateBook(ctx *gin.Context) {

	var data dtos.BookDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	if book, _ := s.bookService.FindByName(data.BookTitle); book != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseJson{
			Success: false,
			Message: "This book already exists",
		})
		return
	}

	userId := ctx.GetString("userId")
	if _, err := s.bookService.CreateBook(data, userId); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Book created successfully",
		Data:    nil,
	})
}

func (s *BookController) EditBook(ctx *gin.Context) {

	id := ctx.Param("id")
	book, err := s.bookService.FindById(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, models.ResponseJson{
				Success: false,
				Message: "Book not found",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
				Success: false,
				Message: err.Error(),
			})
			return
		}
	}

	// Masking data
	book.Category = nil

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Successfully retrieved book data",
		Data:    book,
	})
}

func (s *BookController) UpdateBook(ctx *gin.Context) {

	var data dtos.BookDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	id := ctx.Param("id")
	book, err := s.bookService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if b, _ := s.bookService.FindByName(data.BookTitle); b != nil {
		if b.ID.Hex() != book.ID.Hex() {
			ctx.JSON(http.StatusBadRequest, models.ResponseJson{
				Success: false,
				Message: "This book already exists",
			})
			return
		}
	}

	userId := ctx.GetString("userId")
	if httpStatus, err := s.bookService.UpdateBook(data, *book, userId); err != nil {
		ctx.JSON(httpStatus, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Book updated successfully",
		Data:    nil,
	})
}

func (s *BookController) DeleteBook(ctx *gin.Context) {

	id := ctx.Param("id")
	book, err := s.bookService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if book.BookStatus == 1 {
		ctx.JSON(http.StatusOK, models.ResponseJson{
			Success: false,
			Message: "Cannot be deleted",
		})
		return
	}

	if err := s.bookService.DeleteBook(book.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Book deleted successfully",
		Data:    nil,
	})
}
