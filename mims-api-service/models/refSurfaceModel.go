package models

type RefSurface struct {
	ID                       int      `json:"id"`
	Name                     string   `json:"name"`
	LayerCoefficient         float64  `json:"layer_coefficient"`
	Drainage                 float64  `json:"drainage"`
	Type                     string   `json:"type"`
	SurfaceGroup             string   `json:"surface_group"`
	A                        float64  `json:"a"`
	B                        float64  `json:"b"`
	C                        string   `json:"c"`
	CRT                      *float64 `json:"crt"`
	RRF                      *float64 `json:"rrf"`
	RavelingInitA0           *float64 `json:"raveling_init_a0"`
	RavelingInitA1           *float64 `json:"raveling_init_a1"`
	RavelingProgressA0       *float64 `json:"raveling_progress_a0"`
	RavelingProgressA1       *float64 `json:"raveling_progress_a1"`
	RavelingProgressA2       *float64 `json:"raveling_progress_a2"`
	AscInitHsoldEqual0A0     *float64 `json:"asc_init_hsold_equal_0_a0"`
	AscInitHsoldEqual0A1     *float64 `json:"asc_init_hsold_equal_0_a1"`
	AscInitHsoldEqual0A2     *float64 `json:"asc_init_hsold_equal_0_a2"`
	AscInitHsoldEqual0A3     *float64 `json:"asc_init_hsold_equal_0_a3"`
	AscInitHsoldEqual0A4     *float64 `json:"asc_init_hsold_equal_0_a4"`
	AscInitHsoldOver0A0      *float64 `json:"asc_init_hsold_over_0_a0"`
	AscInitHsoldOver0A1      *float64 `json:"asc_init_hsold_over_0_a1"`
	AscInitHsoldOver0A2      *float64 `json:"asc_init_hsold_over_0_a2"`
	AscInitHsoldOver0A3      *float64 `json:"asc_init_hsold_over_0_a3"`
	AscInitHsoldOver0A4      *float64 `json:"asc_init_hsold_over_0_a4"`
	AscProgressHsoldEqual0A0 *float64 `json:"asc_progress_hsold_equal_0_a0"`
	AscProgressHsoldEqual0A1 *float64 `json:"asc_progress_hsold_equal_0_a1"`
	AscProgressHsoldOver0A0  *float64 `json:"asc_progress_hsold_over_0_a0"`
	AscProgressHsoldOver0A1  *float64 `json:"asc_progress_hsold_over_0_a1"`
	WscInitHsoldEqual0A0     *float64 `json:"wsc_init_hsold_equal_0_a0"`
	WscInitHsoldEqual0A1     *float64 `json:"wsc_init_hsold_equal_0_a1"`
	WscInitHsoldEqual0A2     *float64 `json:"wsc_init_hsold_equal_0_a2"`
	WscInitHsoldOver0A0      *float64 `json:"wsc_init_hsold_over_0_a0"`
	WscInitHsoldOver0A1      *float64 `json:"wsc_init_hsold_over_0_a1"`
	WscInitHsoldOver0A2      *float64 `json:"wsc_init_hsold_over_0_a2"`
	WscProgressHsoldEqual0A0 *float64 `json:"wsc_progress_hsold_equal_0_a0"`
	WscProgressHsoldEqual0A1 *float64 `json:"wsc_progress_hsold_equal_0_a1"`
	WscProgressHsoldOver0A0  *float64 `json:"wsc_progress_hsold_over_0_a0"`
	WscProgressHsoldOver0A1  *float64 `json:"wsc_progress_hsold_over_0_a1"`
	RpdA0                    *float64 `json:"rpd_a0"`
	RpdA1                    *float64 `json:"rpd_a1"`
	RpdA2                    *float64 `json:"rpd_a2"`
	CanDelete                bool     `json:"can_delete"`
	Color                    string   `json:"color"`
}
type RefSurfaceNew struct {
	ID                       int     `json:"id"`
	Name                     string  `json:"name"`
	LayerCoefficient         float64 `json:"layer_coefficient"`
	Drainage                 float64 `json:"drainage"`
	Type                     string  `json:"type"`
	SurfaceGroup             string  `json:"surface_group"`
	A                        float64 `json:"a"`
	B                        float64 `json:"b"`
	C                        string  `json:"c"`
	CRT                      float64 `json:"crt"`
	RRF                      float64 `json:"rrf"`
	RavelingInitA0           float64 `json:"raveling_init_a0"`
	RavelingInitA1           float64 `json:"raveling_init_a1"`
	RavelingProgressA0       float64 `json:"raveling_progress_a0"`
	RavelingProgressA1       float64 `json:"raveling_progress_a1"`
	RavelingProgressA2       float64 `json:"raveling_progress_a2"`
	AscInitHsoldEqual0A0     float64 `json:"asc_init_hsold_equal_0_a0" gorm:"column:asc_init_hsold_equal_0_a0"`
	AscInitHsoldEqual0A1     float64 `json:"asc_init_hsold_equal_0_a1" gorm:"column:asc_init_hsold_equal_0_a1"`
	AscInitHsoldEqual0A2     float64 `json:"asc_init_hsold_equal_0_a2" gorm:"column:asc_init_hsold_equal_0_a2"`
	AscInitHsoldEqual0A3     float64 `json:"asc_init_hsold_equal_0_a3" gorm:"column:asc_init_hsold_equal_0_a3"`
	AscInitHsoldEqual0A4     float64 `json:"asc_init_hsold_equal_0_a4" gorm:"column:asc_init_hsold_equal_0_a4"`
	AscInitHsoldOver0A0      float64 `json:"asc_init_hsold_over_0_a0" gorm:"column:asc_init_hsold_over_0_a0"`
	AscInitHsoldOver0A1      float64 `json:"asc_init_hsold_over_0_a1" gorm:"column:asc_init_hsold_over_0_a1"`
	AscInitHsoldOver0A2      float64 `json:"asc_init_hsold_over_0_a2" gorm:"column:asc_init_hsold_over_0_a2"`
	AscInitHsoldOver0A3      float64 `json:"asc_init_hsold_over_0_a3" gorm:"column:asc_init_hsold_over_0_a3"`
	AscInitHsoldOver0A4      float64 `json:"asc_init_hsold_over_0_a4" gorm:"column:asc_init_hsold_over_0_a4"`
	AscProgressHsoldEqual0A0 float64 `json:"asc_progress_hsold_equal_0_a0" gorm:"column:asc_progress_hsold_equal_0_a0"`
	AscProgressHsoldEqual0A1 float64 `json:"asc_progress_hsold_equal_0_a1" gorm:"column:asc_progress_hsold_equal_0_a1"`
	AscProgressHsoldOver0A0  float64 `json:"asc_progress_hsold_over_0_a0" gorm:"column:asc_progress_hsold_over_0_a0"`
	AscProgressHsoldOver0A1  float64 `json:"asc_progress_hsold_over_0_a1" gorm:"column:asc_progress_hsold_over_0_a1"`
	WscInitHsoldEqual0A0     float64 `json:"wsc_init_hsold_equal_0_a0" gorm:"column:wsc_init_hsold_equal_0_a0"`
	WscInitHsoldEqual0A1     float64 `json:"wsc_init_hsold_equal_0_a1" gorm:"column:wsc_init_hsold_equal_0_a1"`
	WscInitHsoldEqual0A2     float64 `json:"wsc_init_hsold_equal_0_a2" gorm:"column:wsc_init_hsold_equal_0_a2"`
	WscInitHsoldOver0A0      float64 `json:"wsc_init_hsold_over_0_a0" gorm:"column:wsc_init_hsold_over_0_a0"`
	WscInitHsoldOver0A1      float64 `json:"wsc_init_hsold_over_0_a1" gorm:"column:wsc_init_hsold_over_0_a1"`
	WscInitHsoldOver0A2      float64 `json:"wsc_init_hsold_over_0_a2" gorm:"column:wsc_init_hsold_over_0_a2"`
	WscProgressHsoldEqual0A0 float64 `json:"wsc_progress_hsold_equal_0_a0" gorm:"column:wsc_progress_hsold_equal_0_a0"`
	WscProgressHsoldEqual0A1 float64 `json:"wsc_progress_hsold_equal_0_a1" gorm:"column:wsc_progress_hsold_equal_0_a1"`
	WscProgressHsoldOver0A0  float64 `json:"wsc_progress_hsold_over_0_a0" gorm:"column:wsc_progress_hsold_over_0_a0"`
	WscProgressHsoldOver0A1  float64 `json:"wsc_progress_hsold_over_0_a1" gorm:"column:wsc_progress_hsold_over_0_a1"`
	RpdA0                    float64 `json:"rpd_a0"`
	RpdA1                    float64 `json:"rpd_a1"`
	RpdA2                    float64 `json:"rpd_a2"`
}

