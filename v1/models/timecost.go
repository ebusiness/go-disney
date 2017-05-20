package models

// TimeCost -
type TimeCost struct {
	collectionName string   `collectionName:"timecosts"`
	StrID          string   `json:"str_id" bson:"str_id"`
	Cost           *float64 `json:"timeCost" bson:"timeCost"`
}
