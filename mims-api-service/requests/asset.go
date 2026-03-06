package requests

type AssetMap struct {
	KmStart float64 `form:"km_start"`
	KmEnd   float64 `form:"km_end"`
	Year    string  `form:"year"`
	Left    string  `form:"left"`
	Right   string  `form:"right"`
	Bottom  string  `form:"bottom"`
	Top     string  `form:"top"`
	Zoom    string  `form:"zoom"`
}

type Asset struct {
	Page    int     `form:"page"`
	Limit   int     `form:"limit"`
	KmStart float64 `form:"km_start"`
	KmEnd   float64 `form:"km_end"`
	Year    string  `form:"year"`
}

type Condition struct {
	KmStart          float64 `form:"km_start"`
	KmEnd            float64 `form:"km_end"`
	Year             int     `form:"year"`
	ConditionType    int     `form:"condition_type"`
	ConditionOwnerID int     `form:"condition_owner_id"`
}

type ConditionMap struct {
	KmStart          float64 `form:"km_start"`
	KmEnd            float64 `form:"km_end"`
	Year             int     `form:"year"`
	ConditionType    int     `form:"condition_type"`
	ConditionOwnerID int     `form:"condition_owner_id"`
	LaneNo           int     `form:"lane_no"`
	Left             string  `form:"left"`
	Right            string  `form:"right"`
	Bottom           string  `form:"bottom"`
	Top              string  `form:"top"`
}
