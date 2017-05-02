package models

import (
	// "gopkg.in/mgo.v2/bson"
	"time"
)

// Realtime -
type Realtime struct {
	collectionName string         `collectionName:"realtime"`
	Realtime       realTimeDetail `json:"realtime" bson:"realtime"`
	Attraction     `bson:",inline"`
}

type realTimeDetail struct {
	WaitTime          string     `json:"waitTime,omitempty" bson:"waitTime,omitempty"`
	Available         bool       `json:"available" bson:"available"`
	StatusInfo        string     `json:"statusInfo,omitempty" bson:"statusInfo,omitempty"`
	FastpassAvailable bool       `json:"fastpassAvailable" bson:"fastpassAvailable"`
	FastpassInfo      string     `json:"fastpassInfo,omitempty" bson:"fastpassInfo,omitempty"`
	UpdateTime        *time.Time `json:"updateTime" bson:"updateTime"`

	OperationStart *time.Time `json:"operation_start,omitempty" bson:"operation_start,omitempty"`
	OperationEnd   *time.Time `json:"operation_end,omitempty" bson:"operation_end,omitempty"`
	FastpassStart  *time.Time `json:"fastpass_start,omitempty" bson:"fastpass_start,omitempty"`
	FastpassEnd    *time.Time `json:"fastpass_end,omitempty" bson:"fastpass_end,omitempty"`
	CreateTime     time.Time  `json:"createTime" bson:"createTime"`
}

// Waittime -
type Waittime struct {
	collectionName string `collectionName:"realtime"`
	// ID             bson.ObjectId `json:"_id" bson:"_id"`
	WaitTime   string     `json:"waitTime,omitempty" bson:"waitTime,omitempty"`
	UpdateTime *time.Time `json:"updateTime" bson:"updateTime"`
	CreateTime time.Time  `json:"createTime" bson:"createTime"`
}
