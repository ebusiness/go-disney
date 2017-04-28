package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Place -
type Place struct {
	collectionName string        `collectionName:"places"`
	ID             bson.ObjectId `bson:"_id"`
	Language       `bson:",inline"`
}
