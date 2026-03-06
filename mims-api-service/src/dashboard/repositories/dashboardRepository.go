package repositories

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/dashboard/handlers"
	"gitlab.com/mims-api-service/src/dashboard/usecases"
	"gorm.io/gorm"
)

type repository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewRepositoryHandler(conn *gorm.DB) *handlers.Handler {
	usecase := usecases.NewUsecase(&repository{conn})
	handler := handlers.NewHandler(usecase)
	return handler
}

func (r *repository) GetAssetMap() ([]models.TableResult, error) {

	tableResult := []models.TableResult{}
	selectColumn := `distinct ref_at.table_name as table_name,ref_at.table_label as table_label,ref_at.icon_filepath as icon_filepath,ref_at.line_color_code as line_color_code,
	ref_at.id as asset_id`

	query := r.conn
	query = query.Table("road AS r")
	query = query.Select(selectColumn).
		Joins("LEFT JOIN road_asset ra ON r.id = ra.road_id").
		Joins("RIGHT JOIN ref_asset_table ref_at ON ref_at.id = ra.ref_asset_table_id").
		Joins("RIGHT JOIN ref_asset ref_a ON ref_a.id = ref_at.ref_asset_id")

	query = query.Scan(&tableResult)
	err := query.Error
	if err != nil {
		return tableResult, err
	}

	return tableResult, nil
}

func (r *repository) GetTableResult(roadIDs []string, assetIDs []string, depotCodes []string, filter requests.AssetMap) ([]models.TableResult, error) {
	tableResult := []models.TableResult{}
	selectColumn := `distinct ref_at.table_name as table_name,ref_at.table_label as table_label,ref_at.icon_filepath as icon_filepath,ref_at.line_color_code as line_color_code,
	ref_at.id as asset_id,ref_at.id as id`

	query := r.conn
	query = query.Table("road AS r")
	query = query.Select(selectColumn).
		Joins("LEFT JOIN road_info ri ON r.id = ri.road_id").
		Joins("LEFT JOIN road_section rs ON r.road_section_id = rs.id").
		Joins("LEFT JOIN road_asset ra ON r.id = ra.road_id").
		Joins("RIGHT JOIN ref_asset_table ref_at ON ref_at.id = ra.ref_asset_table_id").
		Joins("RIGHT JOIN ref_asset ref_a ON ref_a.id = ref_at.ref_asset_id")
		// if len(roadID) == 0 {
		// 	query = query.Where("ref_at.is_active = true AND ra.status = 'A'  AND ref_at.id IN (?)", assetID)
		// }

	query = query.Where("ref_at.is_active = true AND ra.status = 'A' ")
	// if len(roadID) == 0 {

	// } else {
	// 	query = query.Where("ref_at.is_active = true AND ra.status = 'A' AND r.id IN (?) AND ref_at.id IN (?)", roadID, assetID)
	// }

	if filter.Year != "" {
		query = query.Where("ri.year = ?", filter.Year)
	}

	if filter.KmStart != 0.0 && filter.KmEnd != 0.0 {
		query = query.Where("ri.km_start >= ? AND ri.km_end <= ?", filter.KmStart, filter.KmEnd)
	} else if filter.KmStart != 0.0 {
		query = query.Where("ri.km_start >= ?", filter.KmStart)
	} else if filter.KmEnd != 0.0 {
		query = query.Where("ri.km_end <= ?", filter.KmEnd)
	}

	if len(depotCodes) != 0 {
		query = query.Where("rs.ref_depot_code IN ? ", depotCodes)
	}

	if len(roadIDs) != 0 {
		query.Where("r.id IN (?)", roadIDs)
	}

	if len(assetIDs) != 0 {
		query.Where("ref_at.id IN (?)", assetIDs)
	}

	query = query.Find(&tableResult)
	err := query.Error
	if err != nil {
		return tableResult, err
	}

	return tableResult, nil
}

func (r *repository) GetRoadAssetDetails(ID int, data []responses.RoadAssetDetailColumn) ([]map[string]interface{}, error) {

	results := []map[string]interface{}{}

	if len(data) == 0 {
		return results, nil
	}
	query := r.conn

	sqlJoin := ""
	sqlSelect := "c.id as asset_object_id, ra.id as id, ra.id_parent AS id_parent, ra.status AS status, ref_ds.name AS status_text, ra.updated_by AS updated_by, ra.updated_date AS updated_date, ra.reject_reason AS reject_reason, ra.revision AS revision, ra.is_exclusive_lock AS is_exclusive_lock, c.id_parent AS id_parent_asset"
	refTableCount := 1
	sqlOrders := []string{}
	sqlOrders = append(sqlOrders, "id ASC")
	for _, item := range data {
		if item.TableNameRef != "" {
			tableAlias := fmt.Sprintf("r%d", refTableCount)
			sqlJoin += " " + fmt.Sprintf("LEFT JOIN %s as %s on %s.id = c.%s_id", item.TableNameRef, tableAlias, tableAlias, item.TableNameRef)

			sqlSelect += ", " + fmt.Sprintf("%s.name as %s_name", tableAlias, item.TableNameRef)
			if strings.Contains(item.TableNameRef, "_image") {
				pathCol := strings.Replace(item.TableNameRef, "ref_asset_", "", 1) + "_filepath"
				sqlSelect += ", " + fmt.Sprintf("%s.%s as %s_filepath", tableAlias, pathCol, item.TableNameRef)
				if item.TableNameRef == "ref_asset_sign_image" {
					sqlSelect += ", " + fmt.Sprintf("%s.abbr as %s_abbr", tableAlias, item.TableNameRef)
				}
			}
			refTableCount++
		}

		if strings.Contains(item.ColumnName, "geom_camera") {
			sqlSelect += ", " + fmt.Sprintf("ST_ASTEXT(c.%s) AS %s", item.ColumnName, item.ColumnName)
		} else if strings.Contains(item.ColumnName, "geom") {
			sqlSelect += ", " + fmt.Sprintf("ST_ASTEXT(ST_FORCE2D(c.%s)) AS %s_2d", item.ColumnName, item.ColumnName)
			sqlSelect += ", " + fmt.Sprintf("ST_ASTEXT(ST_FORCE3D(c.%s)) AS %s_3d", item.ColumnName, item.ColumnName)
		} else {
			sqlSelect += ", " + fmt.Sprintf("c.%s as %s", item.ColumnName, item.ColumnName)
		}

		if item.ColumnName == "km" {
			sqlOrders = append(sqlOrders, "km ASC")
		} else if item.ColumnName == "km_start" {
			sqlOrders = append(sqlOrders, "km_start ASC")
		} else if item.ColumnName == "latitude" {
			sqlOrders = append(sqlOrders, "latitude ASC")
		}
	}
	sqlWhere := ""

	sqlWhere += " where  ra.status = 'A' "

	sqlWhere += fmt.Sprintf("and c.is_deleted = false and c.id = %d", ID) // input

	sqlOrders = append(sqlOrders, "revision DESC")
	order := " order by " + helpers.Implode(",", sqlOrders)
	sqlOffset := ""

	sql := "select " + sqlSelect + " from road_asset as ra JOIN " + data[0].TableName + " as c on ra.id = c.road_asset_id LEFT JOIN ref_data_status as ref_ds on ref_ds.status_code = ra.status " + sqlJoin + sqlWhere + order + sqlOffset

	rows, err := query.Raw(sql).Rows()
	if err != nil {
		return results, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return results, err
	}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return results, err
		}

		result := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			result[colName] = *val
		}

		results = append(results, result)
	}
	return results, nil

}

