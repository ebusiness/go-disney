package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Limited -
type Limited struct {
	collectionName string        `collectionName:"limiteds"`
	ID             bson.ObjectId `bson:"_id"`
	Language       `bson:",inline"`
}
