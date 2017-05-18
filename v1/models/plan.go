package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Plan -
type Plan struct {
	collectionName string `collectionName:"cache_plans"`
	PlanTemplate   `bson:",inline"`
	TemplateID     bson.ObjectId `json:"template_id" bson:"template_id"`
	Lang           string        `json:",none" bson:"lang,omitempty"`
}
