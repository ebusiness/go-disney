package models

// AttractionRelations -
type AttractionRelations struct {
	collectionName string   `collectionName:"attraction_relations"`
	Distance       *float64 `json:"distance" bson:"distance"`
	From           string   `json:"from" bson:"from"`
	To             string   `json:"to" bson:"to"`
}
