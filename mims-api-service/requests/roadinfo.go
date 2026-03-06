package requests

import "mime/multipart"

type Road struct {
	Name                      string                `form:"name"`
	RoadSectionID             int                   `form:"road_section_id"`
	RoadID                    int                   `form:"road_id"`
	RoadLevel                 int                   `form:"road_level"`
	KmStart                   float64               `form:"km_start"`
	KmEnd                     float64               `form:"km_end"`
	RampId                    string                `form:"ramp_id"`
	RoadColorCode             string                `form:"road_color_code"`
	RefRoadTypeID             int                   `form:"ref_road_type_id" validate:"nonzero"`
	RegisterDate              string                `form:"register_date"`
	Remark                    string                `form:"remark"`
	YearConstructionCompleted int                   `form:"year_construction_completed"`
	CenterLineShapeFile       *multipart.FileHeader `form:"center_line_shape_file" validate:"nonzero"`
	CenterLaneShapeFile       *multipart.FileHeader `form:"center_lane_shape_file" validate:"nonzero"`
}

type RoadUpdateInit struct {
	Name                      string                `form:"name"`
	KmStart                   float64               `form:"km_start" `
	KmEnd                     float64               `form:"km_end" `
	RampId                    string                `form:"ramp_id"`
	RoadColorCode             string                `form:"road_color_code"`
	RefRoadTypeID             int                   `form:"ref_road_type_id" validate:"nonzero"`
	RegisterDate              string                `form:"register_date" `
	YearConstructionCompleted int                   `form:"year_construction_completed"`
	Remark                    string                `form:"remark"`
	CenterLineShapeFileStatus string                `form:"center_line_shape_file_status"`
	CenterLaneShapeFileStatus string                `form:"center_lane_shape_file_status"`
	CenterLineShapeFile       *multipart.FileHeader `form:"center_line_shape_file"`
	CenterLaneShapeFile       *multipart.FileHeader `form:"center_lane_shape_file"`
}

type RoadUpdate struct {
	Name                      string                `form:"name"`
	KmStart                   float64               `form:"km_start" `
	KmEnd                     float64               `form:"km_end" `
	RampId                    string                `form:"ramp_id"`
	RoadColorCode             string                `form:"road_color_code"`
	RefRoadTypeID             int                   `form:"ref_road_type_id" validate:"nonzero"`
	RegisterDate              string                `form:"register_date" `
	Remark                    string                `form:"remark"`
	YearConstructionCompleted int                   `form:"year_construction_completed"`
	CenterLineShapeFileStatus string                `form:"center_line_shape_file_status"`
	CenterLaneShapeFileStatus string                `form:"center_lane_shape_file_status"`
	CenterLineShapeFile       *multipart.FileHeader `form:"center_line_shape_file"`
	CenterLaneShapeFile       *multipart.FileHeader `form:"center_lane_shape_file"`
}
