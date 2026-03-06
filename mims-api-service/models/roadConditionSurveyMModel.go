package models

type RoadConditionSurveyM struct {
	ID                        int      `gorm:"column:id;primaryKey;autoIncrement"`
	RoadConditionSurvay100mID int      `gorm:"column:road_condition_survey_100m_id"`
	KmStart                   float64  `gorm:"column:km_start"`
	KmEnd                     float64  `gorm:"column:km_end"`
	IRI                       *float64 `gorm:"column:iri"`
	MPD                       *float64 `gorm:"column:mpd"`
	RUT                       *float64 `gorm:"column:rut"`
	IFI                       *float64 `gorm:"column:ifi"`
	// RefGradeIriID         int     `gorm:"column:ref_grade_iri_id"`
	// RefGradeMpdID         int     `gorm:"column:ref_grade_mpd_id"`
	// RefGradeRutID         int     `gorm:"column:ref_grade_rut_id"`
	// RefGradeIfiID         int     `gorm:"column:ref_grade_ifi_id"`
	// RefGradeGnID          int     `gorm:"column:ref_grade_gn_id"`
	SurveyType  string `gorm:"column:survey_type"`
	TheGeom     string `gorm:"column:the_geom"`
	ImgFilepath string `gorm:"column:img_filepath"`
}
type RoadConditionSurveyM2 struct {
	ID                        int      `gorm:"primaryKey;column:id"`
	RoadConditionSurvay100mID int      `gorm:"column:road_condition_survey_100m_id"`
	KmStart                   float64  `gorm:"column:km_start"`
	KmEnd                     float64  `gorm:"column:km_end"`
	IRI                       *float64 `gorm:"column:iri"`
	MPD                       *float64 `gorm:"column:mpd"`
	RUT                       *float64 `gorm:"column:rut"`
	IFI                       *float64 `gorm:"column:ifi"`
	SurveyType                string   `gorm:"column:survey_type"`
	// RefGradeIriID         int     `gorm:"column:ref_grade_iri_id"`
	// RefGradeMpdID         int     `gorm:"column:ref_grade_mpd_id"`
	// RefGradeRutID         int     `gorm:"column:ref_grade_rut_id"`
	// RefGradeIfiID         int     `gorm:"column:ref_grade_ifi_id"`
	// RefGradeGnID          int     `gorm:"column:ref_grade_gn_id"`
	Geojson     []byte `gorm:"column:geojson"`
	ImgFilepath string `gorm:"column:img_filepath"`
}

func (rc *RoadConditionSurveyM) TableName() string {
	return "road_condition_survey_m"
}

func (rc *RoadConditionSurveyM2) TableName() string {
	return "road_condition_survey_m"
}

type RoadConditionSurveyMPreload struct {
	ID                        int      `gorm:"primaryKey;column:id"`
	RoadConditionSurvay100mID int      `gorm:"column:road_condition_survey_100m_id"`
	KmStart                   float64  `gorm:"column:km_start"`
	KmEnd                     float64  `gorm:"column:km_end"`
	IRI                       *float64 `gorm:"column:iri"`
	MPD                       *float64 `gorm:"column:mpd"`
	RUT                       *float64 `gorm:"column:rut"`
	IFI                       *float64 `gorm:"column:ifi"`
	SurveyType                string   `gorm:"column:survey_type"`
	// RefGradeIriID         int        `gorm:"column:ref_grade_iri_id"`
	// RefGradeMpdID         int        `gorm:"column:ref_grade_mpd_id"`
	// RefGradeRutID         int        `gorm:"column:ref_grade_rut_id"`
	// RefGradeIfiID         int        `gorm:"column:ref_grade_ifi_id"`
	// RefGradeGnID          int        `gorm:"column:ref_grade_gn_id"`
	TheGeom     string `gorm:"column:the_geom"`
	ImgFilepath string `gorm:"column:img_filepath"`
	// RefGradeIri           []RefGrade `json:"ref_grade_iri" gorm:"ForeignKey:ID;AssociationForeignKey:RefGradeIriID"`
	// RefGradeMpd           []RefGrade `json:"ref_grade_mpd" gorm:"ForeignKey:ID;AssociationForeignKey:RefGradeMpdID"`
	// RefGradeRut           []RefGrade `json:"ref_grade_rut" gorm:"ForeignKey:ID;AssociationForeignKey:RefGradeRutID"`
	// RefGradeIfi           []RefGrade `json:"ref_grade_ifi" gorm:"ForeignKey:ID;AssociationForeignKey:RefGradeIfiID"`
}

func (rc *RoadConditionSurveyMPreload) TableName() string {
	return "road_condition_survey_m"
}

type RoadConditionSurveyMDashboard struct {
	ID                        int      `gorm:"column:id;primaryKey;autoIncrement"`
	RoadConditionSurvay100mID int      `gorm:"column:road_condition_survey_100m_id"`
	KmStart                   float64  `gorm:"column:km_start"`
	KmEnd                     float64  `gorm:"column:km_end"`
	IRI                       *float64 `gorm:"column:iri"`
	MPD                       *float64 `gorm:"column:mpd"`
	RUT                       *float64 `gorm:"column:rut"`
	IFI                       *float64 `gorm:"column:ifi"`
	// RefGradeIriID         int     `gorm:"column:ref_grade_iri_id"`
	// RefGradeMpdID         int     `gorm:"column:ref_grade_mpd_id"`
	// RefGradeRutID         int     `gorm:"column:ref_grade_rut_id"`
	// RefGradeIfiID         int     `gorm:"column:ref_grade_ifi_id"`
	// RefGradeGnID          int     `gorm:"column:ref_grade_gn_id"`
	SurveyType     string `gorm:"column:survey_type"`
	TheGeom        string `gorm:"column:the_geom"`
	SurveyMTheGeom []byte `gorm:"column:survey_m_the_geom"`
	ImgFilepath    string `gorm:"column:img_filepath"`
}

func (rc *RoadConditionSurveyMDashboard) TableName() string {
	return "road_condition_survey_m"
}
