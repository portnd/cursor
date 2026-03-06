package requests

type ChangedStatusRequest struct {
	Category string `form:"category" validate:"nonzero"`
	Status   string `form:"status" validate:"nonzero"`
	RoadId   string `form:"road_id"`
	Page     string `form:"page"`
	Limit    string `form:"limit"`
}

type ChangedVolumeStatusRequest struct {
	Category    string `form:"category" validate:"nonzero"`
	Status      string `form:"status" validate:"nonzero"`
	RoadGroupID string `form:"road_group_id"`
	Page        string `form:"page"`
	Limit       string `form:"limit"`
}

type ChangeDataStatusRequest struct {
	Category     string `json:"category" validate:"nonzero" extensions:"x-order=0"`
	Status       string `json:"status" validate:"nonzero" extensions:"x-order=1"`
	RejectReason string `json:"reject_reason" extensions:"x-order=2"`
	IdParent     []uint `json:"id_parent" extensions:"x-order=3"`
	AssetId      []uint `json:"asset_id" extensions:"x-order=4"`
	RoadId       []uint `json:"road_id" extensions:"x-order=5"`
}

type ChangeDataVolumeStatusRequest struct {
	Category string `json:"category" validate:"nonzero" extensions:"x-order=0"`
	// StatusCurrent string `json:"status_current" validate:"nonzero" extensions:"x-order=1"`
	Status       string `json:"status" validate:"nonzero" extensions:"x-order=2"`
	RejectReason string `json:"reject_reason" extensions:"x-order=3"`
	IdParent     []uint `json:"id_parent" extensions:"x-order=4"`
	RoadGroupId  []uint `json:"road_group_id" extensions:"x-order=6"`
}
type ChangeAssetDetailRequest struct {
	IdParent uint   `form:"id_parent" validate:"nonzero"`
	AssetId  uint   `form:"asset_id" validate:"nonzero"`
	Page     string `form:"page"`
	Limit    string `form:"limit"`
}

type ChangeSurfaceDetailRequest struct {
	RoadId uint   `form:"road_id" validate:"nonzero"`
	Page   string `form:"page"`
	Limit  string `form:"limit"`
}

type ChangeConditionRequest struct {
	IdParent      uint   `form:"id_parent" validate:"nonzero"`
	ConditionType string `form:"condition_type" validate:"nonzero"`
	Page          string `form:"page"`
	Limit         string `form:"limit"`
}

type ChangeDamageRequest struct {
	IdParent uint   `form:"id_parent" validate:"nonzero"`
	Page     string `form:"page"`
	Limit    string `form:"limit"`
}

type ChangeAADTRequest struct {
	IdParent    uint   `form:"id_parent" validate:"nonzero"`
	RoadGroupID uint   `form:"road_group_id" validate:"nonzero"`
	Status      string `form:"status" validate:"nonzero"`
}

type ChangeAccidentRequest struct {
	IdParent    uint   `form:"id_parent" validate:"nonzero"`
	RoadGroupID uint   `form:"road_group_id" validate:"nonzero"`
	Status      string `form:"status" validate:"nonzero"`
}