func (r *repository) GetRoadAssetSignDataByRoadIDs(t models.TableResult, roadID []string, filter requests.AssetMap) ([]models.RoadAssetSign, error) {
	var result []models.RoadAssetSign
	table := t.TableName
	var selectColumn string
	check, err := r.CheckSignImgId(t.AssetId)
	if err != nil {
		return result, err
	}
	leftJoin := ""
	if check {
		selectColumn = "s.sign_image_filepath,"
		leftJoin = "LEFT JOIN ref_asset_sign_image s ON s.id = t.ref_asset_sign_image_id"
	}
	cols, err := r.FindImgFilepath(t.AssetId)
	if err != nil {
		return result, err
	}
	for i, col := range cols {
		if i != len(cols)-1 {
			selectColumn = selectColumn + " " + col + " " + "AS img_filepath,"
			continue
		}
		selectColumn = selectColumn + " " + col + " " + "AS img_filepath"
	}

	if leftJoin == "" && len(cols) == 0 {
		selectColumn = `ST_AsGeoJSON(t.the_geom) as the_geom` + selectColumn
	} else {
		selectColumn = `ST_AsGeoJSON(t.the_geom) as the_geom` + ", " + selectColumn

	}
	// if leftJoin == "" && len(cols) == 0 {
	// 	selectColumn = `ST_ASTEXT(ST_Force2D(t.the_geom)) as geom_cl,ST_X(ST_Centroid(ST_Transform(ST_SetSRID(t.the_geom, 4326), 4326))) as lon,ST_Y(ST_Centroid(ST_Transform(ST_SetSRID(t.the_geom, 4326), 4326))) as lat` + selectColumn
	// } else {
	// 	selectColumn = `ST_ASTEXT(ST_Force2D(t.the_geom)) as geom_cl,ST_X(ST_Centroid(ST_Transform(ST_SetSRID(t.the_geom, 4326), 4326))) as lon,ST_Y(ST_Centroid(ST_Transform(ST_SetSRID(t.the_geom, 4326), 4326))) as lat` + ", " + selectColumn

	// }
	selectColumn = strings.TrimRight(selectColumn, ",")

	query := r.conn
	query = query.Table(table+" AS t").Select(selectColumn, "t.id, ra.road_id").
		Joins("join road_asset ra on ra.id = t.road_asset_id").
		Joins(leftJoin)
	if len(roadID) > 0 {
		query = query.Where("ra.road_id in (?)", roadID)
	}
	query = query.Where("ra.status = 'A' and t.is_deleted = false")
	if filter.Left != "" && filter.Bottom != "" && filter.Right != "" && filter.Top != "" {
		query = query.Where("the_geom && ST_MakeEnvelope(?, ?, ?, ?, 4326)", filter.Left, filter.Bottom, filter.Right, filter.Top)
	}
	query = query.Scan(&result)
	err = query.Error
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *repository) GetRoadAssetTheGeomCuster(buildQuery string, roadIDs []string, assetIDs []string, depotCodes []string, filter requests.AssetMap) ([]models.RoadAssetGeomCuster, error) {

	var result []models.RoadAssetGeomCuster

	zoomInit := ""
	zoom, err := strconv.Atoi(filter.Zoom)
	if err != nil {
		zoomInit = "0.01"
	}

	if zoom <= 8 {
		zoomInit = "1"
	} else if zoom <= 10 {
		zoomInit = ".1"
	} else {
		zoomInit = ".01"
	}

	// Construct the raw SQL query
	sql := fmt.Sprintf(`SELECT centriod.count_the_geom AS total_the_geom_cluster,
    ST_AsGeoJSON(ST_Centroid(centriod.geom)) AS the_geom_cluster
from (
        SELECT cluster.cid,
            ST_Union(cluster.the_geom) as geom,
            COUNT(*) as count_the_geom
        FROM (
                SELECT ST_ClusterKMeans(geom.the_geom, 1, %s ) OVER () AS cid,
                    geom.the_geom
                FROM (
					%s 
                    ) AS geom
            ) AS cluster
        GROUP BY cid
    ) AS centriod;`, zoomInit, buildQuery)

	// Executing the raw SQL query using GORM's Raw method and scanning the results into the results slice
	query := r.conn.Raw(sql)

	if filter.Year != "" {
		query = query.Where("ri.year = ?", filter.Year)
	}

	if filter.KmStart != 0.0 && filter.KmEnd != 0.0 {
		query = query.Where("ri.km_start >= ? AND ri.km_end <= ?", filter.KmStart, filter.KmEnd)
	} else if filter.KmStart != 0.0 {
		query = query.Where("ri.km_start >= ?", filter.KmStart)
	} else if filter.KmEnd != 0.0 {
		query = query.Where("ri.km_end <= ?", filter.KmEnd)
	}

	if len(depotCodes) != 0 {
		query = query.Where("rs.ref_depot_code IN ? ", depotCodes)
	}

	if len(roadIDs) != 0 {
		query.Where("r.id IN (?)", roadIDs)
	}

	if len(assetIDs) != 0 {
		query.Where("ra.ref_asset_table_id IN (?)", assetIDs)
	}

	err = query.Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (r *repository) GetRoadConditionGradesByID(ownerID int, conditionType string) ([]models.ParamsConditionPreload, error) {
	var paramsCondition []models.ParamsConditionPreload
	query := r.conn
	query = query.Preload("RefOwner")
	query = query.Preload("RefGrade")
	query = query.Where("ref_owner_id = ? AND LOWER(condition_type) = ? ", ownerID, conditionType)
	err := query.Find(&paramsCondition).Error
	if err != nil {
		return paramsCondition, err
	}
	sort.Slice(paramsCondition, func(i, j int) bool {
		return paramsCondition[i].RefGradeID < paramsCondition[j].RefGradeID
	})
	return paramsCondition, nil
}

func (r *repository) GetRefGrade() ([]models.RefGrade, error) {
	var refGrades []models.RefGrade
	err := r.conn.Find(&refGrades).Error
	if err != nil {
		return refGrades, err
	}

	return refGrades, nil
}

func (r *repository) GetRoadLineGradesByID(ownerID int) ([]models.ParamsRoadLinePreload, error) {
	var paramsRoadLine []models.ParamsRoadLinePreload
	query := r.conn
	query = query.Preload("RefOwnerRoadLine")
	query = query.Preload("RefGrade")
	query = query.Where("ref_owner_road_line_id = ?", ownerID)
	err := query.Find(&paramsRoadLine).Error
	if err != nil {
		return paramsRoadLine, err
	}
	sort.Slice(paramsRoadLine, func(i, j int) bool {
		return paramsRoadLine[i].RefGradeID < paramsRoadLine[j].RefGradeID
	})
	return paramsRoadLine, nil
}

// func (r *repository) GetGeomJsonFromTableID(maintenanceRoadHisID int) ([]byte, error) {
// 	type Result struct {
// 		TheGeom []byte `gorm:"column:the_geom"`
// 	}

// 	result := Result{}

// 	err := r.conn.Table("maintenance_road_history").Select("ST_AsGeoJSON(the_geom) AS the_geom").Where("id = ? and status = ?", maintenanceRoadHisID, "A").First(&result).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result.TheGeom, nil
// }

// func (r *repository) GetRoadAssetSignDataByRoadIDs(t models.TableResult, roadID []int) ([]models.RoadAssetSign, error) {
// 	var result []models.RoadAssetSign
// 	table := t.TableName
// 	var selectColumn string
// 	check, err := r.CheckSignImgId(t.AssetId)
// 	if err != nil {
// 		return result, err
// 	}
// 	leftJoin := ""
// 	if check {
// 		selectColumn = "s.sign_image_filepath,"
// 		leftJoin = "LEFT JOIN ref_asset_sign_image s ON s.id = t.ref_asset_sign_image_id"
// 	}
// 	cols, err := r.FindImgFilepath(t.AssetId)
// 	if err != nil {
// 		return result, err
// 	}
// 	for i, col := range cols {
// 		if i != len(cols)-1 {
// 			selectColumn = selectColumn + " " + col + " " + "AS img_filepath,"
// 			continue
// 		}
// 		selectColumn = selectColumn + " " + col + " " + "AS img_filepath"
// 	}
// 	if leftJoin == "" && len(cols) == 0 {
// 		selectColumn = `ST_ASTEXT(ST_Force2D(t.the_geom)) as geom_cl,ST_X(ST_Centroid(ST_Transform(ST_SetSRID(t.the_geom, 4326), 4326))) as lon,ST_Y(ST_Centroid(ST_Transform(ST_SetSRID(t.the_geom, 4326), 4326))) as lat` + selectColumn
// 	} else {
// 		selectColumn = `ST_ASTEXT(ST_Force2D(t.the_geom)) as geom_cl,ST_X(ST_Centroid(ST_Transform(ST_SetSRID(t.the_geom, 4326), 4326))) as lon,ST_Y(ST_Centroid(ST_Transform(ST_SetSRID(t.the_geom, 4326), 4326))) as lat` + ", " + selectColumn

// 	}
// 	selectColumn = strings.TrimRight(selectColumn, ",")

// 	query := r.conn
// 	query = query.Table(table + " AS t").Select(selectColumn).
// 		Joins("join road_asset ra on ra.id = t.road_asset_id").
// 		Joins(leftJoin)
// 	if len(roadID) > 0 {
// 		query = query.Where("ra.road_id in (?)", roadID)
// 	}
// 	query = query.Where("ra.status = 'A' and t.is_deleted = false")
// 	query = query.Scan(&result)

// 	return result, nil
// }

func (r *repository) GetRoadCondition(where string) ([]models.RoadConditionSurveyM2, error) {
	var model []models.RoadConditionSurveyM2
	query := r.conn

	if err := query.Debug().Select("*, st_asgeojson(the_geom) as geojson").Find(&model).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (r *repository) GetRoadCondition100M(where string) ([]models.RoadConditionSurvey100M2, error) {
	var model []models.RoadConditionSurvey100M2
	query := r.conn

	if err := query.Debug().Select("*, st_asgeojson(the_geom) as geojson").Find(&model).Error; err != nil {
		return model, err
	}
	return model, nil
}

func (r *repository) CheckSignImgId(assetID int) (bool, error) {
	count := int64(0)
	if assetID != 0 {
		err := r.conn.Table("ref_asset_table_columns").Where("ref_asset_table_id = (?) and column_name = 'ref_asset_sign_image_id'", assetID).Count(&count).Error
		if err != nil {
			return false, err
		}
		if count > 0 {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, nil
}

func (r *repository) FindImgFilepath(assetID int) ([]string, error) {
	type columnName struct {
		ColumnName string
	}
	var listColumn []columnName
	var result []string

	cols := r.conn

	err := cols.Table("ref_asset_table_columns AS r").Select("column_name").Where("ref_asset_table_id = (?) and column_name LIKE '%_filepath'", assetID).Scan(&listColumn).Error

	if err != nil {
		return result, err
	}

	for _, colName := range listColumn {
		result = append(result, colName.ColumnName)
	}

	return result, nil
}

func (r *repository) GetRefDataFromSelect(tableName string, ID int64) (string, error) {
	query := r.conn

	var result string
	err := query.Table(tableName).Select("name").Where("id = ?", ID).Scan(&result).Error
	if err != nil {
		return "", err
	}

	return result, nil
}

func (r *repository) GetRoadDashboard() ([]responses.RoadGroupDashboard, error) {

	var resp []responses.RoadGroupDashboard

	sql := `SELECT 
    COUNT(road.id) AS total_road, 
    CONCAT('', ltrim(road_group.number, '0')) as road_group_name , 
    road_group.distance,
	road_group.id
FROM 
    public.road
LEFT JOIN 
    road_group ON road.road_group_id = road_group.id
GROUP BY 
    road_group."number", road_group.distance, road_group.id
	ORDER BY 
    road_group.id ASC;
	`
	query := r.conn
	if err := query.Raw(sql).Scan(&resp).Error; err != nil {
		return resp, nil
	}

	return resp, nil
}

func (r *repository) GetPavementSurfaceByRoadGroupID(roadGroupID int) ([]models.PavementSurface, error) {
	var pavementSurface []models.PavementSurface
	query := r.conn
	query = query.Table("road").Select("road.id, sum(abs(road_info.km_end - road_info.km_start)/1000) as length").
		Joins("INNER JOIN road_info on road_info.road_id = road.id and road_info.status = 'A'")

	query = query.Preload("SurfaceData", func(db *gorm.DB) *gorm.DB {
		db = db.Select("road_surface.id ,road_surface.road_id as road_id, sum(abs(road_surface.km_end - road_surface.km_start)/1000) as length")

		db = db.Where("road_surface.status = 'A'").Group("road_surface.id")
		return db
	})

	query = query.Preload("SurfaceData.SurfaceLane", func(db *gorm.DB) *gorm.DB {
		db = db.Select("road_surface_lane.id ,road_surface_lane.road_surface_id ,lane_no ,ref_surface.surface_group as surface_tyepe").
			Joins("LEFT JOIN ref_surface ON road_surface_lane.ref_surface_id = ref_surface.id ")
		return db
	})

	query = query.Where("road.is_active = true and road.road_group_id = ?", roadGroupID)
	query = query.Group("road.id")
	if err := query.Find(&pavementSurface).Error; err != nil {
		return pavementSurface, err
	}

	return pavementSurface, nil
}

func (r *repository) GetVolumeAADTByRoadGroupID(roadGroupID int) ([]responses.VolumeAADTDashboard, error) {

	var volumeAADT []responses.VolumeAADTDashboard
	sql := ` WITH aadt AS (
		SELECT
			road.id AS road_id,
			volume_aadt.total,
			volume_aadt.year,
			volume_aadt.created_date,
			road_group.number AS road_group_number,
			ROW_NUMBER() OVER (
				PARTITION BY road.id, volume_aadt.year
				ORDER BY volume_aadt.created_date DESC
			) AS row_num
		FROM
			public.road
		LEFT JOIN
			volume_aadt ON road.id = volume_aadt.road_id
		LEFT JOIN
			road_group ON road.road_group_id = road_group.id
		WHERE
			road.road_group_id = ?
			AND volume_aadt.status = 'A'
			AND volume_aadt.year IN (
				EXTRACT(YEAR FROM (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Bangkok')),
				EXTRACT(YEAR FROM (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Bangkok') - interval '1 year')
			)
	)
	
	SELECT
		aadt.year,
		CONCAT('', LTRIM(aadt.road_group_number, '0')) AS road_group_name,
		SUM(aadt.total) AS total
	FROM
		aadt
	WHERE
		aadt.row_num = 1
	GROUP BY
		aadt.year,
		aadt.road_group_number
	ORDER BY
		aadt.year ASC;
	`
	query := r.conn
	if err := query.Raw(sql, roadGroupID).Scan(&volumeAADT).Error; err != nil {
		return volumeAADT, err
	}

	return volumeAADT, nil
}

func (r *repository) GetRoadGroups() ([]models.RoadGroup, error) {
	var roadGroups []models.RoadGroup
	if err := r.conn.Find(&roadGroups).Error; err != nil {
		return roadGroups, err
	}
	return roadGroups, nil
}

// func (r *repository) GetAssetRawData(roadIDs []int, filter requests.Asset) ([]models.RoadAssetSummary, error) {

// 	query := r.conn
// 	var roadAssetSummaries []models.RoadAssetSummary
// 	query
// 	query = query.Preload("RefAssetTable",
// 		func(query *gorm.DB) *gorm.DB {
// 			query = query.Select("id,ref_asset_id, seq, is_in_road, table_label, icon_filepath, line_color_code, geom_type ")
// 			query = query.Where("is_active = ?", true).Where("is_in_road = ?", true)
// 			return query.Order("is_in_road, seq")
// 		})
// 	query = query.Preload("RefAssetTable.RoadAsset", func(db *gorm.DB) *gorm.DB {
// 		return db.Where("status = 'A'").Select("id,road_id,ref_asset_table_id ,get_asset_count(ref_asset_table_id, id) AS asset_count")
// 	})
// 	query = query.Select("ID, Name")
// 	// query = query.Order("CASE WHEN(id = 7) THEN 999 ELSE id END")
// 	err := query.Find(&roadAssetSummaries).Error
// 	if err != nil {
// 		return roadAssetSummaries, err
// 	}

// 	return roadAssetSummaries, nil
// }

// func (r *repository) GetAsset(roadIDs []int, filter requests.Asset) ([]models.RefAssetDashboard, int64, error) {

// 	query := r.conn

// 	// Default values for pagination
// 	if filter.Page < 1 {
// 		filter.Page = 1
// 	}
// 	if filter.Limit <= 0 {
// 		filter.Limit = 10 // Default limit
// 	}
// 	offset := (filter.Page - 1) * filter.Limit

// 	var refAssets []models.RefAssetDashboard
// 	var totalCount int64

// 	query = query.Model(&models.RefAssetDashboard{})

// 	query = query.Preload("RefAssetTables", func(db *gorm.DB) *gorm.DB {
// 		return db.Select("id, ref_asset_id, seq, is_in_road, table_label, icon_filepath, line_color_code, geom_type").
// 			Where("is_active = ?", true).Where("is_in_road = ?", true).
// 			Order("is_in_road, seq")
// 	})

// 	query = query.Preload("RefAssetTables.RoadAssets", func(db *gorm.DB) *gorm.DB {
// 		return db.Where("status = 'A'").Select("id,road_id,ref_asset_table_id ,get_asset_count(ref_asset_table_id, id) AS asset_count")
// 	})

// 	query = query.Preload("RefAssetTables.RoadAssets.Road", func(db *gorm.DB) *gorm.DB {

// 		if len(roadIDs) != 0 {
// 			db = db.Where("id IN ?", roadIDs)
// 		}
// 		return db.Select("id, road_group_id, seq, road_code, is_active").Where("is_active = ?", true)
// 	})

// 	query = query.Preload("RefAssetTables.RoadAssets.Road.RoadGroup", func(db *gorm.DB) *gorm.DB {
// 		return db.Select("id, number")
// 	})

// 	query = query.Select("id,name")

// 	query = query.Where("status = ?", 1)

// 	if err := query.Count(&totalCount).Error; err != nil {
// 		log.Println("Error counting total records", err)
// 		return nil, 0, err
// 	}

// 	if err := query.Offset(offset).Limit(filter.Limit).Find(&refAssets).Error; err != nil {
// 		log.Println("Error offset", err)
// 		return nil, 0, err
// 	}

// 	return refAssets, totalCount, nil
// }

func (r *repository) GetAsset(roadIDs []int, depotCodes []string, filter requests.Asset) ([]models.RefAssetDashboard, int64, error) {
	query := r.conn

	// Default values for pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 50 // Default limit
	}
	offset := (filter.Page - 1) * filter.Limit

	var refAssets []models.RefAssetDashboard
	var totalCount int64

	baseQuery := query.Model(&models.RefAssetDashboard{}).
		Where("EXISTS (SELECT 1 FROM ref_asset_table WHERE ref_asset_table.ref_asset_id = ref_asset.id)").
		Where("status = ?", 1)

	if err := baseQuery.Count(&totalCount).Error; err != nil {
		log.Println("Error counting total records with RefAssetTables", err)
		return nil, 0, err
	}

	resultQuery := baseQuery.
		Preload("RefAssetTables", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, ref_asset_id, seq, is_in_road, table_label, icon_filepath, line_color_code, geom_type").
				Where("is_active = ?", true).Where("is_in_road = ?", true).
				Order("is_in_road, seq")
		}).
		Preload("RefAssetTables.RoadAssets", func(db *gorm.DB) *gorm.DB {

			if len(roadIDs) > 0 {
				db = db.Where("road_id IN ?", roadIDs)
			}

			return db.Where("status = 'A'").
				Select("id,road_id,ref_asset_table_id, get_asset_count(ref_asset_table_id, id) AS asset_count")
		}).
		Preload("RefAssetTables.RoadAssets.Road", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, road_group_id,road_section_id, seq, road_code, is_active").Where("is_active = ?", true)
		}).
		Preload("RefAssetTables.RoadAssets.RoadInfo", func(db *gorm.DB) *gorm.DB {
			if filter.Year != "" {
				db = db.Where("year = ?", filter.Year)
			}

			if filter.KmStart != 0.0 && filter.KmEnd != 0.0 {
				db = db.Where("km_start >= ? AND km_end <= ?", filter.KmStart, filter.KmEnd)
			} else if filter.KmStart != 0.0 {
				db = db.Where("km_start >= ?", filter.KmStart)
			} else if filter.KmEnd != 0.0 {
				db = db.Where("km_end <= ?", filter.KmEnd)
			}

			return db
		}).
		Preload("RefAssetTables.RoadAssets.Road.RoadSection", func(db *gorm.DB) *gorm.DB {
			if len(depotCodes) > 0 {
				return db.Where("ref_depot_code IN (?)", depotCodes)
			}
			return db
		}).
		Preload("RefAssetTables.RoadAssets.Road.RoadGroup", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, CONCAT('หมายเลข ', ltrim(number, '0')) AS number")
		}).
		Select("id, name").
		Offset(offset).Limit(filter.Limit)

	if err := resultQuery.Find(&refAssets).Error; err != nil {
		log.Println("Error fetching records with offset", err)
		return nil, 0, err
	}

	return refAssets, totalCount, nil
}

