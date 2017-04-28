package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Tag -
type Tag struct {
	collectionName string        `collectionName:"tags"`
	ID             bson.ObjectId `bson:"_id"`
	Language       `bson:",inline"`
}
