package repository

import (
	"fmt"
	"time"

	"gitlab.com/mims-api-service/models"
	"gorm.io/gorm"
)

func (r *Repository) GetRoadSurface() ([]models.RoadSurface, error) {
	var roadSurface []models.RoadSurface
	query := r.conn
	if err := query.Raw("select * from road_surface where road_id = 2 and status = 'A' and revision = (select max(revision) from road_surface where road_id = 2 and status = 'A') order by km_start").Find(&roadSurface).Error; err != nil {
		return roadSurface, err
	}
	return roadSurface, nil
}

func (r *Repository) GetMaintenanceData(roadId, laneNo int) ([]models.MaintenanceData, error) {
	var maintenanceData []models.MaintenanceData
	query := r.conn

	query = query.Preload("Budget")
	query = query.Preload("BudgetMethod")
	query = query.Preload("MaintenanceRoads", func(db *gorm.DB) *gorm.DB {
		return db.Where("lane_no = ?", laneNo).Where("status = ?", "A").Where("road_id = ?", roadId).Where("intervention_criteria_id != ?", 0)
	})
	query = query.Preload("MaintenanceRoads.RefSurface")
	query = query.Preload("MaintenanceRoads.RoadGroup")
	query = query.Preload("MaintenanceRoads.RoadInfo")
	query = query.Preload("MaintenanceRoads.RoadInfo.Direction")

	if err := query.Where("status = 'A'").Find(&maintenanceData).Error; err != nil {
		return maintenanceData, err
	}
	return maintenanceData, nil
}

func (r *Repository) GetMaintenanceByID(mID int) (models.MaintenanceData, error) {
	var maintenanceData models.MaintenanceData
	query := r.conn

	query = query.Preload("MaintenanceRoad", func(db *gorm.DB) *gorm.DB {
		return db.Where("lane = ?", 1)
	})

	if err := query.Where("is_complete = ?", true).Where("id = ?", mID).First(&maintenanceData).Error; err != nil {
		return maintenanceData, err
	}
	return maintenanceData, nil
}

func (r *Repository) GetRoadGroups() ([]models.RoadGroup, error) {
	var roadGroup []models.RoadGroup
	query := r.conn
	if err := query.Find(&roadGroup).Error; err != nil {
		return roadGroup, err
	}
	return roadGroup, nil
}

func (r *Repository) GetRoadSections() ([]models.RoadSection, error) {
	var roadSection []models.RoadSection
	query := r.conn
	if err := query.Find(&roadSection).Error; err != nil {
		return roadSection, err
	}
	return roadSection, nil
}

func (r *Repository) GetRoadGroupByID(ID int) (models.RoadGroup, error) {
	var roadGroup models.RoadGroup
	query := r.conn
	if err := query.Where("id = ?", ID).First(&roadGroup).Error; err != nil {
		return roadGroup, err
	}
	return roadGroup, nil
}

func (r *Repository) GetRoadInfoByID(roadID int) (models.RoadInfo, error) {
	var roadInfo models.RoadInfo
	query := r.conn
	if err := query.Where("road_id = ?", roadID).First(&roadInfo).Error; err != nil {
		return roadInfo, err
	}
	return roadInfo, nil
}

func (r *Repository) GetRoadByID(roadID int) (models.Road, error) {
	var road models.Road
	query := r.conn
	if err := query.Where("id = ?", roadID).First(&road).Error; err != nil {
		return road, err
	}
	return road, nil
}

func (r *Repository) GetRoadDatebegin(roadID int) ([]models.RoadDatebegin, error) {
	var roadDatebegin []models.RoadDatebegin
	query := r.conn
	if err := query.Where("road_id = ?", roadID).Find(&roadDatebegin).Error; err != nil {
		return roadDatebegin, err
	}
	return roadDatebegin, nil
}

func (r *Repository) GetRoadByRoadID(roadID int) (models.Road, error) {
	var road models.Road
	query := r.conn
	if err := query.Where("id = ?", roadID).First(&road).Error; err != nil {
		return road, err
	}
	return road, nil
}

func (r *Repository) GetRoadGeomLaneByID(roadID, laneNo int) (models.RoadGeom, error) {
	var roadInfo models.RoadGeom
	query := r.conn
	if err := query.Where("road_id = ?", roadID).Where("lane_no = ?", laneNo).Where("status = ?", "A").First(&roadInfo).Error; err != nil {
		return roadInfo, err
	}
	return roadInfo, nil
}

