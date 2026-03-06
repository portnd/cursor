package repositories

import (
	"fmt"
	"strings"
	"time"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/report/handlers"
	"gitlab.com/mims-api-service/src/report/usecases"

	_ "github.com/go-sql-driver/mysql"
	servicesDB "gitlab.com/mims-api-service/services/database"
	"gorm.io/gorm"
)

type Repository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewRepositoryHandler(conn *gorm.DB) *handlers.Handler {
	servicesDB := servicesDB.NewServicesDatabase(conn)
	useCase := usecases.NewUseCase(&Repository{conn}, servicesDB)
	handler := handlers.NewHandler(useCase)
	return handler
}

// ////////////// NEW MIMS ////////////////
func (r *Repository) CreateReportStatus() (int, error) {
	var reportStatus models.ReportStatus
	dateTime := time.Now()
	reportStatus.IsFinish = false
	reportStatus.CreatedAt = dateTime

	query := r.conn
	err := query.Create(&reportStatus).Error
	if err != nil {
		return 0, err
	}
	return reportStatus.Id, err
}

func (r *Repository) UpdateReportStatus(id int, path string) error {
	var reportStatus models.ReportStatus
	query := r.conn
	err := query.Where("id = ?", id).Find(&reportStatus).Error
	if err != nil {
		return err
	}

	dateTime := time.Now()
	reportStatus.Id = id
	reportStatus.IsFinish = true
	reportStatus.Path = path
	reportStatus.UpdatedAt = dateTime

	err = query.Save(&reportStatus).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CheckReportStatusById(id int) (models.ReportStatus, error) {
	var reportStatus models.ReportStatus

	query := r.conn
	err := query.Where("id = ?", id).Find(&reportStatus).Error
	if err != nil {
		return reportStatus, err
	}
	return reportStatus, err
}

func (r *Repository) FilterAssetRoad(isAllData, isOwnerData bool, depotCode string) ([]models.FilterAssetRoad, error) {
	var filterAssetRoad []models.FilterAssetRoad
	query := r.conn
	query = query.Select("road_section.*").
		Joins("JOIN road on road.road_section_id = road_section.id and road.is_active = true").
		Group("road_section.id")
	query = query.Preload("RoadGroup")
	query = query.Preload("RefDivision")
	query = query.Preload("RefDistrict")
	query = query.Preload("RefDepot")
	//ref_depot
	if isOwnerData && !isAllData {
		query = query.Where("ref_depot_code = ?", depotCode)
	}

	if !isOwnerData && !isAllData {
		query = query.Where("ref_depot_code  = ?", "no_data")
	}

	if !isOwnerData && !isAllData {
		return []models.FilterAssetRoad{}, nil
	}
	err := query.Find(&filterAssetRoad).Error
	if err != nil {
		return filterAssetRoad, err
	}
	return filterAssetRoad, nil
}

func (r *Repository) FilterAsset() ([]models.FilterAsset, error) {
	var filterForAsset []models.FilterAsset
	query := r.conn
	query = query.Preload("RefAssetTable", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = true")
	})

	query = query.Where("status = 1")

	err := query.Find(&filterForAsset).Error
	if err != nil {
		return filterForAsset, err
	}
	return filterForAsset, nil
}

func (r *Repository) FilterRoadGroup(isAllData, isOwnerData bool, depotCode string) ([]models.RoadGroup, error) {
	var roadGroup []models.RoadGroup
	query := r.conn
	query = query.Select("road_group.ID, road_group.name, road_group.number, road_group.short_name")
	query = query.Table("road_group")
	query = query.Joins("JOIN road_section ON road_group.ID = road_section.road_group_id").
		Joins("JOIN road ON road.road_group_id = road_group.id and road.is_active = true")
	query = query.Group("road_group.name,road_group.id,road_group.number, road_group.short_name")

	if isOwnerData && !isAllData {
		query = query.Where("road_section.ref_depot_code  = ?", depotCode)
	}

	if !isOwnerData && !isAllData {
		query = query.Where("road_section.ref_depot_code  = ?", "no_data")
	}

	if !isOwnerData && !isAllData {
		return []models.RoadGroup{}, nil
	}
	err := query.Find(&roadGroup).Error
	if err != nil {
		return roadGroup, err
	}
	return roadGroup, nil
}

func (r *Repository) FilterRoadSurface(isAllData, isOwnerData bool, depotCode string) ([]models.FilterRoadSurface, error) {
	var roadSurface []models.FilterRoadSurface
	query := r.conn

	query = query.Preload("Road", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = true")
	})

	query = query.Preload("Road.RoadSection")
	// query = query.Preload("Road.RoadSection.RefDepot, ")
	query = query.Preload("Road.RoadSection.RefDepot", func(db *gorm.DB) *gorm.DB {
		if isOwnerData && !isAllData {
			db = db.Where("depot_code  = ?", depotCode)
		}
		if !isOwnerData && !isAllData {
			db = db.Where("depot_code  = ?", "no_data")
		}
		return db
	})
	query = query.Preload("Road.RoadSection.RoadGroup")
	query = query.Where("status = 'A'")

	err := query.Find(&roadSurface).Error
	if err != nil {
		return roadSurface, err
	}
	return roadSurface, nil
}

func (r *Repository) FilterRoadCondition(isAllData, isOwnerData bool, depotCode string) ([]models.FilterRoadCondition, error) {
	var roadCondition []models.FilterRoadCondition
	query := r.conn

	query = query.Preload("Road", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = true")
	})

	// query = query.Preload("Road.RoadSection")

	query = query.Preload("Road.RoadSection", func(db *gorm.DB) *gorm.DB {
		if isOwnerData && !isAllData {
			db = db.Where("ref_depot_code  = ?", depotCode)
		}
		if !isOwnerData && !isAllData {
			db = db.Where("ref_depot_code  = ?", "no_data")
		}
		return db
	})
	query = query.Preload("Road.RoadSection.RefDepot")
	query = query.Preload("Road.RoadSection.RoadGroup")

	query = query.Where("status = 'A'")

	err := query.Find(&roadCondition).Error
	if err != nil {
		return roadCondition, err
	}
	return roadCondition, nil
}

func (r *Repository) FilterRefOwner() ([]models.RefOwner, error) {
	var refOwner []models.RefOwner
	query := r.conn

	query = query.Where("is_active = true")
	query = query.Order("id ASC")

	err := query.Find(&refOwner).Error
	if err != nil {
		return refOwner, err
	}
	return refOwner, nil
}

