package dtos

type CategoryDto struct {
	CategoryName string `json:"categoryName" bson:"categoryName" binding:"required,max=255"`
}
