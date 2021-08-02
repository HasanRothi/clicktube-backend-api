package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Link struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Link      string             `bson:"link,omitempty"`
	Views     int                `bson:"views"`
	shortLink string             `bson:"shortLink"`
	Date      time.Time          `bson:"date"`
}
