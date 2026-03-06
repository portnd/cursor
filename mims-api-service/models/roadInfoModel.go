package models

import "time"

type RoadInfo struct {
	Id                        int       `json:"id"`
	RoadId                    int       `json:"road_id"`
	Year                      *int      `json:"year"`
	RefDirectionId            int       `json:"ref_direction_id"`
	Name                      string    `json:"name"`
	KmStart                   float32   `json:"km_start"`
	KmEnd                     float32   `json:"km_end"`
	TheGeom                   string    `json:"the_geom"`
	Revision                  int       `json:"revision"`
	Status                    string    `json:"status"`
	RampId                    string    `json:"ramp_id"`
	RoadColorCode             string    `json:"road_color_code"`
	CreatedBy                 int       `json:"created_by"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedBy                 int       `json:"updated_by"`
	UpdatedAt                 time.Time `json:"updated_at"`
	Remark                    string    `json:"remark"`
	RefRoadTypeID             int       `json:"ref_road_type_id"`
	CenterLaneShapeFilePath   string    `json:"center_lane_shape_file_path"`
	CenterLineShapeFilePath   string    `json:"center_line_shape_file_path"`
	YearConstructionCompleted int       `json:"year_construction_completed"`
}

type RoadInfoData struct {
	Id             int          `json:"id"`
	RoadId         int          `json:"road_id"`
	RefDirectionId int          `json:"ref_direction_id"`
	Name           string       `json:"name"`
	RoadColorCode  string       `json:"road_color_code"`
	Direction      RefDirection `json:"direction" gorm:"ForeignKey:RefDirectionId;AssociationForeignKey:Id"`
}

type RoadInfoAddData struct {
	RoadInfo
	OriginToDestination string       `json:"origin_to_destination"`
	RoadCode            string       `json:"road_code"`
	ResponsibleCode     string       `json:"responsible_code"`
	KmRange             string       `json:"km_range"`
	RefDirection        RefDirection `json:"direction" gorm:"ForeignKey:RefDirectionId;AssociationForeignKey:ID"`
	User                UserRes      `json:"user" gorm:"ForeignKey:CreatedBy;AssociationForeignKey:Id"`
	RefRoadType         RefRoadType  `json:"ref_road_type" gorm:"ForeignKey:RefRoadTypeID;AssociationForeignKey:Id"`
}

type RoadInfoDataDirection struct {
	Id             int          `json:"id"`
	RoadId         int          `json:"road_id"`
	RefDirectionId int          `json:"ref_direction_id"`
	Name           string       `json:"name"`
	RoadColorCode  string       `json:"road_color_code"`
	Direction      RefDirection `json:"direction" gorm:"ForeignKey:RefDirectionId;AssociationForeignKey:ID"`
}

func (rt *RoadInfoDataDirection) TableName() string {
	return "road_info"
}

type RoadInfoGeomData struct {
	RoadInfo
	LineString string         `json:"line_string"`
	RoadGeom   []RoadGeomData `json:"road_geom" gorm:"ForeignKey:RoadID;references:RoadId"`
}

type RoadInfoForDashboard struct {
	Id                      int       `json:"id"`
	RoadId                  int       `json:"road_id"`
	Year                    *int      `json:"year"`
	RefDirectionId          int       `json:"ref_direction_id"`
	Name                    string    `json:"name"`
	KmStart                 float32   `json:"km_start"`
	KmEnd                   float32   `json:"km_end"`
	TheGeom                 string    `json:"the_geom"`
	TheGeomJson             string    `json:"the_geom_json"`
	Revision                int       `json:"revision"`
	Status                  string    `json:"status"`
	RampId                  string    `json:"ramp_id"`
	RoadColorCode           string    `json:"road_color_code"`
	CreatedBy               int       `json:"created_by"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedBy               int       `json:"updated_by"`
	UpdatedAt               time.Time `json:"updated_at"`
	Remark                  string    `json:"remark"`
	RefRoadTypeID           int       `json:"ref_road_type_id"`
	CenterLaneShapeFilePath string    `json:"center_lane_shape_file_path"`
	CenterLineShapeFilePath string    `json:"center_line_shape_file_path"`
}

// TableName use to specific table
func (rt *RoadInfo) TableName() string {
	return "road_info"
}

func (rt *RoadInfoData) TableName() string {
	return "road_info"
}

func (rt *RoadInfoForDashboard) TableName() string {
	return "road_info"
}
