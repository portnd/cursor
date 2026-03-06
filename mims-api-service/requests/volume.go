package requests

type VolumeAadtReq struct {
	SurveyedDate string `json:"surveyed_date" validate:"nonzero"`
	Veh1         int    `json:"veh1" validate:"min=0"`
	Veh2         int    `json:"veh2" validate:"min=0"`
	Veh3         int    `json:"veh3" validate:"min=0"`
	Stauts       string `json:"-"`
}

type VolumeAccidentReq struct {
	SurveyedDate string `json:"surveyed_date" validate:"nonzero"`
	Acc1         int    `json:"acc1" validate:"min=0"`
	Acc2         int    `json:"acc2" validate:"min=0"`
	Acc3         int    `json:"acc3" validate:"min=0"`
	Acc4         int    `json:"acc4" validate:"min=0"`
	Stauts       string `json:"-"`
}
