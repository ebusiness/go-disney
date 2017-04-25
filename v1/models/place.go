package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Place -
type Place struct {
	collectionName string        `collectionName:"areas"`
	ID             bson.ObjectId `bson:"_id"`
	Language
}
