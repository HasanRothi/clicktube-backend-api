package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Link struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Link      string             `bson:"link,omitempty"`
	Views     int                `bson:"views"`
	ShortLink string             `bson:"shortLink"`
	Published bool               `bson:"published"`
	Date      time.Time          `bson:"date"`
	UrlKey    string             `bson:"urlKey"`
	Author    primitive.ObjectID `bson:"author"`
}
