package models

import (
	"gopkg.in/mgo.v2/bson"
)

// HotelType -
type HotelType struct {
	collectionName string        `collectionName:"hoteltypes"`
	ID             bson.ObjectId `bson:"_id"`
	Language       `bson:",inline"`
}
