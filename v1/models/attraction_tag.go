package models

import (
	"gopkg.in/mgo.v2/bson"
)

// AttractionTag -
type AttractionTag struct {
	collectionName string        `collectionName:"attraction_tags"`
	ID             bson.ObjectId `json:"_id" bson:"_id"`
	Name           string        `json:"name,omitempty" bson:"name,omitempty"`
	PlanCreation   planCreation  `json:"plan_creation" bson:"plan_creation"`
	Attractions    []Attraction  `json:"attractions,omitempty" bson:"attractions,omitempty"`
}

type planCreation struct {
	Icon  string `json:"icon,omitempty" bson:"icon,omitempty"`
	Color string `json:"color,omitempty" bson:"color,omitempty"`
}
