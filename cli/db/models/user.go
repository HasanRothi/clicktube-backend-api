package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	Gmail      string             `bson:"gmail,omitempty"`
	Password   string             `bson:"password,omitempty"`
	University string             `bson:"university,omitempty"`
	CampusID   string             `bson:"campusId,omitempty"`
	Dept       string             `bson:"dept,omitempty"`
}
