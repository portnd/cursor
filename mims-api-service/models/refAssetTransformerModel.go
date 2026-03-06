package models

type RefAssetTransformer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (rat *RefAssetTransformer) TableName() string {
	return "ref_asset_transformer"
}
