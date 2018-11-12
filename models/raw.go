package models

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

var (
	ErrMandatoryField = "field is mandatory"
)

// RawFood model
type RawFood struct {
	ID             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Code           string        `json:"code" bson:"code"`
	Group          string        `json:"group" bson:"group"`
	Name           string        `json:"name" bson:"name"`
	Varietie       string        `json:"varietie" bson:"varietie"`
	SampleQuantity float32       `json:"sampleQuantity" bson:"sampleQuantity"`
	SampleUnit     string        `json:"sampleUnit" bson:"sampleUnit"`
	Nutritions     []Nutrition   `json:"nutritions" bson:"nutritions"`
	Tags           string        `json:"tags" bson:"tags"`
	Status         string        `json:"status" bson:"status"`
	LastModified   string        `json:"lastModified" bson:"lastModified"`
}

// Nutrition type
type Nutrition struct {
	Type     string  `json:"type" bson:"type"`
	Name     string  `json:"name" bson:"name"`
	Quantity float32 `json:"quantity" bson:"quantity"`
	Aprox    float32 `json:"aprox" bson:"aprox"`
	Unit     string  `json:"unit" bson:"unit"`
	Group    string  `json:"group" bson:"group"`
	Tags     string  `json:"tags" bson:"tags"`
}

// ValidateRawFood validates the RawFood type
func ValidateRawFood(rf RawFood) error {
	if rf.Code == "" {
		return fmt.Errorf("Code " + ErrMandatoryField)
	}
	if rf.Group == "" {
		return fmt.Errorf("Group " + ErrMandatoryField)
	}
	if rf.Name == "" {
		return fmt.Errorf("Name " + ErrMandatoryField)
	}
	if rf.Varietie == "" {
		return fmt.Errorf("Varietie " + ErrMandatoryField)
	}
	if rf.SampleQuantity <= 0 {
		return fmt.Errorf("SampleQuantity " + ErrMandatoryField)
	}
	if rf.SampleUnit == "" {
		return fmt.Errorf("SampleUnit " + ErrMandatoryField)
	}
	return nil
}

// ValidateNutrition validates Nutrition type
func ValidateNutrition(n Nutrition) error {
	return nil
}
