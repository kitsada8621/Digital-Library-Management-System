package controllers

import (
	"dlms/dtos"
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			`success`: false,
			`message`: err.Error(),
			`data`:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		`success`: true,
		`message`: "Category data display successful",
		`data`:    categories,
	})

}

func (s *CategoryController) CreateCategory(ctx *gin.Context) {

	var data dtos.CategoryDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			`success`: false,
			`message`: "Unprocessed Entities",
			`data`:    s.validationService.Validate(err),
		})
		return
	}

	if s.categoryService.CheckDuplicateByName(data.CategoryName) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			`success`: false,
			`message`: "This category already has a name",
			`data`:    nil,
		})
		return
	}

	if _, err := s.categoryService.CreateCategory(data); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			`success`: false,
			`message`: err.Error(),
			`data`:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		`success`: true,
		`message`: "Category created successfully.",
		`data`:    nil,
	})

}

func (s *CategoryController) UpdateCategory(ctx *gin.Context) {

	var data dtos.CategoryDto
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			`success`: false,
			`message`: "Unprocessed Entities",
			`data`:    s.validationService.Validate(err),
		})
		return
	}

	id := ctx.Param("id")
	category, err := s.categoryService.FindCategoryById(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{
				`success`: false,
				`message`: "Category not found",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				`success`: false,
				`message`: err.Error(),
			})
			return
		}
	}

	if c, _ := s.categoryService.FindCategoryByName(data.CategoryName); c != nil {
		if c.ID.Hex() != category.ID.Hex() {
			ctx.JSON(http.StatusBadRequest, gin.H{
				`success`: false,
				`message`: "This category already has a name",
				`data`:    nil,
			})
			return
		}
	}

	if _, err := s.categoryService.UpdateCategory(data, category.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			`success`: false,
			`message`: err.Error(),
			`data`:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		`success`: true,
		`message`: "Category updated successfully.",
		`data`:    nil,
	})

}

func (s *CategoryController) DeleteCategory(ctx *gin.Context) {

	id := ctx.Param("id")
	category, err := s.categoryService.FindCategoryById(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{
				`success`: false,
				`message`: "Category not found",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				`success`: false,
				`message`: err.Error(),
			})
			return
		}
	}

	if err := s.categoryService.DeleteCategory(category.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			`success`: false,
			`message`: err.Error(),
			`data`:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		`success`: true,
		`message`: "Category deleted successfully.",
		`data`:    nil,
	})

}
