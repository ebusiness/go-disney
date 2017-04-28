package models

import (
	"gopkg.in/mgo.v2/bson"
)

// TagType -
type TagType struct {
	collectionName string        `collectionName:"tagtypes"`
	ID             bson.ObjectId `bson:"_id"`
	Language       `bson:",inline"`
}
