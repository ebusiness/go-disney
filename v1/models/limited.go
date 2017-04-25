package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Limited -
type Limited struct {
	collectionName string        `collectionName:"areas"`
	ID             bson.ObjectId `bson:"_id"`
	Language
}
