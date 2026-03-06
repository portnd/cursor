package repositories

import (
	"gitlab.com/mims-api-service/models"
	"gorm.io/gorm"
)

func (r *Repository) GetRoadInfoForAssetSummary(roadSectionID int) ([]models.RoadReportInfo, error) {
	var data []models.RoadReportInfo
	err := r.conn.
		Table("road").
		Select(`
		road.id AS road_id,
		road_group.number AS road_group_number,
        road_group.name AS road_group_name,
		road_section.number AS road_section_number,
		road_section.name_origin_th AS road_section_name_origin_th,
		road_section.name_destination_th AS road_section_name_destination_th,
		road_section.km_start AS km_start,
		road_section.km_end AS km_end,
		road_section.distance AS road_section_distance,
        road_info.name AS road_name,
        road.road_code,
		road_info.road_color_code,
		CASE
        WHEN road_section.name_destination_th IS NOT NULL AND road_section.name_destination_th <> '' THEN
			CONCAT(road_section.name_origin_th , ' - ' ,road_section.name_destination_th)
        ELSE
			road_section.name_destination_th
    	END AS road_main_name
    	`).
		Joins("JOIN road_group ON road_group.id = road.road_group_id").
		Joins("JOIN road_section ON road.road_section_id = road_section.id").
		Joins("JOIN road_info ON road_info.road_id = road.id and road_info.status = 'A'").
		Where("road.road_section_id = ?", roadSectionID).
		Group("road.id, road_group.number, road_group.name, road_section.number, road_section.name_origin_th, road_section.name_destination_th, road_section.km_start, road_section.km_end, road_section.distance, road_info.name, road.road_code, road_info.road_color_code").
		Order("road.id ASC").
		Find(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}
func (r *Repository) GetReportSummayAssetForAssetSummary(roadID []int) ([]models.RefSummaryAsset, error) {
	var resp []models.RefSummaryAsset
	db := r.conn
	db = db.Table("road_asset").Select("ref_asset.id as id, ref_asset.name as name").
		Joins("JOIN ref_asset_table on ref_asset_table.id = road_asset.ref_asset_table_id and ref_asset_table.is_active = true").
		Joins("JOIN ref_asset on ref_asset.id = ref_asset_table.ref_asset_id").
		Where("road_id in ?", roadID). // where in ที่หาได้จาก section
		Where("road_asset.status = ?", "A").
		Where("ref_asset_table.is_in_road = ?", true).
		Group("ref_asset.id")

	db = db.Preload("SummaryAsset", func(db *gorm.DB) *gorm.DB {
		db = db.Table("ref_asset_table").Select(` 
		ref_asset_table.ref_asset_id as ref_asset_id,
		ref_asset_table.id,
		ref_asset_table.table_label,
		array_agg(road_asset.id) as road_asset_id,
		ref_asset_table.table_name,
		MAX(road_asset.updated_date) as last_updated_date,
		MIN(road_asset.updated_date) as first_updated_date
		`).
			Joins("JOIN road_asset on ref_asset_table.id = road_asset.ref_asset_table_id and ref_asset_table.is_active = true").
			Where("road_asset.road_id in ?", roadID). // where in ที่หาได้จาก section
			Where("road_asset.status = ?", "A").
			Where("ref_asset_table.is_in_road = ?", true).
			Order("ref_asset_table.id ASC").
			Group("ref_asset_table.id")
		return db
	})

	if err := db.Find(&resp).Error; err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Repository) GetNoTypeCountForAssetSummary(roadAssetID []int32, tableName string) (*models.CountSummaryAsset, error) {
	var resp models.CountSummaryAsset

	err := r.conn.
		Table(tableName).
		Select(` 
		count(*)
		`).
		Where(tableName+".road_asset_id in ?", roadAssetID).
		Where(tableName+".is_deleted = ?", false).
		Find(&resp).Error

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r *Repository) GetTypeCountForAssetSummary(roadAssetID []int32, tableName string, joinName string) ([]models.CountAndTypeSummaryAsset, error) {
	var resp []models.CountAndTypeSummaryAsset

	selectExpr := tableName + "." + joinName + "_id," +
		joinName + ".name," +
		"COUNT(" + tableName + "." + joinName + "_id) as count"

	err := r.conn.Table(tableName).
		Select(selectExpr).
		Joins("JOIN "+joinName+" ON "+joinName+".id = "+tableName+"."+joinName+"_id").
		Where(tableName+".road_asset_id in ?", roadAssetID).
		Where(tableName+".is_deleted = ?", false).
		Group(tableName + "." + joinName + "_id, " + joinName + ".name").
		Order(tableName + "." + joinName + "_id ASC").
		Find(&resp).Error

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Repository) GetLightCountForAssetSummary(roadAssetID []int32) ([]models.CountLightSummaryAsset, error) {
	var resp []models.CountLightSummaryAsset

	err := r.conn.
		Table("road_asset_in_electricpost").
		Select(`
		road_asset_in_electricpost.ref_asset_electricpost_id, 
		ref_asset_electricpost.name as name, 
		ref_asset_tube.name as type, 
		COUNT(road_asset_in_electricpost.ref_asset_electricpost_id)
		`).
		Joins("JOIN ref_asset_electricpost ON ref_asset_electricpost.id = road_asset_in_electricpost.ref_asset_electricpost_id").
		Joins("JOIN ref_asset_tube ON ref_asset_tube.id = road_asset_in_electricpost.ref_asset_tube_id").
		Where("road_asset_in_electricpost.road_asset_id in ?", roadAssetID).
		Where("road_asset_in_electricpost.is_deleted = ? ", false).
		Group("road_asset_in_electricpost.ref_asset_electricpost_id, ref_asset_electricpost.name, ref_asset_tube.name").
		Order("road_asset_in_electricpost.ref_asset_electricpost_id ASC").
		Find(&resp).Error

	if err != nil {
		return nil, err
	}

	return resp, nil
}
