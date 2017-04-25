package models

import (
	"gopkg.in/mgo.v2/bson"
)

// HotelType -
type HotelType struct {
	collectionName string        `collectionName:"areas"`
	ID             bson.ObjectId `bson:"_id"`
	Place          string        `bson:"place"`
	Language
}
