package models

type RoadConditionSurvey100M struct {
	ID                    int      `gorm:"column:id;primaryKey;autoIncrement"`
	RoadConditionSurveyId int      `gorm:"column:road_condition_survey_id"`
	KmStart               float64  `gorm:"column:km_start"`
	KmEnd                 float64  `gorm:"column:km_end"`
	IRI                   *float64 `gorm:"column:iri"`
	MPD                   *float64 `gorm:"column:mpd"`
	RUT                   *float64 `gorm:"column:rut"`
	IFI                   *float64 `gorm:"column:ifi"`
	SurveyType            string   `gorm:"column:survey_type"`
	TheGeom               string   `gorm:"column:the_geom"`
}

func (rc *RoadConditionSurvey100M) TableName() string {
	return "road_condition_survey_100m"
}

type RoadConditionSurvey100M2 struct {
	ID                    int      `gorm:"column:id;primaryKey;autoIncrement"`
	RoadConditionSurveyId int      `gorm:"column:road_condition_survey_id"`
	KmStart               float64  `gorm:"column:km_start"`
	KmEnd                 float64  `gorm:"column:km_end"`
	IRI                   *float64 `gorm:"column:iri"`
	MPD                   *float64 `gorm:"column:mpd"`
	RUT                   *float64 `gorm:"column:rut"`
	IFI                   *float64 `gorm:"column:ifi"`
	SurveyType            string   `gorm:"column:survey_type"`
	TheGeom               []byte   `gorm:"column:geojson"`
}

func (rc *RoadConditionSurvey100M2) TableName() string {
	return "road_condition_survey_100m"
}

type RoadConditionSurvey100MPreload struct {
	ID                    int                    `gorm:"column:id;primaryKey;autoIncrement"`
	RoadConditionSurveyID int                    `gorm:"column:road_condition_survey_id"`
	RoadConditionSurveyMs []RoadConditionSurveyM `json:"road_condition_survay_m" gorm:"ForeignKey:RoadConditionSurvay100mID;references:ID"`
	KmStart               float64                `gorm:"column:km_start"`
	KmEnd                 float64                `gorm:"column:km_end"`
	IRI                   *float64               `gorm:"column:iri"`
	MPD                   *float64               `gorm:"column:mpd"`
	RUT                   *float64               `gorm:"column:rut"`
	IFI                   *float64               `gorm:"column:ifi"`
	SurveyType            string                 `gorm:"column:survey_type"`
	TheGeom               string                 `gorm:"column:the_geom"`
}

func (rc *RoadConditionSurvey100MPreload) TableName() string {
	return "road_condition_survey_100m"
}

type RoadConditionSurvey100MDashboard struct {
	ID                    int                             `gorm:"column:id;primaryKey;autoIncrement"`
	RoadConditionSurveyID int                             `gorm:"column:road_condition_survey_id"`
	RoadConditionSurveyMs []RoadConditionSurveyMDashboard `json:"road_condition_survay_m" gorm:"ForeignKey:RoadConditionSurvay100mID;references:ID"`
	KmStart               float64                         `gorm:"column:km_start"`
	KmEnd                 float64                         `gorm:"column:km_end"`
	IRI                   *float64                        `gorm:"column:iri"`
	MPD                   *float64                        `gorm:"column:mpd"`
	RUT                   *float64                        `gorm:"column:rut"`
	IFI                   *float64                        `gorm:"column:ifi"`
	SurveyType            string                          `gorm:"column:survey_type"`
	TheGeom               string                          `gorm:"column:the_geom"`
	Survey100MTheGeom     []byte                          `gorm:"column:survey_100m_the_geom"`
}

func (rc *RoadConditionSurvey100MDashboard) TableName() string {
	return "road_condition_survey_100m"
}
