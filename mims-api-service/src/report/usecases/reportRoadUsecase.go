package usecases

import (
	"errors"
	"math"
	"strconv"

	"github.com/jinzhu/copier"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	"gorm.io/gorm"
)

func (u *UseCase) GetReportRoad(roadGroupIDs []int, typ string) (interface{}, error) {

	if roadGroupIDs == nil {
		return nil, responses.NewAppErr(400, "กรุณาเลือกหมายเลขทางหลวง")
	}

	if typ != "excel" && typ != "pdf" && typ != "html" {
		if typ == "" {
			return nil, responses.NewAppErr(400, "กรุณาเลือกประเภทรายงาน")
		}
		return nil, responses.NewAppErr(400, "ประเภทรายงานไม่ถูกต้อง")
	}

	resp, err := u.Repo.GetReportRoad(roadGroupIDs)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		logs.Error(err)
		return responses.RoadConditionDetails{}, err
	}
	var roadGroups []models.RoadListReport
	for _, item := range resp {
		var roadGroup models.RoadListReport
		var roadSectionDataReports []models.RoadSectionDataReport
		for key, section := range item.Sections {
			var roadSectionDataReport models.RoadSectionDataReport
			err = copier.Copy(&roadSectionDataReport, &section)
			if err != nil {
				logs.Error(err)
				return nil, err
			}
			kmStart, err := strconv.ParseInt(roadSectionDataReport.KmStart, 10, 64)
			if err != nil {
				logs.Error(err)
				return nil, err
			}
			kmEnd, err := strconv.ParseInt(roadSectionDataReport.KmEnd, 10, 64)
			if err != nil {
				logs.Error(err)
				return nil, err
			}
			roadSectionDataReport.No = key + 1
			roadSectionDataReport.KmStart = helpers.FormatKM(kmStart)
			roadSectionDataReport.KmEnd = helpers.FormatKM(kmEnd)
			roadSectionDataReport.StrDistance = helpers.FormatNumberFloat(float64(roadSectionDataReport.Distance))
			roadSectionDataReports = append(roadSectionDataReports, roadSectionDataReport)
		}
		err = copier.Copy(&roadGroup, &item)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		kmStart, err := strconv.ParseInt(roadGroup.KmStart, 10, 64)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		kmEnd, err := strconv.ParseInt(roadGroup.KmEnd, 10, 64)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		roadGroup.Sections = roadSectionDataReports
		roadGroup.KmStart = helpers.FormatKM(kmStart)
		roadGroup.KmEnd = helpers.FormatKM(kmEnd)
		roadGroup.StrDistance = helpers.FormatNumberFloat(float64(math.Abs(float64(kmStart)-float64(kmEnd))) / 1000)
		roadGroups = append(roadGroups, roadGroup)
	}
	// return roadGroups, nil

	// generate report
	var pathResult interface{}

	if typ == "excel" {
		pathResult, err = helpers.ExportExcelRoad(roadGroups)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	} else if typ == "pdf" || typ == "html" {
		pathResult, err = helpers.RequestExport(roadGroups, "TEMPLATE_GENARAL_ROAD", typ)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	}

	return pathResult, nil
}