func (r *Repository) FilterRoadDamage(isAllData, isOwnerData bool, depotCode string) ([]models.FilterRoadDamage, error) {
	var filterRoadDamage []models.FilterRoadDamage
	query := r.conn

	query = query.Preload("Road", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = true")
	})

	query = query.Preload("Road.RoadSection")

	// query = query.Preload("Road.RoadSection.RefDepot")
	query = query.Preload("Road.RoadSection.RefDepot", func(db *gorm.DB) *gorm.DB {
		if isOwnerData && !isAllData {
			db = db.Where("depot_code  = ?", depotCode)
		}
		if !isOwnerData && !isAllData {
			db = db.Where("depot_code  = ?", "no_data")
		}
		return db
	})

	query = query.Preload("Road.RoadSection.RoadGroup")

	query = query.Where("status = 'A'")

	err := query.Find(&filterRoadDamage).Error
	if err != nil {
		return filterRoadDamage, err
	}
	return filterRoadDamage, nil
}

func (r *Repository) FilterRefOwnerLine() ([]models.RefOwnerRoadLine, error) {
	var refOwnerRoadLine []models.RefOwnerRoadLine

	query := r.conn

	query = query.Where("is_active = true")
	query = query.Order("id ASC")

	err := query.Find(&refOwnerRoadLine).Error

	if err != nil {
		return refOwnerRoadLine, err
	}
	return refOwnerRoadLine, nil
}

func (r *Repository) FilterRoadRetroReflectivity(isAllData, isOwnerData bool, depotCode string) ([]models.FilterRoadRetroReflectivity, error) {
	var roadRetroReflectivity []models.FilterRoadRetroReflectivity
	query := r.conn

	query = query.Preload("Road", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = true")
	})

	query = query.Preload("Road.RoadSection")

	// query = query.Preload("Road.RoadSection.RefDepot")
	query = query.Preload("Road.RoadSection.RefDepot", func(db *gorm.DB) *gorm.DB {
		if isOwnerData && !isAllData {
			db = db.Where("depot_code  = ?", depotCode)
		}
		if !isOwnerData && !isAllData {
			db = db.Where("depot_code  = ?", "no_data")
		}
		return db
	})

	query = query.Preload("Road.RoadSection.RoadGroup")

	query = query.Where("status = 'A'")

	err := query.Find(&roadRetroReflectivity).Error
	if err != nil {
		return roadRetroReflectivity, err
	}
	return roadRetroReflectivity, nil
}

func (r *Repository) FilterMaintenance(isAllData, isOwnerData bool, depotCode string) ([]models.FilterMaintenance, error) {
	var filterMaintenance []models.FilterMaintenance
	query := r.conn

	query = query.Preload("MaintenanceRoad", "status = 'A'")

	query = query.Preload("MaintenanceRoad.Road")

	query = query.Preload("MaintenanceRoad.Road.RoadSection", func(db *gorm.DB) *gorm.DB {
		if isOwnerData && !isAllData {
			db = db.Where("ref_depot_code  = ?", depotCode)
		}
		if !isOwnerData && !isAllData {
			db = db.Where("ref_depot_code  = ?", "no_data")
		}
		return db
	})

	// query = query.Preload("MaintenanceRoad.Road.RoadGroup.RoadSection")
	query = query.Preload("MaintenanceRoad.Road.RoadSection.RoadGroup")

	query = query.Where("status = 'A'")

	err := query.Find(&filterMaintenance).Error
	if err != nil {
		return filterMaintenance, err
	}
	return filterMaintenance, nil
}

func (r *Repository) FilterAadt(isAllData, isOwnerData bool, depotCode string) ([]models.FilterAadt, error) {
	var filterAadt []models.FilterAadt
	query := r.conn

	query = query.Preload("Road", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = true")
	})

	query = query.Preload("Road.RoadSection", func(db *gorm.DB) *gorm.DB {
		if isOwnerData && !isAllData {
			db = db.Where("ref_depot_code  = ?", depotCode)
		}
		if !isOwnerData && !isAllData {
			db = db.Where("ref_depot_code  = ?", "no_data")
		}
		return db
	})

	query = query.Preload("Road.RoadSection.RoadGroup")
	query = query.Where("status = 'A'")

	err := query.Find(&filterAadt).Error
	if err != nil {
		return filterAadt, err
	}
	return filterAadt, nil
}

func (r *Repository) ReportRetroReflectivity(roadSectionId, filterCriteriaId, year int) (models.ReportRetroReflectivityRoadSection, error) {

	var reportRetroReflectivityRoadSection models.ReportRetroReflectivityRoadSection

	query := r.conn

	var roadInfo []models.RoadInfo
	err := query.Where("status = 'A'").Find(&roadInfo).Error
	if err != nil {
		return reportRetroReflectivityRoadSection, err
	}

	roadInfoMap := map[int]models.RoadInfo{}
	for _, item := range roadInfo {
		roadInfoMap[item.RoadId] = item
	}

	var refOwnerRoadLine models.RefOwnerRoadLine
	err = query.Where("is_active = true").Where("id = ?", filterCriteriaId).Find(&refOwnerRoadLine).Error
	if err != nil {
		return reportRetroReflectivityRoadSection, err
	}

	var roadSection models.RoadSection
	err = query.Where("id = ?", roadSectionId).Find(&roadSection).Error
	if err != nil {
		return reportRetroReflectivityRoadSection, err
	}

	reportRetroReflectivityRoadSection.RoadSection = roadSection

	var roadRoadGruop models.RoadGroup
	err = query.Where("id = ?", roadSection.RoadGroupId).Find(&roadRoadGruop).Error
	if err != nil {
		return reportRetroReflectivityRoadSection, err
	}

	reportRetroReflectivityRoadSection.RoadGroup = roadRoadGruop

	roadIds := []int{}
	err = query.Table("road").Select("id").Where("road_section_id = ?", roadSectionId).Find(&roadIds).Error
	if err != nil {
		return reportRetroReflectivityRoadSection, err
	}

	var retroReflectivity []models.RoadRetroReflectivity
	err = query.Where("road_id IN (?) AND status = 'A' AND year = ?", roadIds, year).Find(&retroReflectivity).Error
	if err != nil {
		return reportRetroReflectivityRoadSection, err
	}

	var reportRetroReflectivitys []models.ReportRetroReflectivity
	for _, item := range retroReflectivity {

		var reportRetroReflectivity models.ReportRetroReflectivity
		var roadRetroReflectivityAllMeter []models.RoadRetroReflectivityAllMeter

		if refOwnerRoadLine.Name == "25 เมตร" {
			var rangIDs []int

			err = query.Select("id").Table("road_retro_reflectivity_range").Where("road_retro_reflectivity_id = ? ", item.ID).Find(&rangIDs).Error
			if err != nil {
				return reportRetroReflectivityRoadSection, err
			}

			err = query.Select("*, st_astext(the_geom) as the_geom").Table("road_retro_reflectivity_m").Where("road_retro_reflectivity_range_id in ? ", rangIDs).Order("km_start").Find(&roadRetroReflectivityAllMeter).Error
			if err != nil {
				return reportRetroReflectivityRoadSection, err
			}
		} else {
			err = query.Select("*, st_astext(the_geom) as the_geom").Table("road_retro_reflectivity_range").Where("road_retro_reflectivity_id = ? ", item.ID).Order("km_start").Find(&roadRetroReflectivityAllMeter).Error
			if err != nil {
				return reportRetroReflectivityRoadSection, err
			}
		}

		reportRetroReflectivity.RoadRetroReflectivity = item
		reportRetroReflectivity.RoadInfo = roadInfoMap[item.RoadID]
		reportRetroReflectivity.RoadRetroReflectivityRange = append(reportRetroReflectivity.RoadRetroReflectivityRange, roadRetroReflectivityAllMeter...)

		reportRetroReflectivitys = append(reportRetroReflectivitys, reportRetroReflectivity)
	}

	reportRetroReflectivityRoadSection.RoadRetroReflectivity = reportRetroReflectivitys

	return reportRetroReflectivityRoadSection, nil
}

