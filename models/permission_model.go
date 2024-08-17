package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Permission struct {
	ID               primitive.ObjectID `json:"_id" bson:"_id"`
	PermissionName   string             `json:"permissionName" bson:"permissionName"`
	PermissionStatus int                `json:"permissionStatus" bson:"permissionStatus"` // 1: read only, 2: crud
}