func (r *Repository) GetRefStructureSurfaceByID(ID int) (models.RefStructureSurface, error) {
	var refStructureSurface models.RefStructureSurface
	query := r.conn
	if err := query.Where("id = ?", ID).First(&refStructureSurface).Error; err != nil {
		return refStructureSurface, err
	}
	return refStructureSurface, nil
}

func (r *Repository) GetRoadSurfaceLaneBySurfaceIDByLane(surfaceID, laneNo int) (models.RoadSurfaceLane, error) {
	var roadSurfaceLane models.RoadSurfaceLane
	query := r.conn
	if err := query.Where("road_surface_id = ?", surfaceID).Where("lane_no = ?", laneNo).First(&roadSurfaceLane).Error; err != nil {
		return roadSurfaceLane, err
	}
	return roadSurfaceLane, nil
}

func (r *Repository) GetRefSurfaceByID(ID int) (models.RefSurface, error) {
	var refSurface models.RefSurface
	query := r.conn
	if err := query.Where("id = ?", ID).First(&refSurface).Error; err != nil {
		return refSurface, err
	}
	return refSurface, nil
}

func (r *Repository) GetRefSurfaceParam(ID int) (models.RefSurfaceParam, error) {
	var refSurfaceParam models.RefSurfaceParam
	query := r.conn
	if err := query.Where("id = ?", ID).First(&refSurfaceParam).Error; err != nil {
		return refSurfaceParam, err
	}
	return refSurfaceParam, nil
}

func (r *Repository) GetRefMaterialBaseByID(ID int) (models.RefMaterialBase, error) {
	var refMaterialBase models.RefMaterialBase
	query := r.conn
	if err := query.Where("id = ?", ID).First(&refMaterialBase).Error; err != nil {
		return refMaterialBase, err
	}
	return refMaterialBase, nil
}

func (r *Repository) GetRefMaterialSubgradeByID(ID int) (models.RefMaterialSubgrade, error) {
	var refMaterialSubgrade models.RefMaterialSubgrade
	query := r.conn
	if err := query.Where("id = ?", ID).First(&refMaterialSubgrade).Error; err != nil {
		return refMaterialSubgrade, err
	}
	return refMaterialSubgrade, nil
}

func (r *Repository) GetRefMaterialSubbaseByID(ID int) (models.RefMaterialSubbase, error) {
	var refMaterialSubbase models.RefMaterialSubbase
	query := r.conn
	if err := query.Where("id = ?", ID).First(&refMaterialSubbase).Error; err != nil {
		return refMaterialSubbase, err
	}
	return refMaterialSubbase, nil
}

func (r *Repository) GetInterventionCriteria() ([]models.InterventionCriteria, error) {
	var interventionCriteria []models.InterventionCriteria
	query := r.conn
	if err := query.Where("is_deleted = ?", false).Find(&interventionCriteria).Error; err != nil {
		return interventionCriteria, err
	}
	return interventionCriteria, nil
}

func (r *Repository) GetCriteriaMethod() ([]models.RefCriteriaMethod, error) {
	var refCriteriaMethod []models.RefCriteriaMethod
	query := r.conn
	if err := query.Where("is_deleted = ?", false).Find(&refCriteriaMethod).Error; err != nil {
		return refCriteriaMethod, err
	}
	return refCriteriaMethod, nil
}

func (r *Repository) GetRoadWorkEffect() (models.SettingRoadWorkEffect, error) {
	var settingRoadWorkEffect models.SettingRoadWorkEffect
	query := r.conn
	if err := query.First(&settingRoadWorkEffect).Error; err != nil {
		return settingRoadWorkEffect, err
	}
	return settingRoadWorkEffect, nil
}

func (r *Repository) GetSettingAadtParams() (models.SettingAadtParams, error) {
	var settingAadtParams models.SettingAadtParams
	query := r.conn
	if err := query.Where("is_latest =?", true).First(&settingAadtParams).Error; err != nil {
		return settingAadtParams, err
	}
	return settingAadtParams, nil
}

func (r *Repository) GetRoadSurfaceByGroupRoadID(roadID int, laneNo int) ([]int, error) {
	var roadSurface []models.RoadSurface
	if err := r.conn.Select("surface_group").Where("status = ? or status = ?", "A", "I").Where("road_id = ?", roadID).Order("surface_group DESC").Group("surface_group").Find(&roadSurface).Error; err != nil {
		return []int{}, err
	}
	surfaceCnt := []int{}
	for _, item := range roadSurface {
		surfaceCnt = append(surfaceCnt, item.SurfaceGroup)
	}
	return surfaceCnt, nil
}

