package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Attraction -
type Attraction struct {
	collectionName string        `collectionName:"areas"`
	ID             bson.ObjectId `bson:"_id"`
}