func (r *Repository) ParamsRoadLine(filterCriteriaId int) ([]models.ParamsRoadLinePreload, error) {
	query := r.conn

	var paramsRoadLinePreload []models.ParamsRoadLinePreload
	err := query.Where("ref_owner_road_line_id = ?", filterCriteriaId).Preload("RefOwnerRoadLine").Preload("RefGrade").Find(&paramsRoadLinePreload).Error
	if err != nil {
		return paramsRoadLinePreload, err
	}

	return paramsRoadLinePreload, nil
}

func (r *Repository) RefStripeColor() ([]models.RefStripeColor, error) {
	query := r.conn

	var refStripeColor []models.RefStripeColor
	err := query.Find(&refStripeColor).Error
	if err != nil {
		return refStripeColor, err
	}

	return refStripeColor, nil
}

func (r *Repository) RefStripeType() ([]models.RefStripeType, error) {
	query := r.conn

	var refStripeType []models.RefStripeType
	err := query.Find(&refStripeType).Error
	if err != nil {
		return refStripeType, err
	}

	return refStripeType, nil
}

func (r *Repository) RefGrade() ([]models.RefGrade, error) {
	query := r.conn

	var refGrade []models.RefGrade
	err := query.Find(&refGrade).Error
	if err != nil {
		return refGrade, err
	}

	return refGrade, nil
}

func (r *Repository) Report12RoadGroup(roadSectionId int) (models.ReportKpiRoadGroup, error) {

	var reportKpi models.ReportKpiRoadGroup

	query := r.conn

	var roadSection models.RoadSection
	err := query.Where("id = ?", roadSectionId).Find(&roadSection).Error
	if err != nil {
		return reportKpi, err
	}

	reportKpi.RoadSection = roadSection

	var roadRoadGruop models.RoadGroup
	err = query.Where("id = ?", roadSection.RoadGroupId).Find(&roadRoadGruop).Error
	if err != nil {
		return reportKpi, err
	}

	reportKpi.RoadGroup = roadRoadGruop

	return reportKpi, nil
}

