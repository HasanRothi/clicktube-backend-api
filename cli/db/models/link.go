package models

type Link struct {
	Link  string `bson:"link,omitempty"`
	Views string `bson:"views"`
	Short string `bson:"short"`
}
