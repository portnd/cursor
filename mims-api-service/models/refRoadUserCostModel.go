package models

type RefRoadUserCostAcc struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	NameEn string `json:"name_en"`
}

type RefRoadUserCostRuc struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	NameEn string `json:"name_en"`
}

func (rrt *RefRoadUserCostAcc) TableName() string {
	return "ref_road_user_cost_acc"
}

func (rrt *RefRoadUserCostRuc) TableName() string {
	return "ref_road_user_cost_ruc"
}
