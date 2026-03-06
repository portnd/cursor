package models

type RefAssetTube struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rat *RefAssetTube) TableName() string {
	return "ref_asset_tube"
}
