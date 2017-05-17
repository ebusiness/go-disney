package models

import (
	// "gopkg.in/mgo.v2/bson"
	"time"
)

// RealWaittime -
type RealWaittime struct {
	collectionName string     `collectionName:"realtime"`
	WaitTime       float64    `json:"waitTime,omitempty" bson:"waitTime,omitempty"`
	CreateTime     time.Time  `json:"createTime" bson:"createTime"`
	Available      bool       `json:"available" bson:"available"`
	OperationEnd   *time.Time `json:"operation_end" bson:"operation_end"`
	// ID             bson.ObjectId `json:"_id" bson:"_id"`
	// UpdateTime *time.Time `json:"updateTime" bson:"updateTime"`
	// OperationStart *time.Time `json:"operation_start" bson:"operation_start"`
}

// PredictionWaittime -
type PredictionWaittime struct {
	collectionName string    `collectionName:"waittimes"`
	WaitTime       float64   `json:"waitTime,omitempty" bson:"waitTime,omitempty"`
	CreateTime     time.Time `json:"createTime" bson:"createTime"`
	// ID             bson.ObjectId `json:"_id" bson:"_id"`
	// UpdateTime *time.Time `json:"updateTime" bson:"updateTime"`
}
