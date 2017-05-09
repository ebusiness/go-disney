package models

import (
	"gopkg.in/mgo.v2/bson"
)

// VisitorTag -
type VisitorTag struct {
	collectionName string        `collectionName:"visitor_tags"`
	ID             bson.ObjectId `json:"_id" bson:"_id"`
	Color          string        `json:"color,omitempty" bson:"color,omitempty"`
	Name           string        `json:"name,omitempty" bson:"name,omitempty"`
}
