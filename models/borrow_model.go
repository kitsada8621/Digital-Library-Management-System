package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Borrow struct {
	ID                 primitive.ObjectID `json:"_id" bson:"_id"`
	UserId             primitive.ObjectID `json:"userId" bson:"userId"`
	User               *User              `json:"user,omitempty" bson:"user,omitempty"`
	BookId             primitive.ObjectID `json:"bookId" bson:"bookId"`
	Book               *Book              `json:"book,omitempty" bson:"book,omitempty"`
	BorrowDate         time.Time          `json:"borrowDate" bson:"borrowDate"`
	ScheduleReturnDate time.Time          `json:"scheduleReturnDate" bson:"scheduleReturnDate"`
	ReturnDate         time.Time          `json:"returnDate" bson:"returnDate"`
	BorrowStatus       int                `json:"borrowStatus" bson:"borrowStatus"` // 1: pending to approve, 2: approve, 3: pending to return, 4 approve return, 5: cancel
	CreatedAt          time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt" bson:"updatedAt"`
}