func (r *Repository) Report12Iri(roadSectionId, year int) (models.ResultIri, error) {

	query := r.conn

	var conditions1000 models.ParamCondition
	var iri100 models.ParamCondition

	var tableResultIri1000 []models.ResultIri1000
	var tableResultIri100 []models.ResultIri100

	query.
		Table("params_condition").
		Where("ref_owner_id = ?", 92).
		First(&conditions1000)
	iriCondition1000Ac := conditions1000.RightValueAc
	iriCondition1000Cc := conditions1000.RightValueCc

	query.
		Table("params_condition").
		Where("ref_owner_id = ?", 91).
		Where("Condition_type = ?", "IRI").
		First(&iri100)
	iriCondition100Ac := iri100.RightValueAc
	iriCondition100Cc := iri100.RightValueCc

	var roadSurveyIri100s []models.RoadSurveyIri100
	query.
		Table("road_condition_survey_100m").
		Select("km_start,km_end,iri,survey_type, road_condition_survey_id").
		Find(&roadSurveyIri100s)

	roadSurveyIri100Map := map[int][]models.RoadSurveyIri100{}
	for _, item := range roadSurveyIri100s {
		roadSurveyIri100Map[item.RoadConditionSurveyId] = append(roadSurveyIri100Map[item.RoadConditionSurveyId], item)
	}

	var roadSurfaceLaneMaximumLanes []models.RoadSurfaceLaneMaximumLane
	query.
		Table("road_surface_lane").
		Select("max(lane_no) as num_lane, road_id").
		Group("road_id").
		Scan(&roadSurfaceLaneMaximumLanes)

	maximumLane := map[int]int{}
	for _, item := range roadSurfaceLaneMaximumLanes {
		maximumLane[item.RoadId] = item.NumLane
	}

	var roadInfoForKpi []models.RoadInfoForKpi
	query.
		Table("road_info").
		Select("road_info.road_id ,road_info.ref_direction_id as ref_direction,road_info.name as name,road_group.number as road_code ,road_section.number as section_code").
		Joins("INNER JOIN road ON road_info.road_id = road.id").
		Joins("INNER JOIN road_section ON road_section.id = road.road_section_id").
		Joins("INNER JOIN road_group ON road.road_group_id = road_group.id").
		Where("road_section.id = ?", roadSectionId).
		Where("road_info.Status = ?", "A").
		Order("road_info.km_start ASC").
		Find(&roadInfoForKpi)

	for _, itemRoadInfo := range roadInfoForKpi {

		var arreyLane []int
		var arreyLine []int
		for index := 0; index < maximumLane[itemRoadInfo.RoadId]; index++ {
			arreyLane = append(arreyLane, index+1)
		}

		for index := 0; index < maximumLane[itemRoadInfo.RoadId]+1; index++ {
			arreyLine = append(arreyLine, index+1)
		}

		var roadConditionSurveys []models.RoadConditionSurveyForKpi

		for _, itemLaneNo := range arreyLane {
			// ข้อมูล 1000 เมตร
			//check รายการล่าสุดในปีนั้น
			query.
				Table("road_condition_survey").
				Select("road_condition_survey.id as id ,road_condition.id as road_condition_id ,road_condition.road_id as road_id ,road_condition.lane_no as lane_no ,road_condition_survey.km_start as km_start,road_condition_survey.km_end as km_end ,road_condition_survey.iri as iri ,road_condition_survey.rut as rut,road_condition_survey.ifi as ifi ,road_condition_survey.survey_type as survey_type, road_condition.surveyed_date as survey_date ,road_condition.year as year").
				Joins("INNER JOIN road_condition ON road_condition_survey.road_condition_id = road_condition.id").
				Where("road_condition.road_id = ?", itemRoadInfo.RoadId).
				Where("road_condition.lane_no = ?", itemLaneNo).
				Where("road_condition.year = ?", year).
				Order("road_condition.lane_no, road_condition.km_start, road_condition_survey.km_start ASC").
				Find(&roadConditionSurveys)

			for _, itemRoadConditionSurvey := range roadConditionSurveys {
				//ได้ข้อมูล 1000 เมตร

				after90Days := itemRoadConditionSurvey.SurveyDate.AddDate(0, 0, 90)

				//set criteria ตาม survey_type
				iriCondition1000 := 0.0
				if itemRoadConditionSurvey.SurveyType == "AC" {
					iriCondition1000 = iriCondition1000Ac //ต้องใช้ค่า ชื่อเกณฑ์ 1 กม ประเภท 1 กม ค่า IRI ผิวทางลาดยาง
				} else {
					iriCondition1000 = iriCondition1000Cc //ต้องใช้ค่า ชื่อเกณฑ์ 1 กม ประเภท 1 กม ค่า IRI ผิวทางคอนกรีต
				}
				if itemRoadConditionSurvey.Iri <= iriCondition1000 {

					for _, itemRoadSurveyIri100 := range roadSurveyIri100Map[itemRoadConditionSurvey.Id] {

						iriCondition100 := 0.0
						if itemRoadConditionSurvey.SurveyType == "AC" {
							iriCondition100 = iriCondition100Ac
						} else {
							iriCondition100 = iriCondition100Cc
						}
						if itemRoadSurveyIri100.Iri <= iriCondition100 {
							//ผ่านเกณฑ์ไม่ออกรายงาน
							continue
						} else {
							//ไม่ผ่านเกณฑ์ 100 เมตร
							var comment []models.Comment
							var commentStrings []models.CommentString
							query.
								Table("maintenance").
								Select("ref_criteria_method.name as maintenance_method,maintenance_road.km_start,maintenance_road.km_end, maintenance.project_end_date").
								Joins("INNER JOIN maintenance_road ON maintenance.id = maintenance_road.maintenance_id").
								Joins("INNER JOIN ref_criteria_method ON maintenance_road.maintenance_method_id = ref_criteria_method.id").
								Where("maintenance_road.road_id = ?", itemRoadInfo.RoadId).
								Where("maintenance_road.lane_no = ?", itemRoadConditionSurvey.LaneNo).
								Where("maintenance_road.intervention_criteria_id != ?", 0).
								Where("maintenance.project_end_date BETWEEN ? AND ? ", itemRoadConditionSurvey.SurveyDate, after90Days).
								Where("maintenance.status = ?", "A").
								Find(&comment)
							for _, itemComment := range comment {
								kmStartMaint := itemComment.KmStart
								kmEndMaint := itemComment.KmEnd
								if itemRoadInfo.RefDirection == 1 {
									if (kmStartMaint <= itemRoadSurveyIri100.KmStart && kmEndMaint >= itemRoadSurveyIri100.KmStart) || (kmStartMaint <= itemRoadSurveyIri100.KmEnd && kmStartMaint >= itemRoadSurveyIri100.KmStart) {
										var commentString models.CommentString
										commentString.KmStart = KmFormat(itemComment.KmStart)
										commentString.KmEnd = KmFormat(itemComment.KmEnd)
										commentString.MaintenanceMethod = helpers.ConvertToThaiFullCalendarNoTime(after90Days)
										commentString.ProjectEndDate = helpers.ConvertToThaiFullCalendarNoTime(itemComment.ProjectEndDate)
										commentStrings = append(commentStrings, commentString)
									} else {
										continue
									}
								} else {
									if (kmStartMaint >= itemRoadSurveyIri100.KmStart && kmEndMaint <= itemRoadSurveyIri100.KmStart) || (kmStartMaint >= itemRoadSurveyIri100.KmEnd && kmStartMaint <= itemRoadSurveyIri100.KmStart) {
										var commentString models.CommentString
										commentString.KmStart = KmFormat(itemComment.KmStart)
										commentString.KmEnd = KmFormat(itemComment.KmEnd)
										commentString.MaintenanceMethod = helpers.ConvertToThaiFullCalendarNoTime(after90Days)
										commentString.ProjectEndDate = helpers.ConvertToThaiFullCalendarNoTime(itemComment.ProjectEndDate)
										commentStrings = append(commentStrings, commentString)
									} else {
										continue
									}
								}
							}
							tableResultIri100 = append(tableResultIri100, models.ResultIri100{RoadId: itemRoadInfo.RoadId, RoadName: itemRoadInfo.Name, LaneNo: itemRoadConditionSurvey.LaneNo, KmStart: KmFormat(itemRoadSurveyIri100.KmStart), KmEnd: KmFormat(itemRoadSurveyIri100.KmEnd), Iri: fmt.Sprintf("%.3f", itemRoadSurveyIri100.Iri), RoadCode: itemRoadInfo.RoadCode, SectionCode: itemRoadInfo.SectionCode, Comment: commentStrings, SurveyType: itemRoadConditionSurvey.SurveyType, MaintExpireDate: helpers.ConvertToThaiFullCalendarNoTime(after90Days), IsExpire: time.Now().After(after90Days)})
						}
					}
				} else { //ไม่ผ่านเกณฑ์ 1000 m
					// ดึงข้อมูลประวัติการซ่อมบำรุง
					var comment []models.Comment
					var commentStrings []models.CommentString
					query.
						Table("maintenance").
						Select("ref_criteria_method.name as maintenance_method,maintenance_road.km_start,maintenance_road.km_end, maintenance.project_end_date").
						Joins("INNER JOIN maintenance_road ON maintenance.id = maintenance_road.maintenance_id").
						Joins("INNER JOIN ref_criteria_method ON maintenance_road.maintenance_method_id = ref_criteria_method.id").
						Where("maintenance_road.road_id = ?", itemRoadInfo.RoadId).
						Where("maintenance_road.lane_no = ?", itemRoadConditionSurvey.LaneNo).
						Where("maintenance_road.intervention_criteria_id != ?", 0).
						Where("maintenance.project_end_date BETWEEN ? AND ? ", itemRoadConditionSurvey.SurveyDate, after90Days).
						Where("maintenance.status = ?", "A").
						Find(&comment)
					for _, itemComment := range comment {
						if itemRoadInfo.RefDirection == 1 {
							if (itemComment.KmStart <= itemRoadConditionSurvey.KmStart && itemComment.KmEnd >= itemRoadConditionSurvey.KmStart) || (itemComment.KmStart <= itemRoadConditionSurvey.KmEnd && itemComment.KmStart >= itemRoadConditionSurvey.KmStart) {
								var commentString models.CommentString
								commentString.KmStart = KmFormat(itemComment.KmStart)
								commentString.KmEnd = KmFormat(itemComment.KmEnd)
								commentString.MaintenanceMethod = helpers.ConvertToThaiFullCalendarNoTime(after90Days)
								commentString.ProjectEndDate = helpers.ConvertToThaiFullCalendarNoTime(itemComment.ProjectEndDate)
								commentStrings = append(commentStrings, commentString)
							} else {
								continue
							}
						} else {
							if (itemComment.KmStart >= itemRoadConditionSurvey.KmStart && itemComment.KmEnd <= itemRoadConditionSurvey.KmStart) || (itemComment.KmStart >= itemRoadConditionSurvey.KmEnd && itemComment.KmStart <= itemRoadConditionSurvey.KmStart) {
								var commentString models.CommentString
								commentString.KmStart = KmFormat(itemComment.KmStart)
								commentString.KmEnd = KmFormat(itemComment.KmEnd)
								commentString.MaintenanceMethod = helpers.ConvertToThaiFullCalendarNoTime(after90Days)
								commentString.ProjectEndDate = helpers.ConvertToThaiFullCalendarNoTime(itemComment.ProjectEndDate)
								commentStrings = append(commentStrings, commentString)
							} else {
								continue
							}
						}
					}
					tableResultIri1000 = append(tableResultIri1000, models.ResultIri1000{RoadId: itemRoadInfo.RoadId, RoadName: itemRoadInfo.Name, LaneNo: itemRoadConditionSurvey.LaneNo, KmStart: KmFormat(itemRoadConditionSurvey.KmStart), KmEnd: KmFormat(itemRoadConditionSurvey.KmEnd), Iri: fmt.Sprintf("%.3f", itemRoadConditionSurvey.Iri), RoadCode: itemRoadInfo.RoadCode, SectionCode: itemRoadInfo.SectionCode, Comment: commentStrings, SurveyType: itemRoadConditionSurvey.SurveyType, MaintExpireDate: helpers.ConvertToThaiFullCalendarNoTime(after90Days), IsExpire: time.Now().After(after90Days)})
				}
			}
		}
	}

	var result models.ResultIri
	result.ResultIri100 = tableResultIri100
	result.ResultIri1000 = tableResultIri1000

	return result, nil
}

