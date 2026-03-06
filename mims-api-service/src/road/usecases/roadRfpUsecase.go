package usecases

import (
	"gitlab.com/mims-api-service/models"
	requests "gitlab.com/mims-api-service/requests"
	responses "gitlab.com/mims-api-service/responses"
)

func (u *roadUseCase) GetRfp(params requests.RoadPrams) (interface{}, error) {
	roadIDs, _ := u.roadRepo.GetRoadID(params)

	// ดึงข้อมูลเกณฑ์หน้าสภาพทาง 1000 m
	conditions1000, err := u.roadRepo.GetParamsConditionByID(92, "")
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	iriCondition1000AC := conditions1000.RightValueAC
	iriCondition1000CC := conditions1000.RightValueCC

	//ดึงข้อมูลเกณฑ์หน้าสภาพทาง 100 m
	// IRI
	iriConditions100, err := u.roadRepo.GetParamsConditionByID(91, "IRI")
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	iriCondition100AC := iriConditions100.RightValueAC
	iriCondition100CC := iriConditions100.RightValueCC

	// RUT
	rutConditions100, err := u.roadRepo.GetParamsConditionByID(91, "RUT")
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	rutCondition100AC := rutConditions100.RightValueAC
	rutCondition100CC := rutConditions100.RightValueCC

	// IFI
	ifiConditions100, err := u.roadRepo.GetParamsConditionByID(91, "IFI")
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	ifiCondition100AC := ifiConditions100.RightValueAC
	ifiCondition100CC := ifiConditions100.RightValueCC

	//ดึงข้อมูลเกณฑ์ G7 100 m
	g7100, err := u.roadRepo.GetParamsRoadLine(13)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	g7ConditionYellow := g7100.LeftValueYellow
	g7ConditionWhite := g7100.LeftValueWhite

	roadSurfaceLanes, err := u.roadRepo.GetRoadSurfaceLaneAll(roadIDs)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}

	var res responses.DataRes
	dataIRI1000 := make(map[int][]int)
	dataIRI100 := make(map[int][]int)
	dataRUT100 := make(map[int][]int)
	dataIFI100 := make(map[int][]int)
	dataG7100 := make(map[int][]int)
	for _, roadID := range roadIDs {
		// IRI 1000 m
		if params.IsIri1000 != nil {
			iri1000Data, err := u.Iri1000(params, roadSurfaceLanes, roadID, iriCondition1000AC, iriCondition1000CC)
			if err != nil {
				return nil, responses.NewAppErr(400, err.Error())
			}
			if iri1000Data {
				dataIRI1000[1] = append(dataIRI1000[1], roadID)
			} else {
				dataIRI1000[0] = append(dataIRI1000[0], roadID)
			}
		}

		// IRI 100 m
		if params.IsIri100 != nil {
			iri100Data, err := u.Iri100(params, roadSurfaceLanes, roadID, iriCondition100AC, iriCondition100CC)
			if err != nil {
				return nil, responses.NewAppErr(400, err.Error())
			}
			if iri100Data {
				dataIRI100[1] = append(dataIRI100[1], roadID)
			} else {
				dataIRI100[0] = append(dataIRI100[0], roadID)
			}
		}

		// RUT 100 m
		if params.IsRut100 != nil {
			rut100Data, err := u.Rut100(params, roadSurfaceLanes, roadID, rutCondition100AC, rutCondition100CC)
			if err != nil {
				return nil, responses.NewAppErr(400, err.Error())
			}
			if rut100Data {
				dataRUT100[1] = append(dataRUT100[1], roadID)
			} else {
				dataRUT100[0] = append(dataRUT100[0], roadID)
			}
		}

		// IFI 100 m
		if params.IsIfi100 != nil {
			ifi100Data, err := u.Ifi100(params, roadSurfaceLanes, roadID, ifiCondition100AC, ifiCondition100CC)
			if err != nil {
				return nil, responses.NewAppErr(400, err.Error())
			}
			if ifi100Data {
				dataIFI100[1] = append(dataIFI100[1], roadID)
			} else {
				dataIFI100[0] = append(dataIFI100[0], roadID)
			}
		}

		// G7 100 m
		if params.IsG7100 != nil {
			g7100Data, err := u.G7100(params, roadSurfaceLanes, roadID, g7ConditionYellow, g7ConditionWhite)
			if err != nil {
				return nil, responses.NewAppErr(400, err.Error())
			}
			if g7100Data {
				dataG7100[1] = append(dataG7100[1], roadID)
			} else {
				dataG7100[0] = append(dataG7100[0], roadID)
			}
		}
	}

	res.Iri1000 = dataIRI1000
	res.Iri100 = dataIRI100
	res.Rut100 = dataRUT100
	res.Ifi100 = dataIFI100
	res.G7100 = dataG7100
	return res, nil
}

