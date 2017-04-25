package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Tag -
type Tag struct {
	collectionName string        `collectionName:"areas"`
	ID             bson.ObjectId `bson:"_id"`
	Language
}
