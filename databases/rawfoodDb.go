package database

import (
	"dietetics/models"
	"errors"
	"fmt"
	"reflect"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//RawFoodDB is to create rawfood db entity
type RawFoodDB struct {
	Session interface{} //This should be type asserted based on the database usage
	DBName  string
}

const (
	ErrRawFoodExists  = "Rawfood code already exists"
	ErrWrongInputType = "Wrong input type"
)

//IsRawFoodExists is to check raw food exists or not
func (r *RawFoodDB) IsRawFoodExists(code string) bool {
	count, err := r.Session.(*mgo.Session).DB(r.DBName).C("rawfood").Find(bson.M{"code": code}).Count()
	if err != nil {
		if err.Error() == "not found" {
			return false
		} else {
			return false
		}
	}
	if count > 0 {
		return true
	}
	return false
}

// AddRawFood is to add raw food
func (r *RawFoodDB) AddRawFood(rawfood models.RawFood) error {
	if !r.IsRawFoodExists(rawfood.Code) {
		if err := r.Session.(*mgo.Session).DB(r.DBName).C("rawfood").Insert(rawfood); err != nil {
			return err
		}
	} else {
		return errors.New(ErrRawFoodExists)
	}
	return nil
}

// UpdateRawFood is to update raw food
func (r *RawFoodDB) UpdateRawFood(selector, data interface{}) error {
	if _, ok := selector.(map[string]interface{}); !ok {
		return errors.New(ErrWrongInputType)
	}
	if _, ok := data.(map[string]interface{}); !ok {
		return errors.New(ErrWrongInputType)
	}

	if err := r.Session.(*mgo.Session).DB(r.DBName).C("rawfood").Update(bson.M(selector.(map[string]interface{})), bson.M(data.(map[string]interface{}))); err != nil {
		return err
	}
	return nil
}

//GetRawFood is to fetch rawfood
func (r *RawFoodDB) GetRawFood(selector interface{}) (*models.RawFood, error) {
	if _, ok := selector.(map[string]interface{}); !ok {
		return nil, errors.New(ErrWrongInputType)
	}
	//var result map[string]interface{}
	var result *models.RawFood

	err := r.Session.(*mgo.Session).DB(r.DBName).C("rawfood").Find(bson.M(selector.(map[string]interface{}))).One(&result)
	if err != nil {
		return result, err
	}
	fmt.Println(result)
	return result, nil
}

//GetRawFoods is to fetch rawfood
func (r *RawFoodDB) GetRawFoods(skip, limit int32, selector interface{}) ([]models.RawFood, error) {
	if _, ok := selector.(map[string]interface{}); !ok {
		return nil, errors.New(ErrWrongInputType)
	}
	var result []models.RawFood
	err := r.Session.(*mgo.Session).DB(r.DBName).C("rawfood").Find(bson.M(selector.(map[string]interface{}))).Skip(int(skip)).Limit(int(limit)).Sort("-_id").All(&result)

	if err != nil {
		return result, err
	}
	return result, nil
}

// DeleteRawFood is to delete raw food
func (r *RawFoodDB) DeleteRawFood(selector interface{}) error {
	if _, ok := selector.(map[string]interface{}); !ok {
		return errors.New(ErrWrongInputType)
	}
	if err := r.Session.(*mgo.Session).DB(r.DBName).C("rawfood").Remove(selector); err != nil {
		return err
	}
	return nil
}

// AddNutrition add nutrition to the existing food
func (r *RawFoodDB) AddNutrition(selector interface{}, nutrition models.Nutrition) error {
	//query := bson.M{"code": q.(string)}
	update, err := ToMap(nutrition, "json")
	if err != nil {
		return err
	}
	update1 := bson.M{"$push": bson.M{"nutritions": update}}
	if err := r.Session.(*mgo.Session).DB(r.DBName).C("rawfood").Update(selector, update1); err != nil {
		return err
	}
	return nil
}

// DeleteNutrition delete nutrition to the existing food
func (r *RawFoodDB) DeleteNutrition(rootSelector, childSelector interface{}) error {
	update := bson.M{"$pull": bson.M{"nutritions": childSelector}}
	if err := r.Session.(*mgo.Session).DB(r.DBName).C("rawfood").Update(rootSelector, update); err != nil {
		return err
	}
	return nil
}

// ToMap converts a struct to a map using the struct's tags.
//
// ToMap uses tags on struct fields to decide which fields to add to the
// returned map.
func ToMap(in interface{}, tag string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" {
			out[tagv] = v.Field(i).Interface()
		}
	}
	return out, nil
}
