package models

import (
	// "time"
	"gopkg.in/mgo.v2/bson"
)

// Attraction -
type Attraction struct {
	collectionName string        `collectionName:"attractions"`
	ID             bson.ObjectId `json:"_id" bson:"_id"`
	StrID          string        `json:"str_id" bson:"str_id"`
	AreaName       string        `json:"area" bson:"area"`
	Category       string        `json:"category" bson:"category"`
	Park           Place         `json:"park" bson:"park"`
	ThumURL        string        `json:"thum_url_pc" bson:"thum_url_pc"`
	YoutubeURL     string        `json:"youtube_url,omitempty" bson:"youtube_url,omitempty"`
	SummaryTag     []summaryTag  `json:"summary_tags,omitempty" bson:"summary_tags,omitempty"`
	Maps           []string      `json:"maps" bson:"maps"`
	TagNames       []string      `json:"tags,omitempty" bson:"tags,omitempty"`
	Images         []string      `json:"images" bson:"main_visual_urls"`
	IsFastpass     bool          `json:"is_fastpass" bson:"is_fastpass"`
	IsLottery      bool          `json:"is_lottery" bson:"is_lottery"`
	IsMustBook     bool          `json:"is_must_book" bson:"is_must_book"`
	Name           string        `json:"name" bson:"name"`
	Note           string        `json:"note" bson:"note"`
	Introductions  string        `json:"introductions" bson:"introductions"`
	Summaries      []summary     `json:"summaries,omitempty" bson:"summaries,omitempty"`
	WaitTime       interface{}   `json:"waitTime,omitempty" bson:"waitTime,omitempty"`
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
