package handler

import (
	"context"
	"dietetics/interfaces"
	"dietetics/models"
	proto "dietetics/proto/food/raw"
	"time"

	"flag"

	"github.com/golang/glog"
)

// RawFood type
type RawFood struct {
	IRawFood interfaces.RawFoodInterface
}

func init() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
}

// AddRawFood is to add rawfood in rawfood collection
func (r *RawFood) AddRawFood(ctx context.Context, req *proto.RawFood, rsp *proto.GeneralResponse) error {
	rf := models.RawFood{}
	ns := make([]models.Nutrition, 0)

	rf.Code = req.Code
	rf.Group = req.Group
	rf.Name = req.Name
	rf.SampleUnit = req.SampleUnit
	rf.SampleQuantity = req.SampleQuantity
	rf.Varietie = req.Varietie
	rf.Tags = req.Tags
	rf.Status = "active"
	rf.LastModified = time.Now().UTC().String()

	for _, v := range req.Nutritions {
		n := models.Nutrition{}
		n.Type = v.Type
		n.Name = v.Name
		n.Unit = v.Unit
		n.Group = v.Group
		n.Quantity = v.Quantity
		n.Tags = v.Tags
		ns = append(ns, n)
	}
	rf.Nutritions = ns
	if err := models.ValidateRawFood(rf); err != nil {
		rsp.Id = "AddRawFood"
		rsp.Code = 404
		rsp.Detail = err.Error()
		rsp.Status = "failure"
		glog.Error(err)

		return err
	}
	if err := r.IRawFood.AddRawFood(rf); err != nil {
		rsp.Id = "AddRawFood"
		rsp.Code = 404
		rsp.Detail = "error in inserting data"
		rsp.Status = "failure"
		glog.Error(err)
		return err
	}
	rsp.Id = "AddRawFood"
	rsp.Code = 201
	rsp.Detail = "raw food successfully inserted"
	rsp.Status = "success"
	glog.Info(rsp.Detail)
	return nil
}

func (r *RawFood) UpdateRawFood(ctx context.Context, in *proto.RawFood, out *proto.GeneralResponse) error {

	return nil
}

// GetRawFood is to fetch rawfood from the collection
func (r *RawFood) GetRawFood(ctx context.Context, req *proto.GetRawFoodRequest, rsp *proto.RawFood) error {
	query := make(map[string]interface{}, 1)
	for _, request := range req.RawFoods {
		query[request.Key] = request.Value
	}
	if rf, err := r.IRawFood.GetRawFood(query); err != nil {
		glog.Error(err)
		return err
	} else {
		rsp.Code = rf.Code
		rsp.Name = rf.Name
		rsp.Group = rf.Group
		rsp.Varietie = rf.Varietie
		rsp.SampleUnit = rf.SampleUnit
		rsp.SampleQuantity = rf.SampleQuantity
		rsp.Tags = rf.Tags
		rsp.Status = rf.Status
		rsp.LastModified = rf.LastModified
		for _, v := range rf.Nutritions {
			pn := &proto.Nutrition{}
			pn.Type = v.Type
			pn.Name = v.Name
			pn.Quantity = v.Quantity
			pn.Unit = v.Unit
			pn.Group = v.Group
			pn.Tags = v.Tags
			rsp.Nutritions = append(rsp.Nutritions, pn)
		}
	}
	return nil
}

