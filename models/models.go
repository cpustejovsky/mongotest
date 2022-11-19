package models

type Animal struct {
	ID      string `json:"id,omitempty" bson:"_id"`
	Species string `json:"species"`
}

type AnimalForm struct {
	Species *string `json:"species"`
}