func (u *roadUseCase) Iri1000(params requests.RoadPrams, roadSurfaceLanes map[int][]models.RoadSurfaceLane, roadID int, iriCondition1000AC, iriCondition1000CC float64) (bool, error) {
	checkIri1000 := 0
	if params.IsIri1000 != nil {
		for _, lane := range roadSurfaceLanes[roadID] {
			items, err := u.roadRepo.GetroadConditionSurvey(roadID, lane.LaneNo)
			if err != nil {
				return false, responses.NewAppErr(400, err.Error())
			}
			for _, item := range items {
				iri1000 := 0.0
				if item.IRI != nil {
					iri := *item.IRI
					if item.SurveyType == "AC" {
						iri1000 = iriCondition1000AC //ต้องใช้ค่า ชื่อเกณฑ์ 1 กม ประเภท 1 กม ค่า IRI ผิวทางลาดยาง
					} else {
						iri1000 = iriCondition1000CC //ต้องใช้ค่า ชื่อเกณฑ์ 1 กม ประเภท 1 กม ค่า IRI ผิวทางคอนกรีต
					}
					if iri <= iri1000 {
						continue
					} else {
						checkIri1000 = checkIri1000 + 1
						break
					}
				} else {
					continue
				}
			}
		}
	}
	if checkIri1000 == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (u *roadUseCase) Iri100(params requests.RoadPrams, roadSurfaceLanes map[int][]models.RoadSurfaceLane, roadID int, iriCondition100AC, iriCondition100CC float64) (bool, error) {
	checkIri100 := 0
	if params.IsIri100 != nil {
		for _, lane := range roadSurfaceLanes[roadID] {
			items, err := u.roadRepo.GetroadConditionSurvey100M(roadID, lane.LaneNo)
			if err != nil {
				return false, responses.NewAppErr(400, err.Error())
			}
			for _, item := range items {
				iri100 := 0.0
				if item.IRI != nil {
					iri := *item.IRI
					if item.SurveyType == "AC" {
						iri100 = iriCondition100AC //ต้องใช้ค่า ชื่อเกณฑ์ 1 กม ประเภท 1 กม ค่า IRI ผิวทางลาดยาง
					} else {
						iri100 = iriCondition100CC //ต้องใช้ค่า ชื่อเกณฑ์ 1 กม ประเภท 1 กม ค่า IRI ผิวทางคอนกรีต
					}
					if iri <= iri100 {
						continue
					} else {
						checkIri100 = checkIri100 + 1
						break
					}
				} else {
					continue
				}
			}
		}
	}
	if checkIri100 == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (u *roadUseCase) Rut100(params requests.RoadPrams, roadSurfaceLanes map[int][]models.RoadSurfaceLane, roadID int, rutCondition100AC, rutCondition100CC float64) (bool, error) {
	checkRut100 := 0
	if params.IsRut100 != nil {
		for _, lane := range roadSurfaceLanes[roadID] {
			items, err := u.roadRepo.GetroadConditionSurvey100M(roadID, lane.LaneNo)
			if err != nil {
				return false, responses.NewAppErr(400, err.Error())
			}
			for _, item := range items {
				rut100 := 0.0
				if item.RUT != nil {
					rut := *item.RUT
					if item.SurveyType == "AC" {
						rut100 = rutCondition100AC //ต้องใช้ค่า ชื่อเกณฑ์ 1 กม ประเภท 1 กม ค่า IRI ผิวทางลาดยาง
					} else {
						rut100 = rutCondition100CC //ต้องใช้ค่า ชื่อเกณฑ์ 1 กม ประเภท 1 กม ค่า IRI ผิวทางคอนกรีต
					}
					if rut <= rut100 {
						continue
					} else {
						checkRut100 = checkRut100 + 1
						break
					}
				} else {
					continue
				}
			}
		}
	}
	if checkRut100 == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (u *roadUseCase) Ifi100(params requests.RoadPrams, roadSurfaceLanes map[int][]models.RoadSurfaceLane, roadID int, ifiCondition100AC, ifiCondition100CC float64) (bool, error) {
	checkIfi100 := 0
	if params.IsIfi100 != nil {
		for _, lane := range roadSurfaceLanes[roadID] {
			items, err := u.roadRepo.GetroadConditionSurvey100M(roadID, lane.LaneNo)
			if err != nil {
				return false, responses.NewAppErr(400, err.Error())
			}
			for _, item := range items {
				ifi100 := 0.0
				if item.IFI != nil {
					ifi := *item.IFI
					if item.SurveyType == "AC" {
						ifi100 = ifiCondition100AC //ต้องใช้ค่า ชื่อเกณฑ์ 1 กม ประเภท 1 กม ค่า IRI ผิวทางลาดยาง
					} else {
						ifi100 = ifiCondition100CC //ต้องใช้ค่า ชื่อเกณฑ์ 1 กม ประเภท 1 กม ค่า IRI ผิวทางคอนกรีต
					}
					if ifi >= ifi100 {
						continue
					} else {
						checkIfi100 = checkIfi100 + 1
						break
					}
				} else {
					continue
				}
			}
		}
	}
	if checkIfi100 == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (u *roadUseCase) G7100(params requests.RoadPrams, roadSurfaceLanes map[int][]models.RoadSurfaceLane, roadID int, g7ConditionYellow, g7ConditionWhite float64) (bool, error) {
	checkRetro100 := 0
	if params.IsIfi100 != nil {
		for _, lane := range roadSurfaceLanes[roadID] {
			items, _ := u.roadRepo.GetRoadRetroReflectivity100M(roadID, lane.LaneNo)
			for _, item := range items {

				if item.RetroAvg != nil {
					retro := *item.RetroAvg
					nameStripID := item.RefStripeColorID
					retroCondition := 0.0
					if nameStripID == 1 { //เส้นสีขาว
						retroCondition = g7ConditionWhite
					} else {
						retroCondition = g7ConditionYellow
					}
					if retro >= retroCondition {
						continue
					} else {
						checkRetro100++
						break
					}
				} else {
					continue
				}
			}
		}
	}
	if checkRetro100 == 0 {
		return true, nil
	} else {
		return false, nil
	}
}