func (r *repository) GetDashboardMaxMinYear() (responses.DashboardYearMaxMin, error) {
	var year responses.DashboardYearMaxMin

	sql := ` WITH combined AS (
			SELECT 
				MIN(year) AS min_year,
				MAX(year) AS max_year
			FROM 
				public.road_info
			WHERE 
				status = 'A'
			
			UNION ALL
			
			SELECT 
				MIN(year) AS min_year,
				MAX(year) AS max_year
			FROM 
				public.road_condition
			WHERE 
				status = 'A'
			
			UNION ALL
			
			SELECT 
				MIN(year) AS min_year,
				MAX(year) AS max_year
			FROM 
				public.road_surface
			WHERE 
				status = 'A'
			
			UNION ALL
			
			SELECT 
				MIN(budget_year) AS min_year,
				MAX(budget_year) AS max_year
			FROM 
				public.maintenance
			WHERE 
				status = 'A'
		)
		
		SELECT 
			MIN(min_year) AS overall_min_year,
			MAX(max_year) AS overall_max_year
		FROM 
			combined;
		`

	query := r.conn
	if err := query.Raw(sql).Scan(&year).Error; err != nil {
		return year, err
	}

	return year, nil

}

func (r *repository) GetDashboardYear(tableName string) ([]int, error) {
	var year []int

	query := r.conn

	switch tableName {
	case "road_info":
		query = query.Table("road_info").Select("distinct year")
	case "road_condition":
		query = query.Table("road_condition").Select("distinct year")
	case "road_surface":
		query = query.Table("road_surface").Select("distinct year")
	case "maintenance":
		query = query.Table("maintenance").Select("distinct budget_year as year")
	}

	err := query.Where("status = 'A'").Find(&year).Error
	if err != nil {
		return year, err
	}

	return year, nil

}

