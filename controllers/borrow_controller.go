package controllers

import (
	"dlms/dtos"
	"dlms/models"
	"dlms/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BorrowController struct {
	validationService services.IValidationService
	userService       services.IUserService
	bookService       services.IBookService
	borrowService     services.IBorrowService
}

func NewBorrowController() BorrowController {
	return BorrowController{
		userService:       services.UserService(),
		validationService: services.ValidationService(),
		bookService:       services.BookService(),
		borrowService:     services.BorrowService(),
	}
}

func (s *BorrowController) BorrowingHistory(ctx *gin.Context) {

	var data dtos.BorrowedHistoryDto
	if err := ctx.ShouldBindQuery(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	userId := ctx.GetString("userId")
	rows, count, err := s.borrowService.BorrowingHistory(data, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	total := int(count)
	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "My book borrowing history",
		Data:    rows,
		Total:   &total,
	})
}

func (s *BorrowController) GetBorrows(ctx *gin.Context) {

	var data dtos.GetAllBorrowedDto
	if err := ctx.ShouldBindQuery(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	borrows, count, err := s.borrowService.GetAllBorrowed(data)
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
		Message: "All borrowed books list",
		Data:    borrows,
		Total:   &total,
	})
}

func (s *BorrowController) BorrowDetails(ctx *gin.Context) {

	id := ctx.Param("id")
	borrow, err := s.borrowService.Details(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Show details of successful book borrowing",
		Data:    borrow,
	})
}

func (s *BorrowController) UpdateBorrow(ctx *gin.Context) {

	id := ctx.Param("id")
	var data dtos.BorrowDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	borrow, err := s.borrowService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	userId := ctx.GetString("userId")
	if borrow.UserId.Hex() != userId {
		ctx.JSON(http.StatusForbidden, models.ResponseJson{
			Success: false,
			Message: "Access is denied",
		})
		return
	}

	if httpStatus, err := s.borrowService.UpdateBorrow(data, *borrow); err != nil {
		ctx.JSON(httpStatus, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Book borrowing information has been successfully edited",
	})
}

func (s *BorrowController) DeleteBorrow(ctx *gin.Context) {

	id := ctx.Param("id")
	borrow, err := s.borrowService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if httpStatus, err := s.borrowService.DeleteBorrow(*borrow); err != nil {
		ctx.JSON(httpStatus, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Book borrowing data has been successfully deleted",
	})

}

func (s *BorrowController) BorrowBook(ctx *gin.Context) {

	var data dtos.BorrowDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	userId, err := primitive.ObjectIDFromHex(ctx.GetString("userId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	book, err := s.bookService.FindById(data.BookId)
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
			Message: "The book has been borrowed",
		})
		return
	}

	if _, err = s.borrowService.CreateBorrow(data, *book, userId); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Book borrowing record has been successfully saved",
		Data:    nil,
	})

}

func (s *BorrowController) ApproveBookBorrowing(ctx *gin.Context) {

	id := ctx.Param("id")
	borrow, err := s.borrowService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	httpStatus, err := s.borrowService.ApproveBookBorrowing(*borrow)
	if err != nil {
		ctx.JSON(httpStatus, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusBadRequest, models.ResponseJson{
		Success: true,
		Message: "Successfully approved book borrowing",
	})
}

func (s *BorrowController) CancelBorrowBook(ctx *gin.Context) {

	id := ctx.Param("id")
	borrow, err := s.borrowService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	userId := ctx.GetString("userId")
	if borrow.UserId.Hex() != userId {
		ctx.JSON(http.StatusForbidden, models.ResponseJson{
			Success: false,
			Message: "Access denied",
		})
		return
	}

	if httpStatus, err := s.borrowService.CancelBorrowBook(*borrow); err != nil {
		ctx.JSON(httpStatus, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Data:    true,
		Message: "The book borrowing status has been canceled.",
	})
}