// GetRawFoods is to get raw foods
func (r *RawFood) GetRawFoods(ctx context.Context, req *proto.GetRawFoodRequests, rsp *proto.RawFoodsResponse) error {
	query := make(map[string]interface{}, 1)
	for _, request := range req.RawFoods {
		query[request.Key] = request.Value
	}
	if rf, err := r.IRawFood.GetRawFoods(req.Skip, req.Limit, query); err != nil {
		glog.Error(err)
		return err
	} else {
		for _, val := range rf {
			rawfood := &proto.RawFood{}
			rawfood.Code = val.Code
			rawfood.Name = val.Name
			rawfood.Group = val.Group
			rawfood.Varietie = val.Varietie
			rawfood.SampleUnit = val.SampleUnit
			rawfood.SampleQuantity = val.SampleQuantity
			rawfood.Tags = val.Tags
			rawfood.Status = val.Status
			rawfood.LastModified = val.LastModified
			for _, v := range val.Nutritions {
				pn := &proto.Nutrition{}
				pn.Type = v.Type
				pn.Name = v.Name
				pn.Quantity = v.Quantity
				pn.Unit = v.Unit
				pn.Group = v.Group
				pn.Tags = v.Tags
				rawfood.Nutritions = append(rawfood.Nutritions, pn)
			}
			rsp.RawFoodRes = append(rsp.RawFoodRes, rawfood)
		}
	}
	return nil
}

// DeleteRawFood deletes an agent given by id
func (r *RawFood) DeleteRawFood(ctx context.Context, req *proto.PairRequest, rsp *proto.GeneralResponse) error {
	query := make(map[string]interface{}, 1)
	query[req.Key] = req.Value

	if err := r.IRawFood.DeleteRawFood(query); err != nil {
		rsp.Id = "RawService.DeleteRawFood"
		rsp.Code = 404
		rsp.Detail = "failure in deleting Rawfood"
		rsp.Status = "failure"
		glog.Error(err)
		return err
	}
	rsp.Id = "RawService.DeleteRawFood"
	rsp.Code = 200
	rsp.Detail = "Rawfood successfully deleted"
	rsp.Status = "success"
	glog.Info(rsp.Detail)
	return nil
}

// AddNutrition is to add nutrition information at a later state
func (r *RawFood) AddNutrition(ctx context.Context, req *proto.AddNutritionRequest, rsp *proto.GeneralResponse) error {
	query := make(map[string]interface{}, 1)
	for _, request := range req.Query {
		query[request.Key] = request.Value
	}
	nutrition := models.Nutrition{}
	nutrition.Name = req.Nutrition.Name
	nutrition.Type = req.Nutrition.Type
	nutrition.Aprox = req.Nutrition.Aprox
	nutrition.Quantity = req.Nutrition.Quantity
	nutrition.Unit = req.Nutrition.Unit
	nutrition.Group = req.Nutrition.Group
	nutrition.Tags = req.Nutrition.Tags

	if err := r.IRawFood.AddNutrition(query, nutrition); err != nil {
		rsp.Id = "RawService.AddNutition"
		rsp.Code = 404
		rsp.Detail = "failure in adding  Nutrition"
		rsp.Status = "failure"
		glog.Error(err)
		return err
	}

	rsp.Id = "RawService.AddNutition"
	rsp.Code = 201
	rsp.Detail = "Nutrition successfully added"
	rsp.Status = "success"
	glog.Info(rsp.Detail)
	return nil
}

// DeleteNutrition is to delete Nutrition information
func (r *RawFood) DeleteNutrition(ctx context.Context, req *proto.DeleteNutritionRequest, rsp *proto.GeneralResponse) error {
	rootSelector := make(map[string]interface{}, 1)
	for _, request := range req.RootSelector {
		rootSelector[request.Key] = request.Value
	}
	childSelector := make(map[string]interface{}, 1)
	for _, request := range req.ChildSelector {
		childSelector[request.Key] = request.Value
	}

	if err := r.IRawFood.DeleteNutrition(rootSelector, childSelector); err != nil {
		rsp.Id = "RawService.DeleteNutrition"
		rsp.Code = 404
		rsp.Detail = "failure in delete  Nutrition"
		rsp.Status = "failure"
		glog.Error(err)
		return err
	}
	rsp.Id = "RawService.DeleteNutrition"
	rsp.Code = 200
	rsp.Detail = "Nutrition successfully deleted"
	rsp.Status = "success"
	glog.Info(rsp.Detail)
	return nil
}
