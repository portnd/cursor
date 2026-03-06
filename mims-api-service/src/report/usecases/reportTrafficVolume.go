package usecases

import (
	"errors"
	"fmt"
	"sort"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	"gorm.io/gorm"
)

func (u *UseCase) GetReportTrafficVolume(roadSectionIDs []string, typ string, year int) (interface{}, error) {

	if len(roadSectionIDs) == 0 {
		return nil, responses.NewAppErr(400, "กรุณาเลือกตอนควบคุม")
	}

	if typ != "excel" && typ != "pdf" && typ != "html" {
		if typ == "" {
			return nil, responses.NewAppErr(400, "กรุณาเลือกประเภทรายงาน")
		}
		return nil, responses.NewAppErr(400, "ประเภทรายงานไม่ถูกต้อง")
	}

	if year == 0 {
		return nil, responses.NewAppErr(400, "กรุณาเลือกปี")
	}

	roadDetails, err := u.Repo.GetRoadDetailsByRoadSectionID(roadSectionIDs)
	if err != nil {
		logs.Error(err)
		return []responses.ReportTrafficVolume{}, err
	}

	if len(roadDetails) == 0 {
		return []responses.ReportTrafficVolume{}, nil

	}

	var roadIDs []int

	for _, roadDetail := range roadDetails {

		roadIDs = append(roadIDs, roadDetail.RoadID)

	}

	trafficVolumes, err := u.Repo.GetReportTrafficVolume(roadIDs, year)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		logs.Error(err)
		return responses.RoadConditionDetails{}, err
	}

	trafficVolumeMap := make(map[int]models.VolumeAadt)
	for _, trafficVolume := range trafficVolumes {
		trafficVolumeMap[trafficVolume.RoadId] = trafficVolume
	}

	var resps []responses.ReportTrafficVolume

	for _, roadDetail := range roadDetails {
		var resp responses.ReportTrafficVolume

		resp.RoadID = roadDetail.RoadID
		resp.RoadGroupName = roadDetail.RoadGroupName
		resp.RoadSectionName = roadDetail.RoadSectionName
		resp.RoadName = roadDetail.RoadName
		resp.KmStart = helpers.FormatKM(int64(roadDetail.KmStart))
		resp.KmEnd = helpers.FormatKM(int64(roadDetail.KmEnd))
		resp.TotalKm = fmt.Sprintf(`%.3f`, roadDetail.TotalKm)
		resp.Year = fmt.Sprintf("%d", year+543)

		if trafficVolume, found := trafficVolumeMap[roadDetail.RoadID]; found {
			resp.Veh1 = helpers.AddCommasToNumber(fmt.Sprintf("%d", trafficVolume.Veh1))
			resp.Veh2 = helpers.AddCommasToNumber(fmt.Sprintf("%d", trafficVolume.Veh2))
			resp.Veh3 = helpers.AddCommasToNumber(fmt.Sprintf("%d", trafficVolume.Veh3))
			resp.Total = helpers.AddCommasToNumber(fmt.Sprintf("%d", trafficVolume.Total))
			if !trafficVolume.SurveyedDate.IsZero() {
				resp.SurveyedDate = helpers.TimeThai(trafficVolume.SurveyedDate)
			}

		} else {
			resp.Veh1 = fmt.Sprintf("%d", 0)
			resp.Veh2 = fmt.Sprintf("%d", 0)
			resp.Veh3 = fmt.Sprintf("%d", 0)
			resp.Total = fmt.Sprintf("%d", 0)
			resp.SurveyedDate = ""

		}

		resps = append(resps, resp)
	}

	sort.Slice(resps, func(i, j int) bool {
		return resps[i].RoadID < resps[j].RoadID
	})

	//generate report
	var pathResult interface{}

	if typ == "excel" {
		pathResult, err = helpers.ExportExcelTrafficVolume(resps)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	} else if typ == "pdf" || typ == "html" {
		pathResult, err = helpers.RequestExport(resps, "TEMPLATE_GENARAL_TRAFFIC_VOLUME", typ)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	}

	return pathResult, nil
}
