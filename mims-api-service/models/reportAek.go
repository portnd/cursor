package models

type RoadConditionReport struct {
	RoadCondition
	Road                 RoadForDashboard             `json:"road" gorm:"ForeignKey:RoadId;AssociationForeignKey:Id"`
	RoadInfo             RoadInfo                     `json:"road_info" gorm:"ForeignKey:RoadId;references:RoadId"`
	RoadConditionSurveys []RoadConditionSurveyPreload ` gorm:"ForeignKey:RoadConditionID;references:ID"`
}

func (rc *RoadConditionReport) TableName() string {
	return "road_condition"
}

type PositionDamage struct {
	LaneNO      int
	Km          int
	StrKm       string
	ImgFilepath string

	ACIcrack       float64
	ACUcrack       float64
	ACRavelling    float64
	ACPatching     float64
	ACPotholeArea  float64
	ACBleeding     float64
	ACPotholeCount float64

	CCTransverseCrack    float64
	CCNonTransverseCrack float64
	CCSpalling           float64
	CCCornerBreaks       float64
	CCJointSealDamage    float64
	CCPatching           float64
	CCScaling            float64
}

type DataRoadDamage struct {
	LaneNo int

	ACIcrack       float64
	ACUcrack       float64
	ACRavelling    float64
	ACPatching     float64
	ACPotholeArea  float64
	ACBleeding     float64
	ACPotholeCount float64

	CCTransverseCrack    float64
	CCNonTransverseCrack float64
	CCSpalling           float64
	CCCornerBreaks       float64
	CCJointSealDamage    float64
	CCPatching           float64
	CCScaling            float64

	Position []CovertPositionDamage
}
type DataRoadDamageModel struct {
	LaneNo int

	ACIcrack       float64
	ACUcrack       float64
	ACRavelling    float64
	ACPatching     float64
	ACPotholeArea  float64
	ACBleeding     float64
	ACPotholeCount float64

	CCTransverseCrack    float64
	CCNonTransverseCrack float64
	CCSpalling           float64
	CCCornerBreaks       float64
	CCJointSealDamage    float64
	CCPatching           float64
	CCScaling            float64

	// Position []CovertPositionDamage
}

type Report11 struct {
	Year              int
	Type              string
	RoadGroupName     string
	RoadSectionNumber string
	RoadSectionName   string
	KmStart           string
	KmEnd             string
	RoadLength        string
	IsNull            bool
	Data              DataRoadDamage
}

type Report10 struct {
	Year              int
	Type              string
	RoadGroupName     string
	RoadSectionNumber string
	RoadSectionName   string
	KmStart           string
	KmEnd             string
	RoadLength        string
	IsNull            bool
	Data              []DataReportDamage
}

type DataReportDamage struct {
	RoadGroupName string
	RoadName      string
	RoadCode      string
	KmStart       int
	KmEnd         int
	StrKmStart    string
	StrKmEnd      string
	RoadLengthStr string
	IsNull        bool

	Detail []DataRoadDamage
}



type DataReportDamageModel struct {
	RoadGroupName string
	RoadName      string
	RoadCode      string
	KmStart       int
	KmEnd         int
	StrKmStart    string
	StrKmEnd      string
	RoadLengthStr string
	IsNull        bool

}

type CovertPositionDamage struct {
	No            int
	Km            int
	StrKm         string
	Surface       string
	DamageType    string
	DamageTypeENG string
	Value         float64
	Unit          string
	Image         string
}