func (r *Repository) GetRoadSurfaceByRoadID(roadID, laneNo, grp int) ([]models.RoadSurfacePrepareData, error) {

	var roadSurfacePrepareData []models.RoadSurfacePrepareData
	query := r.conn
	query = query.Preload("RefStructureSurface")
	query = query.Preload("RefSurfaceShoulderLeft")
	query = query.Preload("RoadSurfaceLane", func(db *gorm.DB) *gorm.DB {
		return db.Where("lane_no = ?", laneNo)
	})
	query = query.Preload("RoadSurfaceLane.RefSurface")
	query = query.Preload("RefMaterialBase")
	query = query.Preload("RefMaterialSubbase")
	query = query.Preload("RefMaterialSubgrade")
	if err := query.Where("status = ? or status = ?", "A", "I").Where("road_id = ?", roadID).Where("surface_group = ?", grp).Find(&roadSurfacePrepareData).Error; err != nil {
		return []models.RoadSurfacePrepareData{}, nil
	}

	if len(roadSurfacePrepareData) > 0 {
		return roadSurfacePrepareData, nil
	}

	return []models.RoadSurfacePrepareData{}, nil
}

func (r *Repository) GetRoadSurfaceFirstGrpByRoadID(roadID int, laneNo int) ([]models.RoadSurfacePrepareData, error) {
	var roadSurface models.RoadSurface
	if err := r.conn.Select("max(surface_group) as surface_group").Where("status = ? or status = ?", "A", "I").Where("road_id = ?", roadID).Find(&roadSurface).Error; err != nil {
		return []models.RoadSurfacePrepareData{}, err
	}
	maxSurfaceGrp := roadSurface.SurfaceGroup

	var roadSurfacePrepareData []models.RoadSurfacePrepareData
	query := r.conn
	query = query.Preload("RefStructureSurface")
	query = query.Preload("RefSurfaceShoulderLeft")
	query = query.Preload("RoadSurfaceLane", func(db *gorm.DB) *gorm.DB {
		return db.Where("lane_no = ?", laneNo)
	})
	query = query.Preload("RefMaterialBase")
	query = query.Preload("RefMaterialSubbase")
	query = query.Preload("RefMaterialSubgrade")
	if err := query.Where("status = ? or status = ?", "A", "I").Where("road_id = ?", roadID).Where("surface_group = ?", maxSurfaceGrp).Find(&roadSurfacePrepareData).Error; err != nil {
		return []models.RoadSurfacePrepareData{}, nil
	}

	if len(roadSurfacePrepareData) > 0 {
		return roadSurfacePrepareData, nil
	}

	return []models.RoadSurfacePrepareData{}, nil

}

func (r *Repository) GetRoadSurfaceByRoadIDYear(roadID int, laneNo int, year int) ([]models.RoadSurfacePrepareData, error) {
	var roadSurfacePrepareData []models.RoadSurfacePrepareData
	query := r.conn
	query = query.Preload("RefStructureSurface")
	query = query.Preload("RefSurfaceShoulderLeft")
	query = query.Preload("RoadSurfaceLane", func(db *gorm.DB) *gorm.DB {
		return db.Where("lane_no = ?", laneNo)
	})
	query = query.Preload("RoadSurfaceLane.RefSurface")
	query = query.Preload("RefMaterialBase")
	query = query.Preload("RefMaterialSubbase")
	query = query.Preload("RefMaterialSubgrade")
	if err := query.Where("status = ?", "A").Where("year = ?", year).Where("road_id = ?", roadID).Find(&roadSurfacePrepareData).Error; err != nil {
		return roadSurfacePrepareData, err
	}
	return roadSurfacePrepareData, nil
}

func (r *Repository) GetInterventionCriteriaParamsById(ID int) (models.SettingInterventionCriteriaParams, error) {
	var interventionCriteria models.SettingInterventionCriteriaParams
	query := r.conn
	err := query.Where("id = ?", ID).First(&interventionCriteria).Error
	if err != nil {
		return interventionCriteria, err
	}
	return interventionCriteria, nil
}

func (r *Repository) GetRefCriteriaMethodByID(ID int) (models.RefCriteriaMethod, error) {
	var refCriteriaMethod models.RefCriteriaMethod
	query := r.conn
	err := query.Where("id = ?", ID).First(&refCriteriaMethod).Error
	if err != nil {
		return refCriteriaMethod, err
	}
	return refCriteriaMethod, nil
}

