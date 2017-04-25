package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Resort -
type Resort struct {
	collectionName string        `collectionName:"areas"`
	ID             bson.ObjectId `bson:"_id"`
	Language
}
