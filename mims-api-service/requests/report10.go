package requests

type Report13 struct {
	TypeReport    string `json:"type" form:"type"`
	YearStart     int    `json:"year_start" form:"year_start"`
	YearEnd       int    `json:"year_end" form:"year_end"`
	Group         int    `json:"group"`
	RoadSectionId int    `json:"road_section_id" form:"road_section_id"`
}
