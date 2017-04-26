package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Tag -
type VisitorTag struct {
	collectionName string        `collectionName:"visitor_tags"`
	ID             bson.ObjectId `json:"_id" bson:"_id"`
	Color          string        `json:"color" bson:"color"`
	Language       `json:",inline" bson:",inline"`
}