func KmFormat(m int) string {
	km := strings.ReplaceAll(fmt.Sprintf("%.3f", float64(m)/1000), ".", "+")
	return km
}

func (r *Repository) Report12Ifi(roadSectionId, year int) ([]models.ResultIfi100, error) {

	query := r.conn

	var ifi100 models.ParamCondition
	var tableResultIfi100 []models.ResultIfi100

	query.
		Table("params_condition").
		Where("ref_owner_id = ?", 91).
		Where("Condition_type = ?", "IFI"). //ชื่อรายการประเภท 1000 m
		First(&ifi100)
	ifiConditionAc := ifi100.LeftValueAc
	ifiConditionCc := ifi100.LeftValueCc

	var roadSurfaceLaneMaximumLanes []models.RoadSurfaceLaneMaximumLane
	query.
		Table("road_surface_lane").
		Select("max(lane_no) as num_lane, road_id").
		Group("road_id").
		Scan(&roadSurfaceLaneMaximumLanes)

	maximumLane := map[int]int{}
	for _, item := range roadSurfaceLaneMaximumLanes {
		maximumLane[item.RoadId] = item.NumLane
	}

	var roadInfoForKpi []models.RoadInfoForKpi
	query.
		Table("road_info").
		Select("road_info.road_id ,road_info.ref_direction_id as ref_direction,road_info.name as name,road_group.number as road_code ,road_section.number as section_code").
		Joins("INNER JOIN road ON road_info.road_id = road.id").
		Joins("INNER JOIN road_section ON road_section.id = road.road_section_id").
		Joins("INNER JOIN road_group ON road.road_group_id = road_group.id").
		Where("road_section.id = ?", roadSectionId).
		Where("road_info.status = ?", "A").
		Order("road_info.km_start ASC").
		Find(&roadInfoForKpi)

	for _, itemRoadInfo := range roadInfoForKpi {

		var arreyLane []int
		var arreyLine []int
		// gen ช่องจราจร 1-num_lane
		for index := 0; index < maximumLane[itemRoadInfo.RoadId]; index++ {
			arreyLane = append(arreyLane, index+1)
		}

		for index := 0; index < maximumLane[itemRoadInfo.RoadId]+1; index++ {
			arreyLine = append(arreyLine, index+1)
		}

		var roadConditionSurveys []models.RoadConditionSurveyForKpi

		for _, itemLaneNo := range arreyLane {
			// ข้อมูล 1000 เมตร
			//check รายการล่าสุดในปีนั้น
			query.
				Table("road_condition_survey").
				Select("road_condition_survey_100m.id as id ,road_condition.id as road_condition_id,road_condition.road_id,road_condition.lane_no,road_condition_survey_100m.km_start,road_condition_survey_100m.km_end,road_condition_survey_100m.rut,road_condition_survey_100m.ifi,road_condition_survey_100m.survey_type, road_condition.surveyed_date as survey_date,road_condition.year").
				Joins("INNER JOIN road_condition ON road_condition_survey.road_condition_id = road_condition.id").
				Joins("INNER JOIN road_condition_survey_100m ON road_condition_survey.id = road_condition_survey_100m.road_condition_survey_id").
				Where("road_condition.road_id = ?", itemRoadInfo.RoadId).
				Where("road_condition.lane_no = ?", itemLaneNo).
				Where("road_condition.year = ?", year).
				Where("road_condition.status = ?", "A").
				Order("road_condition.lane_no, road_condition.km_start, road_condition_survey.km_start, road_condition_survey_100m.km_start ASC").
				Find(&roadConditionSurveys)

			for _, itemRoadConditionSurvey := range roadConditionSurveys {
				after90Days := itemRoadConditionSurvey.SurveyDate.AddDate(0, 0, 90)
				ifiCondition := 0.0
				if itemRoadConditionSurvey.SurveyType == "AC" {
					ifiCondition = ifiConditionAc //รอปรับแก้
				} else {
					ifiCondition = ifiConditionCc //รอปรับแก้
				}
				if itemRoadConditionSurvey.Ifi >= ifiCondition {
					//ผ่านเกณฑ์ไม่ออกรายงาน
					continue
				} else {
					var comment []models.Comment
					var commentStrings []models.CommentString
					query.
						Table("maintenance").
						Select("ref_criteria_method.name as maintenance_method,maintenance_road.km_start,maintenance_road.km_end, maintenance.project_end_date").
						Joins("INNER JOIN maintenance_road ON maintenance.id = maintenance_road.maintenance_id").
						Joins("INNER JOIN ref_criteria_method ON maintenance_road.maintenance_method_id = ref_criteria_method.id").
						Where("maintenance_road.road_id = ?", itemRoadInfo.RoadId).
						Where("maintenance_road.lane_no = ?", itemRoadConditionSurvey.LaneNo).
						Where("maintenance_road.intervention_criteria_id != ?", 0).
						Where("maintenance.project_end_date BETWEEN ? AND ? ", itemRoadConditionSurvey.SurveyDate, after90Days).
						Where("maintenance.status = ?", "A").
						Find(&comment)
					for _, itemComment := range comment {
						if itemRoadInfo.RefDirection == 1 {
							if (itemComment.KmStart <= itemRoadConditionSurvey.KmStart && itemComment.KmEnd >= itemRoadConditionSurvey.KmStart) || (itemComment.KmStart <= itemRoadConditionSurvey.KmEnd && itemComment.KmStart >= itemRoadConditionSurvey.KmStart) {
								var commentString models.CommentString
								commentString.KmStart = KmFormat(itemComment.KmStart)
								commentString.KmEnd = KmFormat(itemComment.KmEnd)
								commentString.MaintenanceMethod = helpers.ConvertToThaiFullCalendarNoTime(after90Days)
								commentString.ProjectEndDate = helpers.ConvertToThaiFullCalendarNoTime(itemComment.ProjectEndDate)
								commentStrings = append(commentStrings, commentString)
							} else {
								continue
							}
						} else {
							if (itemComment.KmStart >= itemRoadConditionSurvey.KmStart && itemComment.KmEnd <= itemRoadConditionSurvey.KmStart) || (itemComment.KmStart >= itemRoadConditionSurvey.KmEnd && itemComment.KmStart <= itemRoadConditionSurvey.KmStart) {
								var commentString models.CommentString
								commentString.KmStart = KmFormat(itemComment.KmStart)
								commentString.KmEnd = KmFormat(itemComment.KmEnd)
								commentString.MaintenanceMethod = helpers.ConvertToThaiFullCalendarNoTime(after90Days)
								commentString.ProjectEndDate = helpers.ConvertToThaiFullCalendarNoTime(itemComment.ProjectEndDate)
								commentStrings = append(commentStrings, commentString)
							} else {
								continue
							}
						}
					}
					tableResultIfi100 = append(tableResultIfi100, models.ResultIfi100{RoadId: itemRoadInfo.RoadId, RoadName: itemRoadInfo.Name, LaneNo: itemRoadConditionSurvey.LaneNo, KmStart: KmFormat(itemRoadConditionSurvey.KmStart), KmEnd: KmFormat(itemRoadConditionSurvey.KmEnd), Ifi: fmt.Sprintf("%.3f", itemRoadConditionSurvey.Ifi), RoadCode: itemRoadInfo.RoadCode, SectionCode: itemRoadInfo.SectionCode, Comment: commentStrings, SurveyType: itemRoadConditionSurvey.SurveyType, MaintExpireDate: helpers.ConvertToThaiFullCalendarNoTime(after90Days), IsExpire: time.Now().After(after90Days)})
				}
			}
		}
	}
	return tableResultIfi100, nil
}

