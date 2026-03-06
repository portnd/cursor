package models

type AssetLocation struct {
	ID                    int    `json:"id"`
	RoadID                int    `json:"road_id"`
	Name                  string `json:"name"`
	Type                  string `json:"type"`
	IconFilepath          string `json:"icon_filepath"`
	ThumbnailIconFilepath string `json:"thumbnail_icon_filepath"`
	LineColorCode         string `json:"line_color_code"`
	Wkt                   []byte `json:"wkt"`
}

type TableResult struct {
	ID            int
	TableName     string
	TableLabel    string
	IconFilepath  string
	LineColorCode string
	AssetId       int
}

type RoadAssetSign struct {
	ID                int     `gorm:"column:id"`
	RoadID            int     `gorm:"column:road_id"`
	TheGeom           []byte  `gorm:"column:the_geom"`
	Longitude         float64 `gorm:"column:lon"`
	Latitude          float64 `gorm:"column:lat"`
	SignImageFilePath string  `gorm:"column:sign_image_filepath"`
	ImgFilePath       string  `gorm:"column:img_filepath"`
}

type RoadAssetGeomCuster struct {
	//AssetID             int    `gorm:"column:asset_id"`
	TheGeomCluster      []byte `gorm:"column:the_geom_cluster"`
	TotalTheGeomCluster int    `gorm:"column:total_the_geom_cluster"`
}
