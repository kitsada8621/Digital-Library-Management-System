package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID  `json:"_id" bson:"_id"`
	RoleId    *primitive.ObjectID `json:"roleId" bson:"roleId"`
	Role      *Role               `json:"role,omitempty" bson:"role,omitempty"`
	Name      string              `json:"name" bson:"name"`
	Phone     string              `json:"phone" bson:"phone"`
	Email     string              `json:"email" bson:"email"`
	Username  string              `json:"username" bson:"username"`
	Password  string              `json:"password,omitempty" bson:"password,omitempty"`
	Token     *string             `json:"token,omitempty" bson:"token,omitempty"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAT time.Time           `json:"updatedAt" bson:"updatedAt"`
}