func (r *Repository) GetRoadCondition(roadID, laneNo int) ([]models.RoadConditionSurveyM, time.Time, error) {
	// query := r.conn
	// var roadCondition models.RoadCondition
	// var roadConditionSurveyM []models.RoadConditionSurveyM
	// err := query.Select("max(year) as year").Where("road_id = ?", roadID).Where("lane_no = ?", laneNo).Where("status = ?", "A").Find(&roadCondition).Error
	// if err != nil {
	// 	return []models.RoadConditionSurveyM{}, time.Time{}, err
	// }

	// err = query.Where("road_id = ?", roadID).Where("year = ?", roadCondition.Year).Where("lane_no = ?", laneNo).Where("status = ?", "A").First(&roadCondition).Error
	// if err != nil {
	// 	return []models.RoadConditionSurveyM{}, time.Time{}, err
	// }

	// var roadConditionSurvey []models.RoadConditionSurvey
	// err = query.Where("road_condition_id = ?", roadCondition.ID).Find(&roadConditionSurvey).Error
	// if err != nil {
	// 	return []models.RoadConditionSurveyM{}, time.Time{}, err
	// }
	// IDs := []int{}
	// for _, item := range roadConditionSurvey {
	// 	IDs = append(IDs, item.ID)
	// }

	// err = query.Where("road_condition_survey_id IN (?)", IDs).Find(&roadConditionSurveyM).Error
	// if err != nil {
	// 	return []models.RoadConditionSurveyM{}, time.Time{}, err
	// }

	// return roadConditionSurveyM, roadCondition.SurveyedDate, nil
	query := r.conn.Debug()
	var roadCondition models.RoadCondition
	var roadConditionSurveyM []models.RoadConditionSurveyM
	err := query.Select("max(year) as year").Where("road_id = ?", roadID).Where("lane_no = ?", laneNo).Where("status = ?", "A").Find(&roadCondition).Error
	if err != nil {
		return []models.RoadConditionSurveyM{}, time.Time{}, err
	}

	err = query.Where("road_id = ?", roadID).Where("year = ?", roadCondition.Year).Where("lane_no = ?", laneNo).Where("status = ?", "A").First(&roadCondition).Error
	if err != nil {
		return []models.RoadConditionSurveyM{}, time.Time{}, err
	}

	// Subquery for MAX Year for Road Conditions
	maxYearSubQuery := query.Model(&models.RoadCondition{}).Select("MAX(year)").Where("road_id = ? AND lane_no = ? AND status = ?", roadID, laneNo, "A")
	// Subquery for getting Road Condition ID
	roadConditionIdSubQuery := query.Model(&models.RoadCondition{}).Select("id").Where("road_id = ? AND lane_no = ? AND status = ? AND year = (?)", roadID, laneNo, "A", maxYearSubQuery).Limit(1)
	// Subquery for getting Road Condition Survey ID
	roadConditionSurveyIdSubQuery := query.Model(&models.RoadConditionSurvey{}).Select("id").Where("road_condition_id IN (?)", roadConditionIdSubQuery)
	// Subquery for getting Road Condition Survey 100m ID
	roadConditionSurvey100mIdSubQuery := query.Model(&models.RoadConditionSurvey100M{}).Select("id").Where("road_condition_survey_id IN (?)", roadConditionSurveyIdSubQuery)

	// Main query to get required records from RoadConditionSurveyM table
	// var roadConditionSurveyM []models.RoadConditionSurveyM
	result := query.Where("road_condition_survey_100m_id IN (?)", roadConditionSurvey100mIdSubQuery).
		Find(&roadConditionSurveyM)

	// Error handling and processing the result
	if result.Error != nil {
		return []models.RoadConditionSurveyM{}, time.Time{}, err
	}
	if err != nil {
		return []models.RoadConditionSurveyM{}, time.Time{}, err
	}
	return roadConditionSurveyM, roadCondition.SurveyedDate, nil
}
func (r *Repository) GetRoadDamage(roadID, laneNo int) ([]models.RoadDamageMPrepareData, time.Time, error) {
	query := r.conn
	var roadDamage models.RoadDamage
	var roadDamageM []models.RoadDamageMPrepareData
	err := query.Select("max(year) as year").Where("road_id = ?", roadID).Where("lane_no = ?", laneNo).Where("status = ?", "A").Find(&roadDamage).Error
	if err != nil {
		return []models.RoadDamageMPrepareData{}, time.Time{}, err
	}
	// fmt.Println("yesr", roadCondition.Year)

	err = query.Where("road_id = ?", roadID).Where("year = ?", roadDamage.Year).Where("lane_no = ?", 1).Where("status = ?", "A").First(&roadDamage).Error
	if err != nil {
		return []models.RoadDamageMPrepareData{}, time.Time{}, err
	}

	var roadDamageRange []models.RoadDamageRange
	err = query.Where("road_damage_id = ?", roadDamage.Id).Find(&roadDamageRange).Error
	if err != nil {
		return []models.RoadDamageMPrepareData{}, time.Time{}, err
	}
	IDs := []int{}
	for _, item := range roadDamageRange {
		IDs = append(IDs, item.Id)
	}

	err = query.Where("road_damage_range_id IN (?)", IDs).Find(&roadDamageM).Error
	if err != nil {
		return []models.RoadDamageMPrepareData{}, time.Time{}, err
	}

	return roadDamageM, roadDamage.SurveyedDate, nil
}

