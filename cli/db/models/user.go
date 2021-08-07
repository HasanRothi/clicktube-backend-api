package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty" validate:"required"`
	Gmail      string             `bson:"gmail,omitempty" validate:"required,email"`
	Password   string             `bson:"password,omitempty" validate:"gte=5,lte=10"`
	University string             `bson:"university,omitempty" validate:"required"`
	CampusID   string             `bson:"campusId,omitempty" validate:"required"`
	Dept       string             `bson:"dept,omitempty" validate:"required"`
	Role       string             `bson:"role,omitempty" validate:"required"`
}
