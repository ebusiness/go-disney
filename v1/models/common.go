package models

// Language -
type Language struct {
	JA string `json:"ja,omitempty" bson:"ja,omitempty"`
	EN string `json:"en,omitempty" bson:"en,omitempty"`
	CN string `json:"cn,omitempty" bson:"cn,omitempty"`
	TW string `json:"tw,omitempty" bson:"tw,omitempty"`
}