type RefSurfaceType struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	SurfaceGroup string `json:"surface_group"`
}

type NewRefSurface struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	LayerCoefficient float64  `json:"layer_coefficient"`
	Drainage         float64  `json:"drainage"`
	Type             string   `json:"type"`
	SurfaceGroup     string   `json:"surface_group"`
	A                float64  `json:"a"`
	B                float64  `json:"b"`
	C                string   `json:"c"`
	CRT              *float64 `json:"crt"`
	RRF              *float64 `json:"rrf"`
	///////////////////////////////////////
	// RavelingInitA0           *float64 `json:"raveling_init_a0"`
	// RavelingInitA1           *float64 `json:"raveling_init_a1"`
	// RavelingProgressA0       *float64 `json:"raveling_progress_a0"`
	// RavelingProgressA1       *float64 `json:"raveling_progress_a1"`
	// RavelingProgressA2       *float64 `json:"raveling_progress_a2"`
	// AscInitHsoldEqual0A0     *float64 `json:"asc_init_hsold_equal_0_a0"`
	// AscInitHsoldEqual0A1     *float64 `json:"asc_init_hsold_equal_0_a1"`
	// AscInitHsoldEqual0A2     *float64 `json:"asc_init_hsold_equal_0_a2"`
	// AscInitHsoldEqual0A3     *float64 `json:"asc_init_hsold_equal_0_a3"`
	// AscInitHsoldEqual0A4     *float64 `json:"asc_init_hsold_equal_0_a4"`
	// AscInitHsoldOver0A0      *float64 `json:"asc_init_hsold_over_0_a0"`
	// AscInitHsoldOver0A1      *float64 `json:"asc_init_hsold_over_0_a1"`
	// AscInitHsoldOver0A2      *float64 `json:"asc_init_hsold_over_0_a2"`
	// AscInitHsoldOver0A3      *float64 `json:"asc_init_hsold_over_0_a3"`
	// AscInitHsoldOver0A4      *float64 `json:"asc_init_hsold_over_0_a4"`
	// AscProgressHsoldEqual0A0 *float64 `json:"asc_progress_hsold_equal_0_a0"`
	// AscProgressHsoldEqual0A1 *float64 `json:"asc_progress_hsold_equal_0_a1"`
	// AscProgressHsoldOver0A0  *float64 `json:"asc_progress_hsold_over_0_a0"`
	// AscProgressHsoldOver0A1  *float64 `json:"asc_progress_hsold_over_0_a1"`
	// WscInitHsoldEqual0A0     *float64 `json:"wsc_init_hsold_equal_0_a0"`
	// WscInitHsoldEqual0A1     *float64 `json:"wsc_init_hsold_equal_0_a1"`
	// WscInitHsoldEqual0A2     *float64 `json:"wsc_init_hsold_equal_0_a2"`
	// WscInitHsoldOver0A0      *float64 `json:"wsc_init_hsold_over_0_a0"`
	// WscInitHsoldOver0A1      *float64 `json:"wsc_init_hsold_over_0_a1"`
	// WscInitHsoldOver0A2      *float64 `json:"wsc_init_hsold_over_0_a2"`
	// WscProgressHsoldEqual0A0 *float64 `json:"wsc_progress_hsold_equal_0_a0"`
	// WscProgressHsoldEqual0A1 *float64 `json:"wsc_progress_hsold_equal_0_a1"`
	// WscProgressHsoldOver0A0  *float64 `json:"wsc_progress_hsold_over_0_a0"`
	// WscProgressHsoldOver0A1  *float64 `json:"wsc_progress_hsold_over_0_a1"`
	// RpdA0                    *float64 `json:"rpd_a0"`
	// RpdA1                    *float64 `json:"rpd_a1"`
	// RpdA2                    *float64 `json:"rpd_a2"`
	CanDelete bool `json:"can_delete"`
}

type SurfaceShoulder struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RefSurfaceLane struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	SurfaceGroup string `json:"surface_group"`
	ColorCode    string `json:"color_code" gorm:"column:color"`
}

func (rs *RefSurface) TableName() string {
	return "ref_surface"
}
func (rs *RefSurfaceNew) TableName() string {
	return "ref_surface"
}

func (rs *RefSurfaceType) TableName() string {
	return "ref_surface_type"
}

func (rs *SurfaceShoulder) TableName() string {
	return "ref_surface"
}

func (rs *RefSurfaceLane) TableName() string {
	return "ref_surface"
}