func (r *repository) GetRefAssetTableColumns(refAssetTableID int) ([]responses.RoadAssetDetailColumn, error) {
	var roadAssetDetailColumn []responses.RoadAssetDetailColumn
	query := r.conn
	sqlSelects := []string{"DISTINCT ref_at.icon_filepath as icon_filepath", "ref_at.line_color_code as line_color_code", "ref_at.table_name as table_name", "ref_atc.column_name as column_name", "ref_atc.table_name_ref as table_name_ref", "ref_atc.column_data_type as column_data_type", "ref_atc.component_title as component_title", "ref_atc.component_type as component_type", "ref_atc.column_seq as seq,ref_at.id as ref_asset_id"}
	query = query.Table("road_asset as ra")
	query = query.Select(helpers.Implode(",", sqlSelects))
	query = query.Joins("LEFT JOIN ref_asset_table as ref_at on ra.ref_asset_table_id = ref_at.id")
	query = query.Joins("LEFT JOIN ref_asset_table_columns as ref_atc on ref_at.id = ref_atc.ref_asset_table_id")

	query = query.Where("ref_at.id = ?", refAssetTableID)
	query = query.Where("ref_atc.is_visible_view = ?", true)
	orWheres := []string{}

	queryWhereString := strings.Join(orWheres, " or ")
	query = query.Where(queryWhereString)

	orWheres2 := []string{}

	queryWhereString2 := strings.Join(orWheres2, " or ")
	query = query.Where(queryWhereString2)
	query = query.Order("seq")
	if err := query.Find(&roadAssetDetailColumn).Error; err != nil {
		return roadAssetDetailColumn, err
	}
	return roadAssetDetailColumn, nil
}

