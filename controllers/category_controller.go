package controllers

import (
	"dlms/dtos"
	"dlms/models"
	"dlms/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryController struct {
	categoryService   services.ICategoryService
	validationService services.IValidationService
}

func NewCategoryController() CategoryController {
	return CategoryController{
		categoryService:   services.CategoryService(),
		validationService: services.ValidationService(),
	}
}

func (s *CategoryController) GetCategories(ctx *gin.Context) {

	categories, err := s.categoryService.GetCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Category data display successful",
		Data:    categories,
	})

}

func (s *CategoryController) CreateCategory(ctx *gin.Context) {

	var data dtos.CategoryDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	if s.categoryService.CheckDuplicateByName(data.CategoryName) {
		ctx.JSON(http.StatusBadRequest, models.ResponseJson{
			Success: false,
			Message: "This category already has a name",
			Data:    nil,
		})
		return
	}

	if _, err := s.categoryService.CreateCategory(data); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Category created successfully.",
		Data:    nil,
	})

}

func (s *CategoryController) UpdateCategory(ctx *gin.Context) {

	var data dtos.CategoryDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, models.ResponseJson{
			Success: false,
			Message: "Unprocessed Entities",
			Data:    s.validationService.Validate(err),
		})
		return
	}

	id := ctx.Param("id")
	category, err := s.categoryService.FindCategoryById(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, models.ResponseJson{
				Success: false,
				Message: "Category not found",
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

	if c, _ := s.categoryService.FindCategoryByName(data.CategoryName); c != nil {
		if c.ID.Hex() != category.ID.Hex() {
			ctx.JSON(http.StatusBadRequest, models.ResponseJson{
				Success: false,
				Message: "This category already has a name",
				Data:    nil,
			})
			return
		}
	}

	if _, err := s.categoryService.UpdateCategory(data, category.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Category updated successfully.",
		Data:    nil,
	})

}

func (s *CategoryController) DeleteCategory(ctx *gin.Context) {

	id := ctx.Param("id")
	category, err := s.categoryService.FindCategoryById(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, models.ResponseJson{
				Success: false,
				Message: "Category not found",
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

	if err := s.categoryService.DeleteCategory(category.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ResponseJson{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseJson{
		Success: true,
		Message: "Category deleted successfully.",
		Data:    nil,
	})

}
