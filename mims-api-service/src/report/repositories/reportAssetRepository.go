package repositories

import (
	"strings"

	"gitlab.com/mims-api-service/models"
)

func (r *Repository) GetRoadInfo(roadID int) ([]models.RoadReportInfo, error) {
	var data []models.RoadReportInfo
	err := r.conn.Debug().
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
		road_info.road_color_code
    	`).
		Joins("JOIN road_group ON road_group.id = road.road_group_id").
		Joins("JOIN road_section ON road.road_section_id = road_section.id").
		Joins("JOIN road_info ON road_info.road_id = road.id").
		Where("road.road_section_id = ? AND road_info.status = ?", roadID, "A").
		Group("road.id, road_group.number, road_group.name, road_section.number, road_section.name_origin_th, road_section.name_destination_th, road_section.km_start, road_section.km_end, road_section.distance, road_info.name, road.road_code, road_info.road_color_code").
		Order("road.id ASC").
		Find(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Repository) GetTableName(roadID, assetID string) (*models.TableName, error) {
	var resp models.TableName

	// err := r.conn.
	// 	Table("road_asset").
	// 	Select(`DISTINCT
	// 	ref_asset_table.table_label,
	// 	ref_asset_table.table_name
	// 	`).
	// 	Joins("JOIN ref_asset_table on road_asset.ref_asset_table_id = ref_asset_table.id").
	// 	//Where("road_asset.road_id = ?", roadID).
	// 	Where("road_asset.ref_asset_table_id = ?", assetID).
	// 	//Where("road_asset.status = ?", "A").
	// 	Find(&resp).Error

	err := r.conn.
		Table("ref_asset_table").
		Select(`DISTINCT
	ref_asset_table.table_label,
	ref_asset_table.table_name
	`).
		Where("id = ?", assetID).
		Find(&resp).Error

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r *Repository) GetColumn(assetID, typ string) ([]models.Column, error) {
	var resp []models.Column
	where := "is_visible_report"
	if typ == "excel" {
		where = "is_visible_view"
	}
	err := r.conn.
		Table("ref_asset_table_columns").
		Select(`
		column_seq,
		column_name,
		component_title,
		component_type
		`).
		Where("ref_asset_table_id = ?", assetID).
		Where("ref_asset_table_columns."+where+" = ?", true).
		Order("column_seq ASC").
		Find(&resp).Error

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Repository) GetAssetName(assetID string) (*models.AssetName, error) {
	var resp models.AssetName

	err := r.conn.
		Table("ref_asset_table").
		Select(`
		table_label
		`).
		Where("id = ?", assetID).
		Scan(&resp).Error

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r *Repository) GetRow(columnName []string, roadID, assetID, tableName, join, typ string) ([]map[string]interface{}, error) {
	var resp []map[string]interface{}

	where := "is_visible_report"
	if typ == "excel" {
		where = "is_visible_view"
	}

	query := "SELECT " + "DISTINCT " + strings.Join(columnName, ",") + "," + tableName + ".id" +
		" FROM " + tableName +
		" JOIN road_asset on road_asset.id = " + tableName + ".road_asset_id" +
		" JOIN ref_asset_table_columns on ref_asset_table_columns.ref_asset_table_id = road_asset.ref_asset_table_id" +
		join +
		" WHERE road_asset.road_id = ?" +
		" AND road_asset.ref_asset_table_id = ?" +
		" AND road_asset.status = ?" +
		" AND ref_asset_table_columns." + where + " = ?" +
		" AND is_deleted = ?"
	rows, err := r.conn.Raw(query, roadID, assetID, "A", true, false).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			switch v := values[i].(type) {
			case []byte:
				row[col] = string(v)
			default:
				row[col] = v
			}
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (r *Repository) GetRoadGeom(roadID string) ([]models.MapGeom, error) {
	var resp []models.MapGeom

	err := r.conn.
		Table("road_info").
		Select(` 
		ST_AsText(the_geom) AS string_geom
		`).
		Where("road_id = ?", roadID).
		Find(&resp).Error
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Repository) GetMapGeom(roadID, assetID, tableName, typ string) ([]models.MapGeom, error) {
	var resp []models.MapGeom
	where := "is_visible_report"
	if typ == "excel" {
		where = "is_visible_view"
	}
	err := r.conn.
		Table(tableName).
		Select(`DISTINCT  
		ST_AsText(the_geom) AS string_geom
		`).
		Joins("JOIN road_asset on road_asset.id = "+tableName+".road_asset_id").
		Joins("JOIN ref_asset_table_columns on ref_asset_table_columns.ref_asset_table_id = road_asset.ref_asset_table_id").
		Where("road_asset.road_id = ?", roadID).
		Where("road_asset.ref_asset_table_id = ?", assetID).
		Where("road_asset.status = ?", "A").
		Where("ref_asset_table_columns."+where+" = ?", true).
		Find(&resp).Error

	if err != nil {
		return nil, err
	}

	return resp, nil
}
