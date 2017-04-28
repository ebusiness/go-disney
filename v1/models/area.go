package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Area -
type Area struct {
	collectionName string        `collectionName:"areas"`
	ID             bson.ObjectId `bson:"_id"`
	Place          bson.ObjectId `bson:"place"`
	Language       `bson:",inline"`
}
