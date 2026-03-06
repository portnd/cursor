package usecases

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	"gorm.io/gorm"
)

func (u *UseCase) GetReport3(roadSectionIDstr, typ string) (interface{}, error) {

	roadSectionID, err := strconv.Atoi(roadSectionIDstr)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}

	roads, err := u.Repo.GetRoadInfoForAssetSummary(roadSectionID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		logs.Error(err)
		return responses.RoadConditionDetails{}, err
	}

	var roadIDs []int
	for _, road := range roads {
		roadIDs = append(roadIDs, road.RoadID)
	}

	var datas []models.DataReportSummaryAsset

	list, err := u.Repo.GetReportSummayAssetForAssetSummary(roadIDs)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		logs.Error(err)
		return responses.RoadConditionDetails{}, err
	}

	var tableTitle []models.TitleSummaryAsset

	if len(list) == 0 {

		for _, road := range roads {
			var data models.DataReportSummaryAsset
			copier.Copy(&datas, &road)

			data.KmStart = helpers.FormatKM(int64(road.KmStart))
			data.KmEnd = helpers.FormatKM(int64(road.KmEnd))
			data.Table = tableTitle

			datas = append(datas, data)
		}

	} else {

		// วนลูปหา road_asset_id ทั้งหมด ที่มี ref_asset_table.id เหมือนกัน ใน list แล้วปั้นเป็น obect ใหม่
		//
		//

		for _, l := range list {
			table := make([]models.ListSummaryAsset, len(l.SummaryAsset))
			for iIndex, i := range l.SummaryAsset {
				switch i.ID {
				case 158:
					{
						count, err := u.Repo.GetLightCountForAssetSummary(i.RoadAssetID)
						if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
							logs.Error(err)
							return responses.RoadConditionDetails{}, err
						}

						for _, j := range count {
							var row []string

							firstUpdatedDateTruncated := time.Date(i.FirstUpdatedDate.Year(), i.FirstUpdatedDate.Month(), i.FirstUpdatedDate.Day(), 0, 0, 0, 0, i.LastUpdatedDate.Location())
							lastUpdatedDateTruncated := time.Date(i.LastUpdatedDate.Year(), i.LastUpdatedDate.Month(), i.LastUpdatedDate.Day(), 0, 0, 0, 0, i.LastUpdatedDate.Location())

							var date string
							if firstUpdatedDateTruncated.Equal(lastUpdatedDateTruncated) {
								date = helpers.TimeThai(i.FirstUpdatedDate)
							} else {
								date = fmt.Sprintf("%s - %s", helpers.TimeThai(i.FirstUpdatedDate), helpers.TimeThai(i.LastUpdatedDate))
							}

							row = append(row, j.Name, j.Type, helpers.FormatNumber(j.Count), "ดวง", date)
							table[iIndex].Row = append(table[iIndex].Row, row)
						}

						table[iIndex].Topic = fmt.Sprint(iIndex+1) + ". " + i.TableLabel
						table[iIndex].Header = append(table[iIndex].Header, "ประเภทสินทรัพย์", "ชนิด", "จำนวน", "หน่วย", "วันที่สำรวจ")
					}

				case 297, 298, 299, 300, 304, 318, 319, 326, 315, 316, 306, 307, 308, 309, 310, 311, 312, 313, 322, 323, 325, 327, 328, 338, 342:
					{
						var unit string
						var joinName string
						switch i.ID {
						case 297, 298, 299, 300:
							unit = "ดวง"
							joinName = "ref_asset_light_type"
						case 304, 318, 319, 326:
							unit = "แห่ง"
							joinName = "ref_asset_building_type"
						case 306, 307, 308, 309, 310, 311, 312, 313:
							unit = "ป้าย"
							joinName = "ref_asset_sign_image"
						case 315, 316:
							unit = "จุด"
							joinName = "ref_asset_kmstone_type"
						case 322:
							unit = "จุด"
							joinName = "ref_asset_qutter_type"
						case 323:
							unit = "จุด"
							joinName = "ref_asset_traffic_camera_type"
						case 325:
							unit = "แห่ง"
							joinName = "ref_asset_weight_station_type"
						case 327:
							unit = "จุด"
							joinName = "ref_asset_noise_barrier"
						case 329:
							unit = "จุด"
							joinName = "ref_asset_fence_type"
						case 338:
							unit = "จุด"
							joinName = "ref_asset_crashcushion_type"
						case 342:
							unit = "จุด"
							joinName = "ref_asset_clerance_type"
						}

						count, err := u.Repo.GetTypeCountForAssetSummary(i.RoadAssetID, i.TableName, joinName)
						if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
							logs.Error(err)
							return responses.RoadConditionDetails{}, err
						}

						for _, j := range count {
							var row []string

							firstUpdatedDateTruncated := time.Date(i.FirstUpdatedDate.Year(), i.FirstUpdatedDate.Month(), i.FirstUpdatedDate.Day(), 0, 0, 0, 0, i.LastUpdatedDate.Location())
							lastUpdatedDateTruncated := time.Date(i.LastUpdatedDate.Year(), i.LastUpdatedDate.Month(), i.LastUpdatedDate.Day(), 0, 0, 0, 0, i.LastUpdatedDate.Location())

							var date string
							if firstUpdatedDateTruncated.Equal(lastUpdatedDateTruncated) {
								date = helpers.TimeThai(i.FirstUpdatedDate)
							} else {
								date = fmt.Sprintf("%s - %s", helpers.TimeThai(i.FirstUpdatedDate), helpers.TimeThai(i.LastUpdatedDate))
							}
							row = append(row, j.Name, helpers.FormatNumber(j.Count), unit, date)
							table[iIndex].Row = append(table[iIndex].Row, row)
						}

						table[iIndex].Topic = fmt.Sprint(iIndex+1) + ". " + i.TableLabel
						table[iIndex].Header = append(table[iIndex].Header, "ประเภทสินทรัพย์", "จำนวน", "หน่วย", "วันที่สำรวจ")
					}
				default:
					{
						count, err := u.Repo.GetNoTypeCountForAssetSummary(i.RoadAssetID, i.TableName)
						if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
							logs.Error(err)
							return responses.RoadConditionDetails{}, err
						}

						table[iIndex].Topic = fmt.Sprint(iIndex+1) + ". " + i.TableLabel
						table[iIndex].Header = append(table[iIndex].Header, "ประเภทสินทรัพย์", "จำนวน", "หน่วย", "วันที่สำรวจ")

						var unit string

						switch i.ID {
						case 159:
							unit = "เครื่อง"
						case 176:
							unit = "ป้าย"
						default:
							unit = "จุด"
						}
						var row []string
						firstUpdatedDateTruncated := time.Date(i.FirstUpdatedDate.Year(), i.FirstUpdatedDate.Month(), i.FirstUpdatedDate.Day(), 0, 0, 0, 0, i.LastUpdatedDate.Location())
						lastUpdatedDateTruncated := time.Date(i.LastUpdatedDate.Year(), i.LastUpdatedDate.Month(), i.LastUpdatedDate.Day(), 0, 0, 0, 0, i.LastUpdatedDate.Location())

						var date string
						if firstUpdatedDateTruncated.Equal(lastUpdatedDateTruncated) {
							date = helpers.TimeThai(i.FirstUpdatedDate)
						} else {
							date = fmt.Sprintf("%s - %s", helpers.TimeThai(i.FirstUpdatedDate), helpers.TimeThai(i.LastUpdatedDate))
						}
						row = append(row, strings.ReplaceAll(strings.ReplaceAll(i.TableLabel, "ข้อมูล", ""), "ตำแหน่ง", ""), helpers.FormatNumber(count.Count), unit, date)
						table[iIndex].Row = append(table[iIndex].Row, row)

					}
				}
			}
			tableTitle = append(tableTitle, models.TitleSummaryAsset{
				Title: l.Name,
				Table: table,
			})

			for _, road := range roads {
				var data models.DataReportSummaryAsset
				copier.Copy(&data, &road)

				data.KmStart = helpers.FormatKM(int64(road.KmStart))
				data.KmEnd = helpers.FormatKM(int64(road.KmEnd))
				data.Table = tableTitle

				datas = append(datas, data)
			}
		}
	}

	var sumData models.DataReportSummaryAsset
	if len(datas) != 0 {
		copier.Copy(&sumData, &datas[0])
	}

	sumData.Table = tableTitle
	sumData.StrRoadLength = fmt.Sprintf(`%.3f`, sumData.RoadSectionDistance)
	// generate report
	var pathResult interface{}

	if typ == "excel" {
		pathResult, err = helpers.ExportExcelType3(&sumData)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	} else {
		pathResult, err = helpers.RequestExport(sumData, "TEMPLATE_GENARAL_TYPE3", typ)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	}

	return pathResult, nil
}
