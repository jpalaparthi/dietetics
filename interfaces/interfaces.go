package interfaces

import (
	"dietetics/models"
)

// Connection interface
type Connection interface {
	GetConnection(...string) (interface{}, error)
}

// RawFoodInterface interfaces
type RawFoodInterface interface {
	AddRawFood(rawfood models.RawFood) error
	AddNutrition(selector interface{}, nutrition models.Nutrition) error
	DeleteNutrition(where interface{}, selector interface{}) error
	GetRawFood(selector interface{}) (*models.RawFood, error)
	GetRawFoods(skip int32, limit int32, selector interface{}) ([]models.RawFood, error)
	UpdateRawFood(selector interface{}, data interface{}) error
	DeleteRawFood(selector interface{}) error
}
