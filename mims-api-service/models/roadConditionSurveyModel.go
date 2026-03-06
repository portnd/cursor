package models

type RoadConditionSurvey struct {
	ID              int      `gorm:"column:id;primaryKey;autoIncrement"`
	RoadConditionId int      `gorm:"column:road_condition_id"`
	KmStart         float64  `gorm:"column:km_start"`
	KmEnd           float64  `gorm:"column:km_end"`
	IRI             *float64 `gorm:"column:iri"`
	MPD             *float64 `gorm:"column:mpd"`
	RUT             *float64 `gorm:"column:rut"`
	IFI             *float64 `gorm:"column:ifi"`
	SurveyType      string   `gorm:"column:survey_type"`
	TheGeom         string   `gorm:"column:the_geom"`
}

func (rc *RoadConditionSurvey) TableName() string {
	return "road_condition_survey"
}

type RoadConditionSurveyPreload struct {
	ID                       int                              `gorm:"primaryKey column:id"`
	RoadConditionID          int                              `gorm:"column:road_condition_id"`
	RoadConditionSurvey100Ms []RoadConditionSurvey100MPreload `json:"road_condition_survay_100_m" gorm:"ForeignKey:RoadConditionSurveyID;references:ID"`
	KmStart                  float64                          `gorm:"column:km_start"`
	KmEnd                    float64                          `gorm:"column:km_end"`
	IRI                      *float64                         `gorm:"column:iri"`
	MPD                      *float64                         `gorm:"column:mpd"`
	RUT                      *float64                         `gorm:"column:rut"`
	IFI                      *float64                         `gorm:"column:ifi"`
	SurveyType               string                           `gorm:"column:survey_type"`
	TheGeom                  string                           `gorm:"column:the_geom"`
}

func (rc *RoadConditionSurveyPreload) TableName() string {
	return "road_condition_survey"
}

type RoadConditionSurveyDashboard struct {
	ID                       int                                `gorm:"primaryKey column:id"`
	RoadConditionID          int                                `gorm:"column:road_condition_id"`
	RoadConditionSurvey100Ms []RoadConditionSurvey100MDashboard `json:"road_condition_survay_100_m" gorm:"ForeignKey:RoadConditionSurveyID;references:ID"`
	KmStart                  float64                            `gorm:"column:km_start"`
	KmEnd                    float64                            `gorm:"column:km_end"`
	IRI                      *float64                           `gorm:"column:iri"`
	MPD                      *float64                           `gorm:"column:mpd"`
	RUT                      *float64                           `gorm:"column:rut"`
	IFI                      *float64                           `gorm:"column:ifi"`
	SurveyType               string                             `gorm:"column:survey_type"`
	TheGeom                  string                             `gorm:"column:the_geom"`
	SurveyTheGeom            []byte                             `gorm:"column:survey_the_geom"`
}

func (rc *RoadConditionSurveyDashboard) TableName() string {
	return "road_condition_survey"
}
