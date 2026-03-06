package responses

import "gitlab.com/mims-api-service/models"

type RoadSurfaceResponds struct {
	Id          int    `json:"id"`
	UpdatedDate string `json:"update_date"`
	UpdatedBy   `json:"update_by"`
	// Status       string `json:"status"`
	// StatusCode   string `json:"status_code"`
	// RejectReason string `json:"reject_reason"`
	Permissions `json:"permissions"`
	Items       `json:"items"`
}

type Items struct {
	KmStart                  float64 `json:"km_start"`
	KmEnd                    float64 `json:"km_end"`
	SurfaceCrossSectionCode  int     `json:"surface_cross_section_code"`
	WidthSurface             float64 `json:"width_surface"`
	ThicknessSurface         float64 `json:"thickness_surface"`
	ThicknessSurfaceConcrete float64 `json:"thickness_surface_concrete"`
	WidthShoulderLeft        float64 `json:"width_shoulder_left"`
	SurfaceShoulderLeft      struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"surface_shoulder_left"`
	WidthShoulderRight   float64 `json:"width_shoulder_right"`
	SurfaceShoulderRight struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"surface_shoulder_right"`
	ThicknessBase         *float64                    `json:"thickness_base"`
	MaterialBase          *models.RefMaterialBase     `json:"material_base"`
	ThicknessSubbase      *float64                    `json:"thickness_subbase"`
	MaterialSubbase       *models.RefMaterialSubbase  `json:"material_subbase"`
	ThicknessSubgrade     *float64                    `json:"thickness_subgrade"`
	MaterialSubgrade      *models.RefMaterialSubgrade `json:"material_subgrade"`
	LaneCount             int                         `json:"lane_count"`
	Lanes                 []Lane                      `json:"lane"`
	ThicknessConcreteSlab *float64                    `json:"thickness_concrete_slab"`
}

type Lane struct {
	Surface struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		SurfaceGroup string `json:"surface_group"`
		ColorCode    string `json:"color_code" gorm:"column:color"`
	} `json:"surface"`
	Direction string `json:"direction"`
	LaneNo    int    `json:"lane_no"`
	GeomCl    string `json:"geom_cl"`
}

type ResPostRoadSurface struct {
	ID []int `json:"id"`
}
