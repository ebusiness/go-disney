package models

// PlanCustomize -
type PlanCustomize struct {
	collectionName string `collectionName:"customize_plans"`
	PlanTemplate   `bson:",inline"`
	Lang           string `json:",none" bson:"lang,omitempty"`
}