func (t *repository) GetRoadAssetDetailColumn(roadAssetID int, refAssetTableID int) ([]responses.RoadAssetDetailColumn, error) {
	var roadAssetDetailColumn []responses.RoadAssetDetailColumn
	query := t.conn
	sqlSelects := []string{"DISTINCT ref_at.icon_filepath as icon_filepath", "ref_at.line_color_code as line_color_code", "ref_at.table_name as table_name", "ref_atc.column_name as column_name", "ref_atc.table_name_ref as table_name_ref", "ref_atc.column_data_type as column_data_type", "ref_atc.component_title as component_title", "ref_atc.component_type as component_type", "ref_atc.column_seq as seq"}
	query = query.Table("road_asset as ra")
	query = query.Select(helpers.Implode(",", sqlSelects))
	query = query.Joins("LEFT JOIN ref_asset_table as ref_at on ra.ref_asset_table_id = ref_at.id")
	query = query.Joins("LEFT JOIN ref_asset_table_columns as ref_atc on ref_at.id = ref_atc.ref_asset_table_id")

	query = query.Where("ra.id = ?", roadAssetID)
	query = query.Where("ref_at.id = ?", refAssetTableID)
	query = query.Where("ref_atc.is_visible_view = ?", true)
	orWheres := []string{}

	queryWhereString := strings.Join(orWheres, " or ")
	query = query.Where(queryWhereString)

	orWheres2 := []string{}

	queryWhereString2 := strings.Join(orWheres2, " or ")
	query = query.Where(queryWhereString2)
	query = query.Order("seq")
	if err := query.Find(&roadAssetDetailColumn).Error; err != nil {
		return roadAssetDetailColumn, err
	}
	return roadAssetDetailColumn, nil
}

