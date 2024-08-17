package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	RoleName    string             `json:"roleName" bson:"roleName"`
	Permissions []Permission       `json:"permissions" bson:"permissions"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}
