package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Attraction -
type Attraction struct {
	collectionName string        `collectionName:"attractions"`
	ID             bson.ObjectId `json:"_id" bson:"_id"`
	StrID          string        `json:"str_id" bson:"str_id"`
	AreaName       string        `json:"area" bson:"area"`
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
}

type summaryTag struct {
	TypeName string   `json:"type" bson:"type"`
	TagNames []string `json:"tags" bson:"tags"`
}

type summary struct {
	Body  string `json:"body,omitempty" bson:"body,omitempty"`
	Title string `json:"title,omitempty" bson:"title,omitempty"`
}
