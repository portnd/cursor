package usecases

import (
	"errors"
	"fmt"
	"image"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	"gorm.io/gorm"
)

func (u *UseCase) Report11(year, roadSectionID, typ string) (interface{}, error) {
	roadSectionIDInt, err := strconv.Atoi(roadSectionID)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	roadIDs, err := u.Repo.GetRoadFromSectionIDForRoadDamage(roadSectionIDInt)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}

	roadSectionInfo, err := u.Repo.GetRoadSectionByIDForRoadDamage(roadSectionIDInt)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	roadGroup, err := u.Repo.GetRoadGroupByIDForRoadDamage(roadSectionInfo.RoadGroupId)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	report11 := models.Report11{}

	report11.RoadGroupName = roadGroup.Number
	report11.RoadSectionNumber = roadSectionInfo.Number
	report11.RoadSectionName = roadSectionInfo.NameOriginTH + " - " + roadSectionInfo.NameDestinationTH
	report11.KmStart = helpers.FormatKM(int64(roadSectionInfo.KmStart))
	report11.KmEnd = helpers.FormatKM(int64(roadSectionInfo.KmEnd))
	report11.RoadLength = helpers.FormatKM(int64(roadSectionInfo.Distance))
	report11.Year = yearInt + 543

	RoadData := []models.DataRoadDamage{}
	for _, roadID := range roadIDs {
		data, err := u.GetReportDamage(yearInt, roadID)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
		RoadData = append(RoadData, data.Detail...)
	}

	// mapRoadDamage := make(map[int]models.DataRoadDamage)

	sum := models.DataRoadDamage{}
	for _, item := range RoadData {
		// sum := mapRoadDamage[item.LaneNo]
		sum.LaneNo = item.LaneNo
		sum.ACIcrack += item.ACIcrack
		sum.ACUcrack += item.ACUcrack
		sum.ACRavelling += item.ACRavelling
		sum.ACPatching += item.ACPatching
		sum.ACPotholeArea += item.ACPotholeArea
		sum.ACBleeding += item.ACBleeding
		sum.ACPotholeCount += item.ACPotholeCount

		sum.CCTransverseCrack += item.CCTransverseCrack
		sum.CCNonTransverseCrack += item.CCNonTransverseCrack
		sum.CCSpalling += item.CCSpalling
		sum.CCCornerBreaks += item.CCCornerBreaks
		sum.CCJointSealDamage += item.CCJointSealDamage
		sum.CCPatching += item.CCPatching
		sum.CCScaling += item.CCScaling
		// mapRoadDamage[item.LaneNo] = sum
	}

	// for _, v := range mapRoadDamage {
	// 	report11.Data = append(report11.Data, v)
	// }
	report11.Data = sum
	// generate report
	var pathResult interface{}

	if typ == "excel" {
		pathResult, err = ExportExcelType11(report11)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	} else {
		pathResult, err = helpers.RequestExport(report11, "TEMPLATE_GENARAL_TYPE11_AEK", typ)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	}
	//_ = pathResult
	return pathResult, nil
}
func (u *UseCase) GetReportDamage(y, r int) (models.DataReportDamage, error) {

	data, detail, position, err := u.Repo.GetReportDamageForRoadDamage(y, r)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		logs.Error(err)
		return models.DataReportDamage{}, err
	}

	if data == nil || len(detail) == 0 || len(position) == 0 {
		roadInfo, err := u.Repo.GetRoadInfoForRoadDamage(r)
		if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
			logs.Error(err)
			return models.DataReportDamage{}, err
		}

		// calulate road distance
		length := math.Abs(float64(roadInfo.KmEnd) - float64(roadInfo.KmStart))
		roadInfo.RoadLengthStr = helpers.FormatNumberFloat(length / 1000)

		newData := models.DataReportDamage{
			RoadGroupName: roadInfo.RoadGroupName,
			RoadName:      roadInfo.RoadName,
			RoadCode:      roadInfo.RoadCode,
			StrKmStart:    helpers.FormatKM(int64(roadInfo.KmStart)),
			StrKmEnd:      helpers.FormatKM(int64(roadInfo.KmEnd)),
			RoadLengthStr: roadInfo.RoadLengthStr,
			IsNull:        true,
			Detail: []models.DataRoadDamage{
				{LaneNo: 0,

					ACIcrack:       0,
					ACUcrack:       0,
					ACRavelling:    0,
					ACPatching:     0,
					ACPotholeArea:  0,
					ACPotholeCount: 0,
					ACBleeding:     0,

					CCTransverseCrack:    0,
					CCNonTransverseCrack: 0,
					CCSpalling:           0,
					CCCornerBreaks:       0,
					CCJointSealDamage:    0,
					CCPatching:           0,
					CCScaling:            0,
					Position: []models.CovertPositionDamage{
						{No: 0,
							Km:         0,
							Surface:    "",
							DamageType: "",
							Value:      0,
							Unit:       "",
							Image:      ""},
					},
				},
			},
		}
		data = &newData
	} else {

		// calulate distance
		length := math.Abs(float64(data.KmEnd) - float64(data.KmStart))
		data.RoadLengthStr = helpers.FormatNumberFloat(length / 1000)
		data.StrKmStart = helpers.FormatKM(int64(data.KmStart))
		data.StrKmEnd = helpers.FormatKM(int64(data.KmEnd))

		//convert position and append in detail
		for index, i := range detail {
			k := 1
			for _, j := range position {
				if i.LaneNo == j.LaneNO {
					var conv models.CovertPositionDamage

					if j.ACIcrack != 0 {
						conv.DamageTypeENG = "ac icrack"
						conv.DamageType = "รอยแตกต่อเนื่อง"
						conv.Unit = "ตร.ม."
						conv.Surface = "ลาดยาง"
						conv.Value = j.ACIcrack

						if conv.Value != 0 {
							conv.StrKm = helpers.FormatKM(int64(j.Km))
							//conv.Image = helpers.ToBase64(j.ImgFilepath)

							if j.ImgFilepath == "" {
								conv.Image = ""
							} else {
								conv.Image = os.Getenv("STORAGE_IP") + "/" + j.ImgFilepath
							}

							conv.No = k
							detail[index].Position = append(detail[index].Position, conv)
							k++
						}

					}
					if j.ACUcrack != 0 {
						conv.DamageTypeENG = "ac ucrack"
						conv.DamageType = "รอยแตกไม่ต่อเนื่อง"
						conv.Unit = "ม."
						conv.Surface = "ลาดยาง"
						conv.Value = j.ACUcrack

						if conv.Value != 0 {
							conv.StrKm = helpers.FormatKM(int64(j.Km))
							//conv.Image = helpers.ToBase64(j.ImgFilepath)

							if j.ImgFilepath == "" {
								conv.Image = ""
							} else {
								conv.Image = os.Getenv("STORAGE_IP") + "/" + j.ImgFilepath
							}

							conv.No = k
							detail[index].Position = append(detail[index].Position, conv)
							k++
						}

					}
					if j.ACRavelling != 0 {
						conv.DamageTypeENG = "ac ravelling"
						conv.DamageType = "ผิวหลุดร่อน"
						conv.Unit = "ตร.ม."
						conv.Surface = "ลาดยาง"
						conv.Value = j.ACRavelling

						if conv.Value != 0 {
							conv.StrKm = helpers.FormatKM(int64(j.Km))
							//conv.Image = helpers.ToBase64(j.ImgFilepath)

							if j.ImgFilepath == "" {
								conv.Image = ""
							} else {
								conv.Image = os.Getenv("STORAGE_IP") + "/" + j.ImgFilepath
							}

							conv.No = k
							detail[index].Position = append(detail[index].Position, conv)
							k++
						}

					}
					if j.ACPatching != 0 {
						conv.DamageTypeENG = "ac patching"
						conv.DamageType = "รอยปะซ่อม"
						conv.Unit = "ตร.ม."
						conv.Surface = "ลาดยาง"
						conv.Value = j.ACPatching

						if conv.Value != 0 {
							conv.StrKm = helpers.FormatKM(int64(j.Km))
							//conv.Image = helpers.ToBase64(j.ImgFilepath)

							if j.ImgFilepath == "" {
								conv.Image = ""
							} else {
								conv.Image = os.Getenv("STORAGE_IP") + "/" + j.ImgFilepath
							}

							conv.No = k
							detail[index].Position = append(detail[index].Position, conv)
							k++
						}

					}
					if j.ACPotholeArea != 0 {
						conv.DamageTypeENG = "ac pothole_area"
						conv.DamageType = "หลุมบ่อ"
						conv.Unit = "ตร.ม."
						conv.Surface = "ลาดยาง"
						conv.Value = j.ACPotholeArea

						if conv.Value != 0 {
							conv.StrKm = helpers.FormatKM(int64(j.Km))
							//conv.Image = helpers.ToBase64(j.ImgFilepath)

							if j.ImgFilepath == "" {
								conv.Image = ""
							} else {
								conv.Image = os.Getenv("STORAGE_IP") + "/" + j.ImgFilepath
							}

							conv.No = k
							detail[index].Position = append(detail[index].Position, conv)
							k++
						}

					}
					if j.ACPotholeCount != 0 {
						conv.DamageTypeENG = "ac pothole_count"
						conv.DamageType = "หลุมบ่อ"
						conv.Unit = "หลุม"
						conv.Surface = "ลาดยาง"
						conv.Value = j.ACPotholeCount

						if conv.Value != 0 {
							conv.StrKm = helpers.FormatKM(int64(j.Km))
							//conv.Image = helpers.ToBase64(j.ImgFilepath)

							if j.ImgFilepath == "" {
								conv.Image = ""
							} else {
								conv.Image = os.Getenv("STORAGE_IP") + "/" + j.ImgFilepath
							}

							conv.No = k
							detail[index].Position = append(detail[index].Position, conv)
							k++
						}

					}
					if j.ACBleeding != 0 {
						conv.DamageTypeENG = "ac bleeding"
						conv.DamageType = "การเยิ้มของลาดยาง"
						conv.Unit = "ตร.ม."
						conv.Surface = "ลาดยาง"
						conv.Value = j.ACBleeding

						if conv.Value != 0 {
							conv.StrKm = helpers.FormatKM(int64(j.Km))
							//conv.Image = helpers.ToBase64(j.ImgFilepath)

							if j.ImgFilepath == "" {
								conv.Image = ""
							} else {
								conv.Image = os.Getenv("STORAGE_IP") + "/" + j.ImgFilepath
							}

							conv.No = k
							detail[index].Position = append(detail[index].Position, conv)
							k++
						}

					}
					if j.CCTransverseCrack != 0 {
						conv.DamageTypeENG = "cc transverse crack"
						conv.DamageType = "รอยแตกตามขวาง"
						conv.Unit = "ม."
						conv.Surface = "คอนกรีต"
						conv.Value = j.CCTransverseCrack

						if conv.Value != 0 {
							conv.StrKm = helpers.FormatKM(int64(j.Km))
							//conv.Image = helpers.ToBase64(j.ImgFilepath)

							if j.ImgFilepath == "" {
								conv.Image = ""
							} else {
								conv.Image = os.Getenv("STORAGE_IP") + "/" + j.ImgFilepath
							}

							conv.No = k
							detail[index].Position = append(detail[index].Position, conv)
							k++
						}

					}
					if j.CCNonTransverseCrack != 0 {
						conv.DamageTypeENG = "cc non transverse crack"
						conv.DamageType = "รอยแตกตามยาว และแนวทแยง"
						conv.Unit = "ม."
						conv.Surface = "คอนกรีต"
						conv.Value = j.CCNonTransverseCrack

						if conv.Value != 0 {
							conv.StrKm = helpers.FormatKM(int64(j.Km))
							//conv.Image = helpers.ToBase64(j.ImgFilepath)

							if j.ImgFilepath == "" {
								conv.Image = ""
							} else {
								conv.Image = os.Getenv("STORAGE_IP") + "/" + j.ImgFilepath
							}

							conv.No = k
							detail[index].Position = append(detail[index].Position, conv)
							k++
						}

					}
					if j.CCScaling != 0 {
						conv.DamageTypeENG = "cc scaling"
						conv.DamageType = "ผิวหลุดร่อน"
						conv.Unit = "ตร.ม"
						conv.Surface = "คอนกรีต"
						conv.Value = j.CCScaling

						if conv.Value != 0 {
							conv.StrKm = helpers.FormatKM(int64(j.Km))
							//conv.Image = helpers.ToBase64(j.ImgFilepath)

							if j.ImgFilepath == "" {
								conv.Image = ""
							} else {
								conv.Image = os.Getenv("STORAGE_IP") + "/" + j.ImgFilepath
							}

							conv.No = k
							detail[index].Position = append(detail[index].Position, conv)
							k++
						}

					}
					if j.CCSpalling != 0 {
						conv.DamageTypeENG = "cc spalling"
						conv.DamageType = "รอยบิ่นกะเทาะ"
						conv.Unit = "ตร.ม"
						conv.Surface = "คอนกรีต"
						conv.Value = j.CCSpalling

						if conv.Value != 0 {
							conv.StrKm = helpers.FormatKM(int64(j.Km))
							//conv.Image = helpers.ToBase64(j.ImgFilepath)

							if j.ImgFilepath == "" {
								conv.Image = ""
							} else {
								conv.Image = os.Getenv("STORAGE_IP") + "/" + j.ImgFilepath
							}

							conv.No = k
							detail[index].Position = append(detail[index].Position, conv)
							k++
						}

					}
					if j.CCCornerBreaks != 0 {
						conv.DamageTypeENG = "cc cornerbreaks"
						conv.DamageType = "รอยแตกที่มุม"
						conv.Unit = "จุด"
						conv.Surface = "คอนกรีต"
						conv.Value = j.CCCornerBreaks

						if conv.Value != 0 {
							conv.StrKm = helpers.FormatKM(int64(j.Km))
							//conv.Image = helpers.ToBase64(j.ImgFilepath)

							if j.ImgFilepath == "" {
								conv.Image = ""
							} else {
								conv.Image = os.Getenv("STORAGE_IP") + "/" + j.ImgFilepath
							}

							conv.No = k
							detail[index].Position = append(detail[index].Position, conv)
							k++
						}

					}
					if j.CCJointSealDamage != 0 {
						conv.DamageTypeENG = "cc joint seal damage"
						conv.DamageType = "วัสดุยาแนวรอยต่อเสียหาย"
						conv.Unit = "ม."
						conv.Surface = "คอนกรีต"
						conv.Value = j.CCJointSealDamage

						if conv.Value != 0 {
							conv.StrKm = helpers.FormatKM(int64(j.Km))
							//conv.Image = helpers.ToBase64(j.ImgFilepath)

							if j.ImgFilepath == "" {
								conv.Image = ""
							} else {
								conv.Image = os.Getenv("STORAGE_IP") + "/" + j.ImgFilepath
							}

							conv.No = k
							detail[index].Position = append(detail[index].Position, conv)
							k++
						}

					}
					if j.CCPatching != 0 {
						conv.DamageTypeENG = "cc patching"
						conv.DamageType = "รอยปะซ่อม"
						conv.Unit = "ตร.ม."
						conv.Surface = "คอนกรีต"
						conv.Value = j.CCPatching

						if conv.Value != 0 {
							conv.StrKm = helpers.FormatKM(int64(j.Km))
							//conv.Image = helpers.ToBase64(j.ImgFilepath)

							if j.ImgFilepath == "" {
								conv.Image = ""
							} else {
								conv.Image = os.Getenv("STORAGE_IP") + "/" + j.ImgFilepath
							}

							conv.No = k
							detail[index].Position = append(detail[index].Position, conv)
							k++
						}
					}
				}
			}
		}
		// add detail in data
		data.Detail = detail
	}

	return *data, nil
}

