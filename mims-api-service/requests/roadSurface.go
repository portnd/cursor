package requests

type RoadSurface struct {
	IDs                       []int      `json:"id"`
	KmStart                   []float64  `json:"km_start"`
	KmEnd                     []float64  `json:"km_end"`
	SurfaceCrosssectionCode   []int      `json:"surface_crosssection_code"`
	WidthSurface              []float64  `json:"width_surface"`
	ThicknessSurface          []float64  `json:"thickness_surface"`
	ThicknessSurfaceConcrete  []*float64 `json:"thickness_surface_concrete"`
	WidthShoulderLeft         []float64  `json:"width_shoulder_left"`
	RefSurfaceIdShoulderLeft  []int      `json:"ref_surface_shoulder_id_left"`
	WidthShoulderRight        []float64  `json:"width_shoulder_right"`
	RefSurfaceIdShoulderRight []int      `json:"ref_surface_shoulder_id_right"`
	ThicknessBase             []*float64 `json:"thickness_base"`
	MaterialBase              []*int     `json:"material_base"`
	ThicknessSubbase          []*float64 `json:"thickness_subbase"`
	MaterialSubbase           []*int     `json:"material_subbase"`
	ThicknessSubgrade         []*float64 `json:"thickness_subgrade"`
	MaterialSubgrade          []*int     `json:"material_subgrade"`
	LaneSurface               [][]int    `json:"lane_surface"`
	RoadId                    int        `json:"road_id"`
	ThicknessConcreteSlab     []*float64 `json:"thickness_concrete_slab"`
	RoadSurfaceIDParents      []int      `json:"road_surface_id"`

	// Status      string
	// CreatedBy   int       `json:"created_by"`
	// CreatedDate time.Time `json:"created_date"`
	// UpdatedBy   int        `json:"updated_by"`
	// UpdatedDate time.Time `json:"updated_date"`
	// Revision    int
	// IDParent    int
	// Year        int

}
