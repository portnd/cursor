package responses

type AssetMap struct {
	ID                    int      `json:"id"`
	RoadID                int      `json:"road_id"`
	AssetTableID          int      `json:"asset_table_id"`
	Name                  string   `json:"name"`
	IconFilepath          string   `json:"icon_filepath"`
	ThumbnailIconFilepath string   `json:"thumbnail_icon_filepath"`
	LineColorCode         string   `json:"line_color_code"`
	IsCluster             bool     `json:"is_cluster"`
	Cluster               int      `json:"cluster"`
	TheGeom               GeomJSON `json:"the_geom"`
}
