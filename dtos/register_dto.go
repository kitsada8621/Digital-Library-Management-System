package dtos

type RegisterDto struct {
	Name            string `json:"name" bson:"name" binding:"required,max=255"`
	Phone           string `json:"phone" bson:"phone" binding:"required,max=50"`
	Email           string `json:"email" bson:"email" binding:"required,email,max=255"`
	Username        string `json:"username" bson:"username" binding:"required,max=255"`
	Password        string `json:"password" bson:"password" binding:"required,max=255"`
	ConfirmPassword string `json:"confirmPassword" bson:"confirmPassword" binding:"required,eqfield=Password,max=255"`
}
