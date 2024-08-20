package controllers

import (
	"dlms/models"
	"dlms/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReturnController struct {
	borrowService services.IBorrowService
}

func NewReturnController() ReturnController {
	return ReturnController{
		borrowService: services.BorrowService(),
	}
}

func (s *ReturnController) ReturnBook(ctx *gin.Context) {
	id := ctx.Param("id")
	borrow, err := s.borrowService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if httpStatus, err := s.borrowService.ReturnBook(*borrow); err != nil {
		ctx.JSON(httpStatus, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "The request to return the book was successfully submitted.",
	})
}

func (s *ReturnController) ApproveBookReturn(ctx *gin.Context) {
	id := ctx.Param("id")
	borrow, err := s.borrowService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if httpStatus, err := s.borrowService.ApproveBookReturn(*borrow); err != nil {
		ctx.JSON(httpStatus, models.ResponseJson{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Book return approval successful",
	})

}