func (r *Repository) Report12Rut(roadSectionId, year int) ([]models.ResultRut100, error) {
	query := r.conn

	var rut100 models.ParamCondition

	//rut
	query.
		Table("params_condition").
		Where("ref_owner_id = ?", 91).
		Where("Condition_type = ?", "RUT"). //ชื่อรายการประเภท 1000 m
		First(&rut100)
	rutConditionAc := rut100.RightValueAc

	//สร้างตารางเก็บผลลัพธ์
	var tableResultRut100 []models.ResultRut100

	var roadSurfaceLaneMaximumLanes []models.RoadSurfaceLaneMaximumLane
	query.
		Table("road_surface_lane").
		Select("max(lane_no) as num_lane, road_id").
		Group("road_id").
		Scan(&roadSurfaceLaneMaximumLanes)

	maximumLane := map[int]int{}
	for _, item := range roadSurfaceLaneMaximumLanes {
		maximumLane[item.RoadId] = item.NumLane
	}

	var roadInfoForKpi []models.RoadInfoForKpi
	query.
		Table("road_info").
		Select("road_info.road_id ,road_info.ref_direction_id as ref_direction,road_info.name as name,road_group.number as road_code ,road_section.number as section_code").
		Joins("INNER JOIN road ON road_info.road_id = road.id").
		Joins("INNER JOIN road_section ON road_section.id = road.road_section_id").
		Joins("INNER JOIN road_group ON road.road_group_id = road_group.id").
		Where("road_section.id = ?", roadSectionId).
		Where("road_info.status = ?", "A").
		Order("road_info.km_start ASC").
		Find(&roadInfoForKpi)

	for _, itemRoadInfo := range roadInfoForKpi {

		var arreyLane []int
		var arreyLine []int

		// gen ช่องจราจร 1-num_lane
		for index := 0; index < maximumLane[itemRoadInfo.RoadId]; index++ {
			arreyLane = append(arreyLane, index+1)
		}

		for index := 0; index < maximumLane[itemRoadInfo.RoadId]+1; index++ {
			arreyLine = append(arreyLine, index+1)
		}

		var roadConditionSurveys []models.RoadConditionSurveyForKpi

		for _, itemLaneNo := range arreyLane {
			// ข้อมูล 1000 เมตร
			//check รายการล่าสุดในปีนั้น
			query.
				Table("road_condition_survey").
				Select("road_condition_survey_100m.id as id ,road_condition.id as road_condition_id,road_condition.road_id,road_condition.lane_no,road_condition_survey_100m.km_start,road_condition_survey_100m.km_end,road_condition_survey_100m.rut,road_condition_survey_100m.ifi,road_condition_survey_100m.survey_type, road_condition.surveyed_date as survey_date,road_condition.year").
				Joins("INNER JOIN road_condition ON road_condition_survey.road_condition_id = road_condition.id").
				Joins("INNER JOIN road_condition_survey_100m ON road_condition_survey.id = road_condition_survey_100m.road_condition_survey_id").
				Where("road_condition.road_id = ?", itemRoadInfo.RoadId).
				Where("road_condition.lane_no = ?", itemLaneNo).
				Where("road_condition.year = ?", year).
				Where("road_condition.status = ?", "A").
				Order("road_condition.lane_no, road_condition.km_start, road_condition_survey.km_start, road_condition_survey_100m.km_start  ASC").
				Find(&roadConditionSurveys)

			for _, itemRoadConditionSurvey := range roadConditionSurveys {
				after7Days := itemRoadConditionSurvey.SurveyDate.AddDate(0, 0, 7)
				rutCondition := 0.0
				if itemRoadConditionSurvey.SurveyType == "AC" {
					rutCondition = rutConditionAc
				} else {
					continue
				}
				if itemRoadConditionSurvey.Rut <= rutCondition {
					//ผ่านเกณฑ์ไม่ออกรายงาน
					continue
				} else {
					//ไม่ผ่านเกณฑ์ RFP
					var comment []models.Comment
					var commentStrings []models.CommentString
					query.
						Table("maintenance").
						Select("ref_criteria_method.name as maintenance_method,maintenance_road.km_start,maintenance_road.km_end, maintenance.project_end_date").
						Joins("INNER JOIN maintenance_road ON maintenance.id = maintenance_road.maintenance_id").
						Joins("INNER JOIN ref_criteria_method ON maintenance_road.maintenance_method_id = ref_criteria_method.id").
						Where("maintenance_road.road_id = ?", itemRoadInfo.RoadId).
						Where("maintenance_road.lane_no = ?", itemRoadConditionSurvey.LaneNo).
						Where("maintenance_road.intervention_criteria_id != ?", 0).
						Where("maintenance.project_end_date BETWEEN ? AND ? ", itemRoadConditionSurvey.SurveyDate, after7Days).
						Where("maintenance.status = ?", "A").
						Find(&comment)
					for _, itemComment := range comment {
						if itemRoadInfo.RefDirection == 1 {
							if (itemComment.KmStart <= itemRoadConditionSurvey.KmStart && itemComment.KmEnd >= itemRoadConditionSurvey.KmStart) || (itemComment.KmStart <= itemRoadConditionSurvey.KmEnd && itemComment.KmStart >= itemRoadConditionSurvey.KmStart) {
								var commentString models.CommentString
								commentString.KmStart = KmFormat(itemComment.KmStart)
								commentString.KmEnd = KmFormat(itemComment.KmEnd)
								commentString.MaintenanceMethod = helpers.ConvertToThaiFullCalendarNoTime(after7Days)
								commentString.ProjectEndDate = helpers.ConvertToThaiFullCalendarNoTime(itemComment.ProjectEndDate)
								commentStrings = append(commentStrings, commentString)
							} else {
								continue
							}
						} else {
							if (itemComment.KmStart >= itemRoadConditionSurvey.KmStart && itemComment.KmEnd <= itemRoadConditionSurvey.KmStart) || (itemComment.KmStart >= itemRoadConditionSurvey.KmEnd && itemComment.KmStart <= itemRoadConditionSurvey.KmStart) {
								var commentString models.CommentString
								commentString.KmStart = KmFormat(itemComment.KmStart)
								commentString.KmEnd = KmFormat(itemComment.KmEnd)
								commentString.MaintenanceMethod = helpers.ConvertToThaiFullCalendarNoTime(after7Days)
								commentString.ProjectEndDate = helpers.ConvertToThaiFullCalendarNoTime(itemComment.ProjectEndDate)
								commentStrings = append(commentStrings, commentString)
							} else {
								continue
							}
						}
					}
					tableResultRut100 = append(tableResultRut100, models.ResultRut100{RoadId: itemRoadInfo.RoadId, RoadName: itemRoadInfo.Name, LaneNo: itemRoadConditionSurvey.LaneNo, KmStart: KmFormat(itemRoadConditionSurvey.KmStart), KmEnd: KmFormat(itemRoadConditionSurvey.KmEnd), Rut: fmt.Sprintf("%.3f", itemRoadConditionSurvey.Rut), RoadCode: itemRoadInfo.RoadCode, SectionCode: itemRoadInfo.SectionCode, Comment: commentStrings, SurveyType: itemRoadConditionSurvey.SurveyType, MaintExpireDate: helpers.ConvertToThaiFullCalendarNoTime(after7Days), IsExpire: time.Now().After(after7Days)})
				}
			}
		}
	}
	return tableResultRut100, nil
}

