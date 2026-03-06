package models

import "time"

type SettingDeteriorationParams struct {
	Id        int       `json:"id"`
	Params    string    `json:"params"`
	IsLatest  bool      `json:"is_latest"`
	IsDeleted bool      `json:"is_deleted"`
	UpdatedBy int       `json:"updated_by"`
	CreatedBy int       `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type SettingDeterioration struct {
	Id             int       `json:"id"`
	RoadGroupId    int       `json:"road_group_id"`
	ParamsAsphalt  string    `json:"params_Asphalt"`
	ParamsConcrete string    `json:"params_concrete"`
	IsLatest       bool      `json:"is_latest"`
	IsDeleted      bool      `json:"is_deleted"`
	UpdatedBy      int       `json:"updated_by"`
	CreatedBy      int       `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedAt      time.Time `json:"created_at"`
}

type DeteriorationParams struct {
	Asphalt  []DeteriorationAsphalt  `json:"asphalt"`
	Concrete []DeteriorationConcrete `json:"concrete"`
}

type DeteriorationAsphalt struct {
	RoadGroupId int     `json:"road_group_id"`
	Tlf         float64 `json:"tlf"`
	Cdb         float64 `json:"cdb"`
	Cds         float64 `json:"cds"`
	Comp        float64 `json:"comp"`
	Kvi         float64 `json:"kvi"`
	Kvp         float64 `json:"kvp"`
	Kpi         float64 `json:"kpi"`
	Kpp         float64 `json:"kpp"`
	Krid        float64 `json:"krid"`
	Krst        float64 `json:"krst"`
	Krpd        float64 `json:"krpd"`
	Kgm         float64 `json:"kgm"`
	Kgp         float64 `json:"kgp"`
	Kcia        float64 `json:"kcia"`
	Cmod        float64 `json:"cmod"`
	Kciw        float64 `json:"kciw"`
	Kcpa        float64 `json:"kcpa"`
	Kcpw        float64 `json:"kcpw"`
}

type DeteriorationConcrete struct {
	RoadGroupId int     `json:"road_group_id"`
	PSteel      float64 `json:"p_steel"`
	Ec          float64 `json:"ec"`
	Mi          float64 `json:"mi"`
	Fi          float64 `json:"fi"`
	Kjrc        float64 `json:"kjrc"`
	BStress     float64 `json:"b_stress"`
	JtSpace     float64 `json:"jt_space"`
	Kjrf        float64 `json:"kjrf"`
	Widened     float64 `json:"widened"`
	PredSeal    float64 `json:"pred_seal"`
	DwlCor      float64 `json:"dwl_cor"`
	Kjrs        float64 `json:"kjrs"`
	Kjrr        float64 `json:"kjrr"`
}

func (b *SettingDeteriorationParams) TableName() string {
	return "setting_deterioration_params"
}

func (b *SettingDeterioration) TableName() string {
	return "setting_deterioration"
}
