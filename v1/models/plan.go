package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Plan -
type Plan struct {
	collectionName string `collectionName:"cache_plans"`
	PlanTemplate   `bson:",inline"`
	TemplateID     *bson.ObjectId `json:",none" bson:"template_id,omitempty"`
	Lang           string         `json:",none" bson:"lang,omitempty"`
}