func (r *Repository) Report12G7(roadSectionId, year int) ([]models.ResultG7100, error) {

	query := r.conn

	var g7100 models.G7Condition

	query.
		Table("params_road_line").
		Where("Ref_owner_road_line_id = ?", 13).
		First(&g7100)
	g7ConditionYellow := g7100.LeftValueYellow
	g7ConditionWhite := g7100.LeftValueWhite

	//สร้างตารางเก็บผลลัพธ์
	var tableResultG7100 []models.ResultG7100

	var roadSurfaceLaneMaximumLanes []models.RoadSurfaceLaneMaximumLane
	query.
		Table("road_surface_lane").
		Select("max(lane_no) as num_lane, road_id").
		Group("road_id").
		Scan(&roadSurfaceLaneMaximumLanes)

	maximumLane := map[int]int{}
	for _, item := range roadSurfaceLaneMaximumLanes {
		maximumLane[item.RoadId] = item.NumLane
	}

	var roadInfoForKpi []models.RoadInfoForKpi
	query.
		Table("road_info").
		Select("road_info.road_id ,road_info.ref_direction_id as ref_direction,road_info.name as name,road_group.number as road_code ,road_section.number as section_code").
		Joins("INNER JOIN road ON road_info.road_id = road.id").
		Joins("INNER JOIN road_section ON road_section.id = road.road_section_id").
		Joins("INNER JOIN road_group ON road.road_group_id = road_group.id").
		Where("road_section.id = ?", roadSectionId).
		Where("road_info.status = ?", "A").
		Order("road_info.km_start ASC").
		Find(&roadInfoForKpi)

	for _, itemRoadInfo := range roadInfoForKpi {

		var arreyLane []int
		var arreyLine []int
		// gen ช่องจราจร 1-num_lane
		for index := 0; index < maximumLane[itemRoadInfo.RoadId]; index++ {
			arreyLane = append(arreyLane, index+1)
		}

		for index := 0; index < maximumLane[itemRoadInfo.RoadId]+1; index++ {
			arreyLine = append(arreyLine, index+1)
		}

		var roadG7100s []models.RoadG7100

		for _, itemLaneNo := range arreyLine {
			//check รายการล่าสุดในปีนั้น
			query.
				Table("road_retro_reflectivity").
				Select("road_retro_reflectivity.id as id , road_retro_reflectivity.road_id,road_retro_reflectivity.line_no,road_retro_reflectivity_range.km_start,road_retro_reflectivity_range.km_end,road_retro_reflectivity_range.retro_avg as retro,ref_stripe_color.name_th as name_strip, road_retro_reflectivity.surveyed_date as survey_date,road_retro_reflectivity.year").
				Joins("INNER JOIN road_retro_reflectivity_range ON road_retro_reflectivity.id = road_retro_reflectivity_range.road_retro_reflectivity_id").
				Joins("INNER JOIN ref_stripe_color ON ref_stripe_color.id = road_retro_reflectivity_range.ref_stripe_color_id").
				Where("road_retro_reflectivity.road_id = ?", itemRoadInfo.RoadId).
				Where("road_retro_reflectivity.line_no = ?", itemLaneNo).
				Where("road_retro_reflectivity.year = ?", year).
				Where("road_retro_reflectivity.status = ?", "A").
				Order("road_retro_reflectivity.line_no, road_retro_reflectivity.km_start, road_retro_reflectivity_range.km_start ASC").
				Find(&roadG7100s)

			for _, itemRoadG7100 := range roadG7100s {
				after7Days := itemRoadG7100.SurveyDate.AddDate(0, 0, 7)
				retroCondition := 0.0
				if itemRoadG7100.NameStrip == "เส้นสีขาว" {
					retroCondition = g7ConditionWhite
				} else {
					retroCondition = g7ConditionYellow
				}
				if itemRoadG7100.Retro >= retroCondition {
					//ผ่านเกณฑ์ไม่ออกรายงาน
					continue
				} else {
					//ไม่ผ่านเกณฑ์ RFP
					var comment []models.Comment
					var commentStrings []models.CommentString
					query.
						Table("maintenance").
						Select("ref_criteria_method.name as maintenance_method,maintenance_road.km_start,maintenance_road.km_end, maintenance.project_end_date").
						Joins("INNER JOIN maintenance_road ON maintenance.id = maintenance_road.maintenance_id").
						Joins("INNER JOIN ref_criteria_method ON maintenance_road.maintenance_method_id = ref_criteria_method.id").
						Where("maintenance_road.road_id = ?", itemRoadInfo.RoadId).
						Where("maintenance_road.lane_no = ?", itemRoadG7100.LineNo).
						Where("maintenance_road.intervention_criteria_id != ?", 0).
						Where("maintenance.project_end_date BETWEEN ? AND ? ", itemRoadG7100.SurveyDate, after7Days).
						Where("maintenance.status = ?", "A").
						Find(&comment)
					for _, itemComment := range comment {
						if itemRoadInfo.RefDirection == 1 {
							if (itemComment.KmStart <= itemRoadG7100.KmStart && itemComment.KmEnd >= itemRoadG7100.KmStart) || (itemComment.KmStart <= itemRoadG7100.KmEnd && itemComment.KmStart >= itemRoadG7100.KmStart) {
								var commentString models.CommentString
								commentString.KmStart = KmFormat(itemComment.KmStart)
								commentString.KmEnd = KmFormat(itemComment.KmEnd)
								commentString.MaintenanceMethod = helpers.ConvertToThaiFullCalendarNoTime(after7Days)
								commentString.ProjectEndDate = helpers.ConvertToThaiFullCalendarNoTime(itemComment.ProjectEndDate)
								commentStrings = append(commentStrings, commentString)
							} else {
								continue
							}
						} else {
							if (itemComment.KmStart >= itemRoadG7100.KmStart && itemComment.KmEnd <= itemRoadG7100.KmStart) || (itemComment.KmStart >= itemRoadG7100.KmEnd && itemComment.KmStart <= itemRoadG7100.KmStart) {
								var commentString models.CommentString
								commentString.KmStart = KmFormat(itemComment.KmStart)
								commentString.KmEnd = KmFormat(itemComment.KmEnd)
								commentString.MaintenanceMethod = helpers.ConvertToThaiFullCalendarNoTime(after7Days)
								commentString.ProjectEndDate = helpers.ConvertToThaiFullCalendarNoTime(itemComment.ProjectEndDate)
								commentStrings = append(commentStrings, commentString)
							} else {
								continue
							}
						}
					}
					tableResultG7100 = append(tableResultG7100, models.ResultG7100{RoadId: itemRoadInfo.RoadId, RoadName: itemRoadInfo.Name, LineNo: itemRoadG7100.LineNo, KmStart: KmFormat(itemRoadG7100.KmStart), KmEnd: KmFormat(itemRoadG7100.KmEnd), Retro: fmt.Sprintf("%.3f", itemRoadG7100.Retro), RoadCode: itemRoadInfo.RoadCode, SectionCode: itemRoadInfo.SectionCode, Comment: commentStrings, NameStrip: itemRoadG7100.NameStrip, MaintExpireDate: helpers.ConvertToThaiFullCalendarNoTime(after7Days), IsExpire: time.Now().After(after7Days)})
				}
			}
		}
	}
	return tableResultG7100, nil
}

//////////////// NEW MIMS ////////////////