func (r *Repository) GetVolumeAadt(roadId int) (models.VolumeAadt, error) {
	query := r.conn
	var volumeAadt models.VolumeAadt
	err := query.Select("max(year) as year").Where("road_id = ?", roadId).Where("status = ?", "A").Find(&volumeAadt).Error
	if err != nil {
		return models.VolumeAadt{}, err
	}
	// fmt.Println("yesr", roadCondition.Year)

	err = query.Where("road_id = ?", roadId).Where("year = ?", volumeAadt.Year).Where("status = ?", "A").Order("surveyed_date DESC").First(&volumeAadt).Error
	if err != nil {
		return models.VolumeAadt{}, err
	}
	return volumeAadt, nil
}

func (r *Repository) GetVolumeRain(roadGroupID int) (models.VolumeRain, error) {
	query := r.conn
	var volumeRain models.VolumeRain
	err := query.Select("max(year) as year").Where("road_group_id = ?", roadGroupID).Where("status = ?", "A").Find(&volumeRain).Error
	if err != nil {
		return models.VolumeRain{}, err
	}

	err = query.Where("road_group_id = ?", roadGroupID).Where("year = ?", volumeRain.Year).Where("status = ?", "A").First(&volumeRain).Error
	if err != nil {
		return models.VolumeRain{}, err
	}
	return volumeRain, nil
}

func (r *Repository) GetSettingRoadWorkEffect() (models.SettingRoadWorkEffectParams, error) {
	var settingRoadWorkEffectParams models.SettingRoadWorkEffectParams
	query := r.conn
	err := query.Where("is_latest = ?", true).First(&settingRoadWorkEffectParams).Error
	if err != nil {
		return models.SettingRoadWorkEffectParams{}, err
	}
	return settingRoadWorkEffectParams, nil
}

func (r *Repository) GetRrefSurfacrParamsById(ID int) (models.RefSurfaceParam, error) {
	var refSurfaceParam models.RefSurfaceParam
	query := r.conn
	err := query.Where("id = ?", ID).First(&refSurfaceParam).Error
	if err != nil {
		return refSurfaceParam, err
	}

	return refSurfaceParam, nil
}

func (r *Repository) GetRefMaterialBaseConcretePavement() (models.RefMaterialBase, error) {
	var refMaterialBase models.RefMaterialBase
	query := r.conn
	err := query.Where("id = ?", 3).First(&refMaterialBase).Error
	if err != nil {
		return refMaterialBase, err
	}
	return refMaterialBase, nil
}

func (r *Repository) GetSettingAadtGrowthRate(RoadGrpID int) (models.SettingAadtGrowthRate, error) {
	var growthRate models.SettingAadtGrowthRate
	query := r.conn
	err := query.Where("road_group_id = ?", RoadGrpID).Where("is_latest = ?", true).First(&growthRate).Error
	if err != nil {
		return growthRate, err
	}
	return growthRate, nil
}

func (r *Repository) GetSettingAadtParameter(RoadGrpID int) (models.AadtParameterData, error) {
	var aadtParameterData models.AadtParameterData
	query := r.conn
	err := query.Table("setting_aadt_parameter").Select("*").Where("road_group_id = ?", RoadGrpID).First(&aadtParameterData).Error
	if err != nil {
		return aadtParameterData, err
	}
	return aadtParameterData, nil
}