func (r *repository) GetRoadConditionDashboard(roadIDs []string, depotCodes []string, filter requests.Condition) ([]models.RoadConditionDashboard, error) {
	var roadCondition []models.RoadConditionDashboard

	query := r.conn.Model(&models.RoadConditionDashboard{})

	query = query.Preload("Road.RoadSection", func(db *gorm.DB) *gorm.DB {
		if len(depotCodes) > 0 {
			return db.Where("ref_depot_code IN (?)", depotCodes)
		}
		return db
	})

	query = query.Preload("RoadConditionSurveys")
	query = query.Preload("RoadConditionSurveys.RoadConditionSurvey100Ms")
	query = query.Preload("RoadConditionSurveys.RoadConditionSurvey100Ms.RoadConditionSurveyMs")

	// Applying filters
	if len(roadIDs) > 0 {
		query = query.Where("road_condition.road_id IN ?", roadIDs)
	}
	// Filter by depot_code via JOIN — depot filter in Preload only filters nested data, not the main rows
	if len(depotCodes) > 0 {
		query = query.
			Joins("JOIN road AS rcd_road ON road_condition.road_id = rcd_road.id").
			Joins("JOIN road_section AS rcd_road_section ON rcd_road.road_section_id = rcd_road_section.id").
			Where("rcd_road_section.ref_depot_code IN ?", depotCodes)
	}
	if filter.Year != 0 {
		query = query.Where("road_condition.year = ?", filter.Year)
	}
	if filter.KmStart != 0.0 && filter.KmEnd != 0.0 {
		query = query.Where("road_condition.km_start >= ? AND road_condition.km_end <= ?", filter.KmStart, filter.KmEnd)
	} else if filter.KmStart != 0.0 {
		query = query.Where("road_condition.km_start >= ?", filter.KmStart)
	} else if filter.KmEnd != 0.0 {
		query = query.Where("road_condition.km_end <= ?", filter.KmEnd)
	}
	query = query.Where("road_condition.status = ?", "A")

	// Subquery to find the latest road condition per road ID based on created_date
	subquery := r.conn.Table("road_condition").
		Select("road_id, MAX(created_date) AS latest_date").
		Where("status = ?", "A").
		Group("road_id,lane_no")
	if filter.Year != 0 {
		subquery = subquery.Where("year = ?", filter.Year)
	}

	// Main query with subquery to get the latest records
	query = query.Joins("JOIN (?) AS latest ON road_condition.road_id = latest.road_id AND road_condition.created_date = latest.latest_date", subquery)

	// Execute the query
	err := query.Find(&roadCondition).Error
	if err != nil {
		return nil, err
	}

	return roadCondition, nil
}

func (r *repository) GetRoadRetroReflectivityDashboard(roadIDs []string, depotCodes []string, filter requests.Condition) ([]models.RoadRetroReflectivityDashboard, error) {
	var retroReflectivity []models.RoadRetroReflectivityDashboard

	query := r.conn.Model(&models.RoadRetroReflectivityDashboard{})

	query = query.Preload("Road.RoadSection", func(db *gorm.DB) *gorm.DB {
		if len(depotCodes) > 0 {
			return db.Where("ref_depot_code IN (?)", depotCodes)
		}
		return db
	})

	query = query.Preload("RoadRetroReflectivityRanges")
	query = query.Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs")

	// Applying filters
	if len(roadIDs) > 0 {
		query = query.Where("road_retro_reflectivity.road_id IN ?", roadIDs)
	}
	// Filter by depot_code via JOIN — depot filter in Preload only filters nested data, not the main rows
	if len(depotCodes) > 0 {
		query = query.
			Joins("JOIN road AS rrd_road ON road_retro_reflectivity.road_id = rrd_road.id").
			Joins("JOIN road_section AS rrd_road_section ON rrd_road.road_section_id = rrd_road_section.id").
			Where("rrd_road_section.ref_depot_code IN ?", depotCodes)
	}
	if filter.Year != 0 {
		query = query.Where("road_retro_reflectivity.year = ?", filter.Year)
	}
	if filter.KmStart != 0.0 && filter.KmEnd != 0.0 {
		query = query.Where("road_retro_reflectivity.km_start >= ? AND road_retro_reflectivity.km_end <= ?", filter.KmStart, filter.KmEnd)
	} else if filter.KmStart != 0.0 {
		query = query.Where("road_retro_reflectivity.km_start >= ?", filter.KmStart)
	} else if filter.KmEnd != 0.0 {
		query = query.Where("road_retro_reflectivity.km_end <= ?", filter.KmEnd)
	}
	query = query.Where("road_retro_reflectivity.status = ?", "A")

	// Subquery to find the latest road condition per road ID based on created_date
	subquery := r.conn.Table("road_retro_reflectivity").
		Select("road_id, MAX(created_date) AS latest_date").
		Where("status = ?", "A").
		Group("road_id,line_no")
	if filter.Year != 0 {
		subquery = subquery.Where("year = ?", filter.Year)
	}

	// Main query with subquery to get the latest records
	query = query.Joins("JOIN (?) AS latest ON road_retro_reflectivity.road_id = latest.road_id AND road_retro_reflectivity.created_date = latest.latest_date", subquery)

	// Execute the query
	err := query.Find(&retroReflectivity).Error
	if err != nil {
		return nil, err
	}

	return retroReflectivity, nil
}

