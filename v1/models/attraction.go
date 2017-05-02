package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Attraction -
type Attraction struct {
	collectionName string        `collectionName:"attractions"`
	ID             bson.ObjectId `json:"_id" bson:"_id"`
	StrID          string        `json:"str_id" bson:"str_id"`
	Area           Area          `json:"area" bson:"area"`
	Park           Place         `json:"park" bson:"park"`
	ThumURL        string        `json:"thum_url_pc" bson:"thum_url_pc"`
	YoutubeURL     string        `json:"youtube_url,omitempty" bson:"youtube_url,omitempty"`
	SummaryTag     []summaryTag  `json:"summary_tags,omitempty" bson:"summary_tags,omitempty"`
	Maps           []string      `json:"maps" bson:"maps"`
	Tags           []Tag         `json:"tags,omitempty" bson:"tags,omitempty"`
	Images         []string      `json:"images" bson:"main_visual_urls"`
	IsFastpass     bool          `json:"is_fastpass" bson:"is_fastpass"`
	Name           Language      `json:"name" bson:"name"`
	Note           Language      `json:"note" bson:"note"`
	Introductions  Language      `json:"introductions" bson:"introductions"`
	Summaries      []summary     `json:"summaries,omitempty" bson:"summaries,omitempty"`
}

type summaryTag struct {
	Type TagType `json:"type" bson:"type"`
	Tags []Tag   `json:"tags" bson:"tags"`
}

type summary struct {
	Body  Language `json:"body" bson:"body"`
	Title Language `json:"title" bson:"title"`
}
