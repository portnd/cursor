package models

type ParamsCondition struct {
	ID         int
	RefOwnerID int
	RefGradeID int `gorm:"column:ref_grade_id"`
	// RefConditionRangeID int
	LeftValueAC      float64 `gorm:"column:left_value_ac"`
	LeftConditionAC  string  `gorm:"column:left_condition_ac"`
	RightValueAC     float64 `gorm:"column:right_value_ac"`
	RightConditionAC string  `gorm:"column:right_condition_ac"`
	LeftValueCC      float64 `gorm:"column:left_value_cc"`
	LeftConditionCC  string  `gorm:"column:left_condition_cc"`
	RightValueCC     float64 `gorm:"column:right_value_cc"`
	RightConditionCC string  `gorm:"column:right_condition_cc"`
	ConditionType    string
}

type ParamsRoadLine struct {
	ID                 int
	RefOwnerRoadLineID int `gorm:"column:ref_owner_road_line_id"`
	RefGradeID         int
	// RefConditionRangeID  int
	LeftValueWhite       float64 `gorm:"column:left_value_white"`
	LeftConditionWhite   string  `gorm:"column:left_condition_white"`
	RightValueWhite      float64 `gorm:"column:right_value_white"`
	RightConditionWhite  string  `gorm:"column:right_condition_white"`
	LeftValueYellow      float64 `gorm:"column:left_value_yellow"`
	LeftConditionYellow  string  `gorm:"column:left_condition_yellow"`
	RightValueYellow     float64 `gorm:"column:right_value_yellow"`
	RightConditionYellow string  `gorm:"column:right_condition_yellow"`
	// ConditionType    string
}
type ParamsRoadLinePreload struct {
	ParamsRoadLine
	RefOwnerRoadLine RefOwnerRoadLine `gorm:"foreignKey:RefOwnerRoadLineID;references:ID"`
	RefGrade         RefGrade         `gorm:"foreignKey:RefGradeID;references:ID"`
}

type ParamsConditionPreload struct {
	ID               int             `gorm:"column:id"`
	RefOwnerID       int             `gorm:"column:ref_owner_id"`
	RefOwner         RefOwnerPreload `gorm:"foreignKey:RefOwnerID;references:ID"`
	RefGradeID       int             `gorm:"column:ref_grade_id"`
	RefGrade         RefGrade        `gorm:"foreignKey:RefGradeID;references:ID"`
	LeftValueAC      float64         `gorm:"column:left_value_ac"`
	LeftConditionAC  string          `gorm:"column:left_condition_ac"`
	RightValueAC     float64         `gorm:"column:right_value_ac"`
	RightConditionAC string          `gorm:"column:right_condition_ac"`
	LeftValueCC      float64         `gorm:"column:left_value_cc"`
	LeftConditionCC  string          `gorm:"column:left_condition_cc"`
	RightValueCC     float64         `gorm:"column:right_value_cc"`
	RightConditionCC string          `gorm:"column:right_condition_cc"`
	ConditionType    string
}

type RefConditionRange struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (pc *RefConditionRange) TableName() string {
	return "ref_condition_range"
}

type RefReflectivityRange struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (pc *RefReflectivityRange) TableName() string {
	return "ref_reflectivity_range"
}

func (pc *ParamsConditionPreload) TableName() string {
	return "params_condition"
}

func (pc *ParamsCondition) TableName() string {
	return "params_condition"
}

func (pc *ParamsRoadLinePreload) TableName() string {
	return "params_road_line"
}

func (pc *ParamsRoadLine) TableName() string {
	return "params_road_line"
}
