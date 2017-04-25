package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Service -
type Service struct {
	collectionName string        `collectionName:"areas"`
	ID             bson.ObjectId `bson:"_id"`
	Language
}
