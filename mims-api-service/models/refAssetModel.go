package models

type RefAsset struct {
	ID        int    `json:"id" extensions:"x-order=0"`
	Name      string `json:"name" extensions:"x-order=1"`
	Seq       int    `json:"-"`
	Status    int    `json:"-"`
	CanDelete bool   `json:"can_delete" extensions:"x-order=2"`
}

func (ra *RefAsset) TableName() string {
	return "ref_asset"
}

type RefAssetWeightStationType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rat *RefAssetWeightStationType) TableName() string {
	return "ref_asset_weight_station_type"
}

type RefAssetReflecType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (ract *RefAssetReflecType) TableName() string {
	return "ref_asset_reflec_type"
}

type RefAssetNoiseBarrier struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (raf *RefAssetNoiseBarrier) TableName() string {
	return "ref_asset_noise_barrier"
}

type RefAssetBuildingType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rae *RefAssetBuildingType) TableName() string {
	return "ref_asset_building_type"
}

type RefAssetTrafficCameraType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rae *RefAssetTrafficCameraType) TableName() string {
	return "ref_asset_traffic_camera_type"
}

type RefAssetKmstoneType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rad *RefAssetKmstoneType) TableName() string {
	return "ref_asset_kmstone_type"
}

type RefAssetFenceType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rac *RefAssetFenceType) TableName() string {
	return "ref_asset_fence_type"
}

type RefAssetCleranceType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (racp *RefAssetCleranceType) TableName() string {
	return "ref_asset_clerance_type"
}

type RefAssetCrashcushionType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (ract *RefAssetCrashcushionType) TableName() string {
	return "ref_asset_crashcushion_type"
}

type RefAssetLightWatt struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rac *RefAssetLightWatt) TableName() string {
	return "ref_asset_light_watt"
}

type RefAssetLightType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rab *RefAssetLightType) TableName() string {
	return "ref_asset_light_type"
}

type RefAssetOwner struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rab *RefAssetOwner) TableName() string {
	return "ref_asset_owner"
}

type RefAssetQutterType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (racp *RefAssetQutterType) TableName() string {
	return "ref_asset_qutter_type"
}

type RefAssetCrashcushion struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rac *RefAssetCrashcushion) TableName() string {
	return "ref_asset_crashcushion"
}

type RefAssetSignType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rast *RefAssetSignType) TableName() string {
	return "ref_asset_sign_type"
}

type RefAssetPosition struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rap *RefAssetPosition) TableName() string {
	return "ref_asset_position"
}

type RefAssetArea struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (ract *RefAssetArea) TableName() string {
	return "ref_asset_area"
}

type RefAssetSign struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (ras *RefAssetSign) TableName() string {
	return "ref_asset_sign"
}

type RefAssetSignImage struct {
	ID                int    `json:"id" example:"1"`
	Name              string `json:"name" example:"ลดความเร็ว"`
	Abbr              string `json:"abbr" example:"ต"`
	SignImageFilepath string `json:"sign_image_filepath" gorm:"column:sign_image_filepath" example:"public://attachments/settings/sign/20210209_152112_T-reducespeed.png"`
	StatusCode        string `json:"-"`
	Status            int    `json:"-"`
}

func (rasi *RefAssetSignImage) TableName() string {
	return "ref_asset_sign_image"
}

type RefAssetGuardrail struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rag *RefAssetGuardrail) TableName() string {
	return "ref_asset_guardrail"
}