func ExportExcelType11(data models.Report11) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE11_AEK_EXCEL"))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	logoPath := os.Getenv("LOGO")

	fileLogo, err := os.Open(fmt.Sprint(logoPath))
	if err != nil {
		return nil, err
	}
	defer fileLogo.Close()

	imageLogoConfig, _, err := image.DecodeConfig(fileLogo)
	if err != nil {
		return nil, err
	}

	err = f.AddPicture("Sheet1", "B1", fmt.Sprint(logoPath), &excelize.GraphicOptions{ScaleX: 140 / float64(imageLogoConfig.Width), ScaleY: 120 / float64(imageLogoConfig.Height)})
	if err != nil {
		fmt.Println("Error adding picture:", err)
		return nil, err
	}

	f.SetCellValue("Sheet1", "C3", fmt.Sprintf("หมายเลขทางหลวง : %s ตอนควบคุม : %s ", data.RoadGroupName, data.RoadSectionNumber))
	f.SetCellValue("Sheet1", "C4", "ชื่อสายทาง : "+data.RoadSectionName)
	f.SetCellValue("Sheet1", "C5", "กม.เริ่มต้น "+data.KmStart+" กม.สิ้นสุด "+data.KmEnd+" ระยะทาง "+data.RoadLength+" กม.")

	if !data.IsNull {
		iRow := 7
		// iStrRow := fmt.Sprint(iRow)
		// f.SetCellValue("Sheet1", "B"+iStrRow, "ช่องจราจรที่ "+fmt.Sprint(i.LaneNo))

		f.SetCellValue("Sheet1", "F"+fmt.Sprint(iRow+3), data.Data.ACIcrack)          //รอยแตกต่อเนื่อง
		f.SetCellValue("Sheet1", "M"+fmt.Sprint(iRow+3), data.Data.CCTransverseCrack) //แตกตามขวาง

		f.SetCellValue("Sheet1", "F"+fmt.Sprint(iRow+4), data.Data.ACUcrack)             //รอบแตกไม่ต่อเนื่อง
		f.SetCellValue("Sheet1", "M"+fmt.Sprint(iRow+4), data.Data.CCNonTransverseCrack) //แนวทแยง

		f.SetCellValue("Sheet1", "F"+fmt.Sprint(iRow+5), data.Data.ACRavelling) //ผิวหลุด
		f.SetCellValue("Sheet1", "M"+fmt.Sprint(iRow+5), data.Data.CCCornerBreaks)

		f.SetCellValue("Sheet1", "F"+fmt.Sprint(iRow+6), data.Data.ACPatching) //รอยปะ
		f.SetCellValue("Sheet1", "M"+fmt.Sprint(iRow+6), data.Data.CCJointSealDamage)

		f.SetCellValue("Sheet1", "F"+fmt.Sprint(iRow+7), data.Data.ACPotholeArea) //หลุมบ่อ
		f.SetCellValue("Sheet1", "M"+fmt.Sprint(iRow+7), data.Data.CCPatching)    //รอยแตกตามมุม

		f.SetCellValue("Sheet1", "F"+fmt.Sprint(iRow+8), data.Data.ACBleeding) //เสียรูป
		f.SetCellValue("Sheet1", "M"+fmt.Sprint(iRow+8), data.Data.CCSpalling) //รอยปะ

		f.SetCellValue("Sheet1", "F"+fmt.Sprint(iRow+9), data.Data.ACPotholeCount) //เยิ้ม
		f.SetCellValue("Sheet1", "M"+fmt.Sprint(iRow+9), data.Data.CCScaling)      //ผิวหลุด

	}

	helpers.AddFooter(f, "Sheet1")

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	reportName, err := helpers.ReportName("TEMPLATE_GENARAL_TYPE11_AEK_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}

func (u *UseCase) Report10(year, roadSectionID, typ string) (interface{}, error) {
	roadSectionIDInt, err := strconv.Atoi(roadSectionID)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	roadIDs, err := u.Repo.GetRoadFromSectionIDForRoadDamage(roadSectionIDInt)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}

	roadSectionInfo, err := u.Repo.GetRoadSectionByIDForRoadDamage(roadSectionIDInt)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	roadGroup, err := u.Repo.GetRoadGroupByIDForRoadDamage(roadSectionInfo.RoadGroupId)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return nil, responses.NewAppErr(400, err.Error())
	}
	report10 := models.Report10{}

	report10.RoadGroupName = roadGroup.Number
	report10.RoadSectionNumber = roadSectionInfo.Number
	report10.RoadSectionName = roadSectionInfo.NameOriginTH + " - " + roadSectionInfo.NameDestinationTH
	report10.KmStart = helpers.FormatKM(int64(roadSectionInfo.KmStart))
	report10.KmEnd = helpers.FormatKM(int64(roadSectionInfo.KmEnd))
	report10.RoadLength = helpers.FormatNumberFloat(float64(roadSectionInfo.Distance))
	report10.Year = yearInt + 543

	RoadData := []models.DataReportDamage{}
	for _, roadID := range roadIDs {
		data, err := u.GetReportDamage(yearInt, roadID)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
		if !data.IsNull {
			RoadData = append(RoadData, data)
		}
	}

	// mapRoadDamage := make(map[int]models.DataRoadDamage)

	// for _, item := range RoadData {
	// 	sum := mapRoadDamage[item.LaneNo]
	// 	sum.LaneNo = item.LaneNo
	// 	sum.ACIcrack += item.ACIcrack
	// 	sum.ACUcrack += item.ACUcrack
	// 	sum.ACRavelling += item.ACRavelling
	// 	sum.ACPatching += item.ACPatching
	// 	sum.ACPotholeArea += item.ACPotholeArea
	// 	sum.ACBleeding += item.ACBleeding
	// 	sum.ACPotholeCount += item.ACPotholeCount

	// 	sum.CCTransverseCrack += item.CCTransverseCrack
	// 	sum.CCNonTransverseCrack += item.CCNonTransverseCrack
	// 	sum.CCSpalling += item.CCSpalling
	// 	sum.CCCornerBreaks += item.CCCornerBreaks
	// 	sum.CCJointSealDamage += item.CCJointSealDamage
	// 	sum.CCPatching += item.CCPatching
	// 	sum.CCScaling += item.CCScaling
	// 	mapRoadDamage[item.LaneNo] = sum
	// }

	// for _, v := range mapRoadDamage {
	// 	report10.Data = append(report10.Data, v)
	// }
	report10.Data = RoadData
	// generate report
	var pathResult interface{}

	if typ == "excel" {
		pathResult, err = ExportExcelType10(report10)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	} else {
		pathResult, err = helpers.RequestExport(report10, "TEMPLATE_GENARAL_TYPE10_AEK", typ)
		if err != nil {
			return nil, responses.NewAppErr(400, err.Error())
		}
	}
	//_ = pathResult
	return pathResult, nil
}

