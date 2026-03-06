package models

type RefOwner struct {
	ID                  int    `json:"id" example:"1"`
	Name                string `json:"name" example:"BEM"`
	IsActive            bool   `json:"-"`
	RefConditionRangeID int    `json:"ref_condition_range_id"`
}

type RefOwnerPreload struct {
	ID                  int               `json:"id" example:"1"`
	Name                string            `json:"name" example:"BEM"`
	IsActive            bool              `json:"-"`
	RefConditionRangeID int               `json:"ref_condition_range_id"`
	RefConditionRange   RefConditionRange `json:"ref_condition_range" gorm:"foreignKey:RefConditionRangeID;references:ID"`
}

func (ro *RefOwnerPreload) TableName() string {
	return "ref_owner"
}

func (ro *RefOwner) TableName() string {
	return "ref_owner"
}

type RefOwnerRoadLine struct {
	ID                     int    `json:"id" example:"1"`
	Name                   string `json:"name" example:"BEM"`
	IsActive               bool   `json:"-"`
	RefReflectivityRangeID int    `json:"ref_reflectivity_range_id"`
}

type RefOwnerRoadLineInit struct {
	ID                     int    `json:"id" example:"1"`
	Name                   string `json:"name" example:"BEM"`
	IsActive               bool   `json:"-"`
	RefReflectivityRangeID int    `json:"ref_reflectivity_range_id"`
}
type RefOwnerRoadLinePreload struct {
	ID                     int                  `json:"id" example:"1"`
	Name                   string               `json:"name" example:"BEM"`
	IsActive               bool                 `json:"-"`
	RefReflectivityRangeID int                  `json:"ref_reflectivity_range_id"`
	RefReflectivityRange   RefReflectivityRange `json:"ref_reflectivity_range"`
}

func (ro *RefOwnerRoadLine) TableName() string {
	return "ref_owner_road_line"
}

func (ro *RefOwnerRoadLineInit) TableName() string {
	return "ref_owner_road_line"
}
func (ro *RefOwnerRoadLinePreload) TableName() string {
	return "ref_owner_road_line"
}
