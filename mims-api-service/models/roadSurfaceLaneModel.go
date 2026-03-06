package models

// Todo ...
type RoadSurfaceLane struct {
	Id                 int    `json:"id"`
	RoadSurfaceId      int    `json:"road_surface_id"`
	LaneNo             int    `json:"lane_no"`
	RoadId             int    `json:"road_id"`
	RefSurfaceId       int    `json:"ref_surface_id"`
	TheGeom            string `json:"the_geom"`
	RefSurfaceParamsID int    `json:"ref_surface_params_id"`
}

type RoadSurfaceLanePrePareData struct {
	RoadSurfaceLane
	RefSurface RefSurface `json:"ref_surface" gorm:"ForeignKey:RefSurfaceParamsID;AssociationForeignKey:ID"`
}

type RoadSurfaceLanePrePareDataLane struct {
	RoadSurfaceLane
	RefSurface RefSurfaceLane `json:"ref_surface" gorm:"ForeignKey:RefSurfaceId;AssociationForeignKey:ID"`
}

type RoadSurfaceLaneCountLane struct {
	Id            int `json:"id"`
	RoadSurfaceId int `json:"road_surface_id"`
	LaneNo        int `json:"lane_no"`
}

// TableName use to specific table
func (b *RoadSurfaceLane) TableName() string {
	return "road_surface_lane"
}

func (b *RoadSurfaceLaneCountLane) TableName() string {
	return "road_surface_lane"
}
