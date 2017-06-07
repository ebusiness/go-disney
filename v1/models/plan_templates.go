package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// PlanTemplate -
type PlanTemplate struct {
	collectionName string        `collectionName:"plan_templates"`
	ID             bson.ObjectId `json:"_id" bson:"_id"`
	Name           string        `json:"name" bson:"name"`
	Introduction   string        `json:"introduction" bson:"introduction"`
	Start          *time.Time    `json:"start,omitempty" bson:"start,omitempty"`
	Route          []PlanRoute   `json:"route" bson:"route"`
}

// PlanRoute -
type PlanRoute struct {
	StrID          string     `json:"str_id" bson:"str_id"`
	TimeCost       float64    `json:"timeCost" bson:"timeCost"`
	DistanceToNext float64    `json:"distanceToNext" bson:"distanceToNext"`
	WalktimeToNext float64    `json:"walktimeToNext"`
	WaitTime       float64    `json:"waitTime" bson:"waitTime"`
	Attraction     Attraction `json:"attraction" bson:"attraction"`
	Schedule       Schedule   `json:"schedule" bson:"schedule"`
}