func (r *repository) GetRoadConditionDashboardMap(roadIDs []string, depotCodes []string, filter requests.ConditionMap) ([]models.RoadConditionDashboard, error) {
	var roadCondition []models.RoadConditionDashboard

	query := r.conn.Debug().Model(&models.RoadConditionDashboard{})

	query = query.Preload("Road.RoadSection", func(db *gorm.DB) *gorm.DB {
		if len(depotCodes) > 0 {
			return db.Where("ref_depot_code IN (?)", depotCodes)
		}
		return db
	})

	// if filter.Left != "" && filter.Bottom != "" && filter.Right != "" && filter.Top != "" {
	// 	query = query.Where("the_geom && ST_MakeEnvelope(?, ?, ?, ?, 4326)", filter.Left, filter.Bottom, filter.Right, filter.Top)
	// }

	query = query.Preload("RoadConditionSurveys", func(db *gorm.DB) *gorm.DB {

		if filter.Left != "" && filter.Bottom != "" && filter.Right != "" && filter.Top != "" {
			return db.Select("*, ST_AsGeoJSON(ST_ASTEXT(ST_FORCE2D(the_geom))) AS survey_the_geom").Where("the_geom && ST_MakeEnvelope(?, ?, ?, ?, 4326)", filter.Left, filter.Bottom, filter.Right, filter.Top)
		}

		return db.Select("*, ST_AsGeoJSON(ST_ASTEXT(ST_FORCE2D(the_geom))) as survey_the_geom ")
	})

	query = query.Preload("RoadConditionSurveys.RoadConditionSurvey100Ms", func(db *gorm.DB) *gorm.DB {

		if filter.Left != "" && filter.Bottom != "" && filter.Right != "" && filter.Top != "" {
			return db.Select("*, ST_AsGeoJSON(ST_ASTEXT(ST_FORCE2D(the_geom))) AS survey_100m_the_geom").Where("the_geom && ST_MakeEnvelope(?, ?, ?, ?, 4326)", filter.Left, filter.Bottom, filter.Right, filter.Top)
		}

		return db.Select("*, ST_AsGeoJSON(ST_ASTEXT(ST_FORCE2D(the_geom))) as survey_100m_the_geom")
	})

	query = query.Preload("RoadConditionSurveys.RoadConditionSurvey100Ms.RoadConditionSurveyMs", func(db *gorm.DB) *gorm.DB {

		if filter.Left != "" && filter.Bottom != "" && filter.Right != "" && filter.Top != "" {
			return db.Select("*, ST_AsGeoJSON(ST_ASTEXT(ST_FORCE2D(the_geom))) AS survey_m_the_geom").Where("the_geom && ST_MakeEnvelope(?, ?, ?, ?, 4326)", filter.Left, filter.Bottom, filter.Right, filter.Top)
		}

		return db.Select("*, ST_AsGeoJSON(ST_ASTEXT(ST_FORCE2D(the_geom))) as survey_m_the_geom")
	})

	// Applying filters
	if len(roadIDs) > 0 {
		query = query.Where("road_condition.road_id IN ?", roadIDs)
	}
	// Filter by depot_code via JOIN — depot filter in Preload only filters nested data, not the main rows
	if len(depotCodes) > 0 {
		query = query.
			Joins("JOIN road AS rc_road ON road_condition.road_id = rc_road.id").
			Joins("JOIN road_section AS rc_road_section ON rc_road.road_section_id = rc_road_section.id").
			Where("rc_road_section.ref_depot_code IN ?", depotCodes)
	}
	if filter.Year != 0 {
		query = query.Where("road_condition.year = ?", filter.Year)
	}
	if filter.KmStart != 0.0 && filter.KmEnd != 0.0 {
		query = query.Where("road_condition.km_start >= ? AND road_retro_reflectivity.km_end <= ?", filter.KmStart, filter.KmEnd)
	} else if filter.KmStart != 0.0 {
		query = query.Where("road_condition.km_start >= ?", filter.KmStart)
	} else if filter.KmEnd != 0.0 {
		query = query.Where("road_condition.km_end <= ?", filter.KmEnd)
	}

	if filter.LaneNo != 0 {
		query = query.Where("road_condition.lane_no = ?", filter.LaneNo)
	}

	query = query.Where("road_condition.status = ?", "A")

	// Subquery to find the latest road condition per road ID based on created_date
	subquery := r.conn.Table("road_condition").
		Select("road_id, MAX(created_date) AS latest_date").
		Where("status = ?", "A").
		Group("road_id,lane_no")
	if filter.Year != 0 {
		subquery = subquery.Where("year = ?", filter.Year)
	}

	// Main query with subquery to get the latest records
	query = query.Joins("JOIN (?) AS latest ON road_condition.road_id = latest.road_id AND road_condition.created_date = latest.latest_date", subquery)

	// Execute the query
	err := query.Find(&roadCondition).Error
	if err != nil {
		return nil, err
	}

	return roadCondition, nil
}

func (r *repository) GetRoadRetroReflectivityDashboardMap(roadIDs []string, depotCodes []string, filter requests.ConditionMap) ([]models.RoadRetroReflectivityDashboard, error) {
	var retroReflectivity []models.RoadRetroReflectivityDashboard

	query := r.conn.Model(&models.RoadRetroReflectivityDashboard{})

	query = query.Preload("Road.RoadSection", func(db *gorm.DB) *gorm.DB {
		if len(depotCodes) > 0 {
			return db.Where("ref_depot_code IN (?)", depotCodes)
		}
		return db
	})

	query = query.Preload("RoadRetroReflectivityRanges", func(db *gorm.DB) *gorm.DB {
		if filter.Left != "" && filter.Bottom != "" && filter.Right != "" && filter.Top != "" {
			return db.Select("*, ST_AsGeoJSON(ST_ASTEXT(ST_FORCE2D(the_geom))) AS retro_range_the_geom").Where("the_geom && ST_MakeEnvelope(?, ?, ?, ?, 4326)", filter.Left, filter.Bottom, filter.Right, filter.Top)
		}
		return db.Select("*, ST_AsGeoJSON(ST_ASTEXT(ST_FORCE2D(the_geom))) as retro_range_the_geom")
	})

	query = query.Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs", func(db *gorm.DB) *gorm.DB {
		if filter.Left != "" && filter.Bottom != "" && filter.Right != "" && filter.Top != "" {
			return db.Select("*, ST_AsGeoJSON(ST_ASTEXT(ST_FORCE2D(the_geom))) AS retro_m_the_geom").Where("the_geom && ST_MakeEnvelope(?, ?, ?, ?, 4326)", filter.Left, filter.Bottom, filter.Right, filter.Top)
		}
		return db.Select("*, ST_AsGeoJSON(ST_ASTEXT(ST_FORCE2D(the_geom))) as retro_m_the_geom")
	})

	// Applying filters
	if len(roadIDs) > 0 {
		query = query.Where("road_retro_reflectivity.road_id IN ?", roadIDs)
	}
	// Filter by depot_code via JOIN — depot filter in Preload only filters nested data, not the main rows
	if len(depotCodes) > 0 {
		query = query.
			Joins("JOIN road AS rr_road ON road_retro_reflectivity.road_id = rr_road.id").
			Joins("JOIN road_section AS rr_road_section ON rr_road.road_section_id = rr_road_section.id").
			Where("rr_road_section.ref_depot_code IN ?", depotCodes)
	}
	if filter.Year != 0 {
		query = query.Where("road_retro_reflectivity.year = ?", filter.Year)
	}

	if filter.KmStart != 0.0 && filter.KmEnd != 0.0 {
		query = query.Where("road_retro_reflectivity.km_start >= ? AND road_retro_reflectivity.km_end <= ?", filter.KmStart, filter.KmEnd)
	} else if filter.KmStart != 0.0 {
		query = query.Where("road_retro_reflectivity.km_start >= ?", filter.KmStart)
	} else if filter.KmEnd != 0.0 {
		query = query.Where("road_retro_reflectivity.km_end <= ?", filter.KmEnd)
	}

	if filter.LaneNo != 0 {
		query = query.Where("road_retro_reflectivity.line_no = ?", filter.LaneNo)
	}

	query = query.Where("road_retro_reflectivity.status = ?", "A")

	// Subquery to find the latest road condition per road ID based on created_date
	subquery := r.conn.Table("road_retro_reflectivity").
		Select("road_id, MAX(created_date) AS latest_date").
		Where("status = ?", "A").
		Group("road_id,line_no")
	if filter.Year != 0 {
		subquery = subquery.Where("year = ?", filter.Year)
	}

	// Main query with subquery to get the latest records
	query = query.Joins("JOIN (?) AS latest ON road_retro_reflectivity.road_id = latest.road_id AND road_retro_reflectivity.created_date = latest.latest_date", subquery)

	// Execute the query
	err := query.Find(&retroReflectivity).Error
	if err != nil {
		return nil, err
	}

	return retroReflectivity, nil
}

