package models

import (
	// "gopkg.in/mgo.v2/bson"
	"time"
)

// RealWaittime -
type RealWaittime struct {
	collectionName string `collectionName:"realtime"`
	// ID             bson.ObjectId `json:"_id" bson:"_id"`
	WaitTime string `json:"waitTime,omitempty" bson:"waitTime,omitempty"`
	// UpdateTime *time.Time `json:"updateTime" bson:"updateTime"`
	CreateTime time.Time `json:"createTime" bson:"createTime"`
	Available  bool      `json:"available" bson:"available"`
	// OperationStart *time.Time `json:"operation_start" bson:"operation_start"`
	OperationEnd *time.Time `json:"operation_end" bson:"operation_end"`
}

// PredictionWaittime -
type PredictionWaittime struct {
	collectionName string `collectionName:"waittimes"`
	// ID             bson.ObjectId `json:"_id" bson:"_id"`
	WaitTime string `json:"waitTime,omitempty" bson:"waitTime,omitempty"`
	// UpdateTime *time.Time `json:"updateTime" bson:"updateTime"`
	CreateTime time.Time `json:"createTime" bson:"createTime"`
}
