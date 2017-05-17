package models

import (
// "time"
// "gopkg.in/mgo.v2/bson"
)

// Attraction -
type Attraction struct {
	collectionName string `collectionName:"attractions"`
	// ID             bson.ObjectId `json:"_id" bson:"_id"`
	// Park          Place        `json:"park" bson:"park"`
	StrID         string       `json:"str_id,omitempty" bson:"str_id,omitempty"`
	AreaName      string       `json:"area,omitempty" bson:"area,omitempty"`
	Category      string       `json:"category" bson:"category"`
	IsAvailable   bool         `json:"is_available" bson:"is_available"`
	YoutubeURL    string       `json:"youtube_url,omitempty" bson:"youtube_url,omitempty"`
	SummaryTag    []summaryTag `json:"summary_tags,omitempty" bson:"summary_tags,omitempty"`
	Maps          []string     `json:"maps,omitempty" bson:"maps,omitempty"`
	TagNames      []string     `json:"tags,omitempty" bson:"tags,omitempty"`
	Images        []string     `json:"images" bson:"main_visual_urls"`
	IsLottery     bool         `json:"is_lottery,omitempty" bson:"is_lottery,omitempty"`
	IsMustBook    bool         `json:"is_must_book,omitempty" bson:"is_must_book,omitempty"`
	Name          string       `json:"name,omitempty" bson:"name,omitempty"`
	Note          string       `json:"note,omitempty" bson:"note,omitempty"`
	Introductions string       `json:"introductions,omitempty" bson:"introductions,omitempty"`
	Summaries     []summary    `json:"summaries,omitempty" bson:"summaries,omitempty"`
	Realtime      interface{}  `json:"realtime,omitempty" bson:"realtime,omitempty"`
	// IsFastpass    bool         `json:"is_fastpass,omitempty" bson:"is_fastpass,omitempty"`
}

type summaryTag struct {
	TypeName string   `json:"type" bson:"type"`
	TagNames []string `json:"tags" bson:"tags"`
}

type summary struct {
	Body  string `json:"body,omitempty" bson:"body,omitempty"`
	Title string `json:"title,omitempty" bson:"title,omitempty"`
}

// type waitTime struct {
// 	WaitTime          string     `json:"waitTime" bson:"waitTime,omitempty"`
// 	Available         bool       `json:"available" bson:"available"`
// 	StatusInfo        string     `json:"statusInfo" bson:"statusInfo,omitempty"`
// 	FastpassAvailable bool       `json:"fastpassAvailable" bson:"fastpassAvailable"`
// 	FastpassInfo      string     `json:"fastpassInfo" bson:"fastpassInfo,omitempty"`
// 	UpdateTime        *time.Time `json:"updateTime" bson:"updateTime"`
//
// 	OperationStart *time.Time `json:"operation_start" bson:"operation_start,omitempty"`
// 	OperationEnd   *time.Time `json:"operation_end" bson:"operation_end,omitempty"`
// 	FastpassStart  *time.Time `json:"fastpass_start" bson:"fastpass_start,omitempty"`
// 	FastpassEnd    *time.Time `json:"fastpass_end" bson:"fastpass_end,omitempty"`
// 	CreateTime     *time.Time  `json:"createTime" bson:"createTime"`
// }
