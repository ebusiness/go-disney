package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Resort -
type Resort struct {
	collectionName string        `collectionName:"resorts"`
	ID             bson.ObjectId `bson:"_id"`
	Language       `bson:",inline"`
}
