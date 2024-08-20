package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	CreatedBy       primitive.ObjectID `json:"createdBy" bson:"createdBy"`
	CategoryId      primitive.ObjectID `json:"categoryId" bson:"categoryId"`
	Category        *Category          `json:"category,omitempty" bson:"category,omitempty"`
	Author          string             `json:"author" bson:"author"`
	BookTitle       string             `json:"bookTitle" bson:"bookTitle"`
	BookDescription string             `json:"bookDesc" bson:"bookDesc"`
	BorrowCount     int                `json:"borrowCount" bson:"borrowCount"`
	BookStatus      int                `json:"bookStatus" bson:"bookStatus"` // 1:Unavailable status  0: Vacant status
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
}