func (r *Repository) GetLoadEquivalent(ID int) (models.RefAadtParameterVehicleType, error) {
	// data, err := u.Repo.GetLoadEquivalent(ID)
	// if err != nil {
	// 	return 0, responses.NewAppErr(400, err.Error())
	// }
	// return data, nil
	var refAadtParameterVehicleType models.RefAadtParameterVehicleType
	query := r.conn
	err := query.Where("id = ?", ID).First(&refAadtParameterVehicleType).Error
	if err != nil {
		return models.RefAadtParameterVehicleType{}, err
	}
	return refAadtParameterVehicleType, nil
}

func (r *Repository) CreateFullPrepareData(data models.PrepareData) (models.PrepareData, error) {
	query := r.conn
	var prepareData models.PrepareData
	sql := fmt.Sprintf("INSERT INTO prepare_data(the_geom,maintenance_analysis_id,is_selected,group_name,road_name,road_group_id,road_id,lane_no,lane_km_start,lane_km_end,lane_length,lane_width,km_start,km_end,length,area,last_type,type,analyst_year,year_road_begin,year_last_overlay,year_last_seal,year_last_mol_rcl,year_last_reconstruction,age,rut,iri,ifi,number_of_pothole,area_ac_icrack,percent_ac_icrack,area_ac_ucrack,percent_ac_ucrack,percent_ac_ravelling,cc_transverse_crack,cc_faulting,cc_spalling,current_surface_id,current_surface_name,current_surface_type,current_surface_surface_group,current_surface_layer_coefficient,current_surface_drainage,current_surface_a,current_surface_b,current_surface_c_base,current_surface_c_exp,current_surface_crt,current_surface_rrf,hsold,hsnew,snp_surface,snp_base,snp_subbase,snp,aadt,truck_factor,esal,yax,data,created_at,updated_at) VALUES(%s, %d, %t, '%s', '%s', %d, %d, %d, %f,%f,%f,%f,%f,%f,%f,%f, %d, %d, %d, %d, %d, %d, %d, %d, %d,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%d,'%s','%s','%s',%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,'%s','%s','%s')",
		data.TheGeom, data.MaintenanceAnalysisID, data.IsSelected, data.GroupName, data.RoadName, data.RoadGroupID, data.RoadID, data.LaneNo, data.LaneKmStart, data.LaneKmEnd, data.LaneLength, data.LaneWidth, data.KmStart, data.KmEnd, data.Length, data.Area, data.LastType, data.Type, data.AnalystYear, data.YearRoadBegin, data.YearLastOverlay, data.YearLastSeal, data.YearLastMolRcl, data.YearLastReconstruction, data.Age, data.Rut, data.Iri, data.Ifi, data.NumberOfPothole, data.AreaAcIcrack, data.PercentAcIcrack, data.AreaAcUcrack, data.PercentAcUcrack, data.PercentAcRavelling, data.CcTransverseCrack, data.CcFaulting, data.CcSpalling, data.CurrentSurfaceID, data.CurrentSurfaceName, data.CurrentSurfaceType, data.CurrentSurfaceSurfaceGroup, data.CurrentSurfaceLayerCoefficient, data.CurrentSurfaceDrainage, data.CurrentSurfaceA, data.CurrentSurfaceB, data.CurrentSurfaceCBase, data.CurrentSurfaceCExp, data.CurrentSurfaceCRT, data.CurrentSurfaceRRF, data.Hsold, data.Hsnew, data.SNPSurface, data.SNPBase, data.SNPSubbase, data.SNP, data.AADT, data.TruckFactor, data.ESAL, data.YAX, data.Data,
		data.CreatedAt.Format("2006-01-02 15:04:05"),
		data.UpdatedAt.Format("2006-01-02 15:04:05"))
	if err := query.Exec(sql).Error; err != nil {
		return models.PrepareData{}, err
	}
	return prepareData, nil
}

func (r *Repository) GetRoadGeomByRoadIDLaneNo(roadID int, laneNo int) (models.RoadGeom, error) {
	var roadGeom models.RoadGeom
	query := r.conn
	err := query.Where("road_id = ?", roadID).Where("lane_no = ?", laneNo).Where("status = ?", "A").First(&roadGeom).Error
	if err != nil {
		return roadGeom, err
	}
	return roadGeom, nil
}

func (r *Repository) GetMaxLane(roadIDs []int) (int, error) {
	var roadGeom models.RoadGeom
	query := r.conn
	err := query.Select("max(lane_no) as lane_no").Where("road_id in (?)", roadIDs).Where("status = ?", "A").Find(&roadGeom).Error
	if err != nil {
		return 0, err
	}
	return roadGeom.LaneNo, nil
}
