package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Service -
type Service struct {
	collectionName string        `collectionName:"services"`
	ID             bson.ObjectId `bson:"_id"`
	Language       `bson:",inline"`
}