func (r *repository) GetRoadConditionList(rcStatus string, roadId int) ([]models.RoadConditionList, error) {
	var rcList []models.RoadConditionList
	err := r.conn.Table("road_info AS ri").
		Select("DISTINCT ri.road_id AS road_id, rc.year AS year, rc.id AS id, rc.id_parent AS id_parent, rc.revision AS revision, ref_d.id AS direction_id, ref_d.name AS direction_name, rc.lane_no AS lane_no, rc.surveyed_date AS surveyed_date").
		Joins("LEFT JOIN ref_direction ref_d ON ri.ref_direction_id = ref_d.id").
		Joins("LEFT JOIN road_condition rc ON ri.road_id = rc.road_id").
		Where(rcStatus).
		Where("ri.road_id = ?", roadId).
		Order("year DESC, lane_no ASC, surveyed_date DESC, revision DESC").
		Find(&rcList).Error
	if err != nil {
		return nil, err
	}
	return rcList, nil

}

func (t *repository) GetAllRoadConditionByIdParent(idParent int) (models.RoadConditionAll, error) {
	var roadCondition models.RoadConditionAll
	err := t.conn.
		Preload("RoadConditionSurveys").
		Preload("RoadConditionSurveys.RoadConditionSurvey100Ms").
		Preload("RoadConditionSurveys.RoadConditionSurvey100Ms.RoadConditionSurveyMs").
		Where("id_parent = ? AND status = ?", idParent, "A").
		Find(&roadCondition).Error // Use Find to retrieve the data and store it in 'roadCondition'.

	// Check for errors in querying the database and return them if present.
	if err != nil {
		return roadCondition, err // Return the zero value of roadCondition and the error.
	}

	// If no error, return the retrieved roadCondition data and nil for the error.
	return roadCondition, nil
}

func (r *repository) GetRoadByID(roadID int) (models.RoadInfoGeomDirection, error) {
	var roadInfoGeomDirection models.RoadInfoGeomDirection
	query := r.conn
	query = query.Preload("RoadInfo")
	query = query.Preload("RoadInfo.Direction", func(db *gorm.DB) *gorm.DB {
		return db.Order("id")
	})
	query = query.Preload("RoadGeom", func(db *gorm.DB) *gorm.DB {
		return db.Order("lane_no")
	})
	if err := query.Where("id = ?", roadID).Where("is_active = ?", true).First(&roadInfoGeomDirection).Error; err != nil {
		return roadInfoGeomDirection, err
	}
	return roadInfoGeomDirection, nil
}

func (t *repository) GetUserDepartmentById(userId int) (models.UserDepartment, error) {
	var user models.UserDepartment
	if err := t.conn.Where("id = ?", userId).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *repository) GetRoadRetroReflectivityList(roadID int) ([]models.RetroReflectivityList, error) {
	var rcList []models.RetroReflectivityList
	err := r.conn.Table("road_info AS ri").
		Select("DISTINCT ri.road_id AS road_id, rrt.year AS year, rrt.id AS id, rrt.id_parent AS id_parent, rrt.revision AS revision, ref_d.id AS direction_id, ref_d.name AS direction_name, rrt.line_no AS line_no, rrt.surveyed_date AS surveyed_date,rrt.updated_date AS updated_date").
		Joins("LEFT JOIN ref_direction ref_d ON ri.ref_direction_id = ref_d.id").
		Joins("LEFT JOIN road_retro_reflectivity rrt ON ri.road_id= rrt.road_id  ").
		Where("rrt.status = 'A'").
		Where("ri.road_id = ?", roadID).
		Order("direction_id, line_no ASC, surveyed_date DESC, updated_date DESC").
		Find(&rcList).Error
	if err != nil {
		return nil, err
	}
	return rcList, nil

}

func (r *repository) GetRoadRetroReflectivityDetailsByIdParent(idParent int, refStripeTypeIDs []string) (models.RoadRetroReflectivityPreload, error) {
	var rrs models.RoadRetroReflectivityPreload

	query := r.conn.
		Preload("RoadRetroReflectivityRanges", func(db *gorm.DB) *gorm.DB {
			// Apply condition only if refStripeTypeIDs is not empty.
			if len(refStripeTypeIDs) > 0 {
				db = db.Where("ref_stripe_type_id IN (?)", refStripeTypeIDs)
			}
			return db
		}).
		Preload("RoadRetroReflectivityRanges.RefStripeColor").
		Preload("RoadRetroReflectivityRanges.RefStripeType").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs.RefStripeColor").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs.RefStripeType").
		Preload("RoadRetroReflectivityRanges.RoadRetroReflectivityMs").
		Where("id_parent = ? AND status = ?", idParent, "A")

	// Execute the query.
	err := query.Find(&rrs).Error
	if err != nil {
		// Return early if there's an error.
		return models.RoadRetroReflectivityPreload{}, err
	}

	// Return the result and nil error if the operation was successful.
	return rrs, nil
}

func (r *repository) GetDataList(model interface{}, where string) error {
	query := r.conn
	if where != "" {
		query = query.Where(where)
	}

	return query.Find(model).Error
}
