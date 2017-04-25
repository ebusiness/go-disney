package models

import (
	"gopkg.in/mgo.v2/bson"
)

// TagType -
type TagType struct {
	collectionName string        `collectionName:"areas"`
	ID             bson.ObjectId `bson:"_id"`
	Language
}
