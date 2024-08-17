package dtos

type LoginDto struct {
	Username string `json:"username" bson:"username" binding:"required,max=255"`
	Password string `json:"password" bson:"password" binding:"required,max=255"`
}