func ExportExcelType10(data models.Report10) (interface{}, error) {
	filePath := os.Getenv("GENARAL_EXCEL")
	f, err := excelize.OpenFile(os.Getenv("TEMPLATE_GENARAL_TYPE10_AEK_EXCEL"))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	folderPath := "storages/road/temp" + uuid.New().String() + "/"
	// create temp folder
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	// remove temp folder
	defer os.RemoveAll(folderPath)

	templateIndex, _ := f.GetSheetIndex("Sheet1")
	sheetList := make(map[string]int)

	// Initialize a map to keep track of sheets that already have the picture
	for _, v := range data.Data {
		logoPath := os.Getenv("LOGO")

		fileLogo, err := os.Open(fmt.Sprint(logoPath))
		if err != nil {
			return nil, err
		}
		defer fileLogo.Close()

		imageLogoConfig, _, err := image.DecodeConfig(fileLogo)
		if err != nil {
			return nil, err
		}

		// Create a new sheet by copying the template sheet
		sheetName := v.RoadName
		if utf8.RuneCountInString(sheetName) > 31 {
			runes := []rune(sheetName)
			sheetName = string(runes[:25])
		}
		if count, exists := sheetList[sheetName]; exists {
			sheetList[sheetName]++
			if len(sheetName) > 75 {
				sheetName = sheetName[:75]
			}
			sheetName = fmt.Sprintf("%s (%d)", sheetName, count)
			newSheetIndex, err := f.NewSheet(sheetName)
			if err != nil {
				fmt.Println("Error creating new sheet:", err)
				return nil, err
			}
			log := fmt.Sprintf("template index : %d new index : %d", templateIndex, newSheetIndex)
			fmt.Println(log)
			err = f.CopySheet(templateIndex, newSheetIndex)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		} else {
			newSheetIndex, err := f.NewSheet(sheetName)
			if err != nil {
				fmt.Println("Error creating new sheet:", err)
				return nil, err
			}

			err = f.CopySheet(templateIndex, newSheetIndex)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			sheetList[sheetName]++

		}

		var wg sync.WaitGroup

		for _, i := range v.Detail {
			for _, j := range i.Position {
				wg.Add(1) // add number of go routine to wait
				go func(imgURL string) {
					defer wg.Done() // decrease the number of go routine to wait after go routine finished

					// trim img path
					trimmedURL := strings.TrimPrefix(imgURL, os.Getenv("STORAGE_IP")+"/")

					// open file
					srcImage, err := imaging.Open(trimmedURL)
					if err != nil {
						fmt.Println("can not open file", err)
						return
					}

					fileName := filepath.Base(trimmedURL)
					//saveAs := filepath.Join(folderPath, fileName)
					saveAs := folderPath + fileName
					// convert to .png
					err = imaging.Save(srcImage, saveAs)
					if err != nil {
						fmt.Println("can not save file", err)
						return
					}
				}(j.Image)
			}
		}

		wg.Wait() // waiting for go routine finish

		f.SetCellValue(sheetName, "C3", fmt.Sprintf("หมายเลขทางหลวง : %s ตอนควบคุม : %s ", data.RoadGroupName, data.RoadSectionNumber))
		f.SetCellValue(sheetName, "C4", "ชื่อสายทาง : "+v.RoadName)
		f.SetCellValue(sheetName, "C5", "กม.เริ่มต้น "+v.StrKmStart+" กม.สิ้นสุด "+v.StrKmEnd+" ระยะทาง "+v.RoadLengthStr+" กม.")
		f.SetCellValue(sheetName, "B6", "ชื่อสายทาง : "+v.RoadName)

		if !data.IsNull {
			iRow := 7
			iStrRow := fmt.Sprint(iRow)
			for iIndex, i := range v.Detail {
				f.SetCellValue(sheetName, "B"+iStrRow, "ช่องจราจรที่ "+fmt.Sprint(i.LaneNo))
				jRow := iRow + 3
				jStrRow := fmt.Sprint(jRow)

				for jIndex, j := range i.Position {
					f.SetCellValue(sheetName, "B"+jStrRow, j.No)
					f.SetCellValue(sheetName, "C"+jStrRow, j.StrKm)
					f.SetCellValue(sheetName, "E"+jStrRow, j.Surface)
					f.SetCellValue(sheetName, "G"+jStrRow, j.DamageTypeENG+"\n"+j.DamageType)
					f.SetCellValue(sheetName, "J"+jStrRow, j.Value)
					f.SetCellValue(sheetName, "K"+jStrRow, j.Unit)
					// trim img path
					fileName := filepath.Base(j.Image)

					picRow := fmt.Sprint(jRow + 1)
					f.AddPicture(sheetName, "L"+picRow, folderPath+fileName, &excelize.GraphicOptions{OffsetX: 8, OffsetY: -119, ScaleX: 0.078, ScaleY: 0.075})
					f.SetCellHyperLink(sheetName, "L"+jStrRow, j.Image, "External")

					if jIndex < len(i.Position)-1 {
						f.DuplicateRow(sheetName, jRow)
						jRow++
						jStrRow = fmt.Sprint(jRow)
					}
				}

				if iIndex < len(v.Detail)-1 {

					for count := 1; count <= 14; count++ {
						f.DuplicateRowTo(sheetName, 6+count, jRow+1+count)
					}
					iRow = jRow + 2
					iStrRow = fmt.Sprint(iRow)
				}
			}
		}

		helpers.AddFooter(f, sheetName)
		// Define the target dimensions in inches
		targetWidthInches := 1.55
		targetHeightInches := 1.0

		// Convert inches to pixels (assuming 96 DPI)
		dpi := 96.0
		targetWidthPixels := targetWidthInches * dpi
		targetHeightPixels := targetHeightInches * dpi

		// Calculate the scaling factors based on the image's original dimensions
		scaleX := targetWidthPixels / float64(imageLogoConfig.Width)
		scaleY := targetHeightPixels / float64(imageLogoConfig.Height)

		// Add the picture with the calculated scaling factors
		err = f.AddPicture(sheetName, "B2", fmt.Sprint(logoPath), &excelize.GraphicOptions{
			ScaleX: scaleX,
			ScaleY: scaleY,
		})
		if err != nil {
			fmt.Println("Error adding picture:", err)
			return nil, err
		}
	}
	f.DeleteSheet("Sheet1")
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return nil, err
	}

	reportName, err := helpers.ReportName("TEMPLATE_GENARAL_TYPE10_AEK_EXCEL")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	name := code + "_" + reportName

	f.SaveAs(os.Getenv("GENARAL_EXCEL") + name + ".xlsx")

	return os.Getenv("STORAGE_IP") + "/" + filePath + name + ".xlsx", nil
}
