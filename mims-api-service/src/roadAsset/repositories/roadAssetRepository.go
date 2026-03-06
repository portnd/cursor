package repositories

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadAsset/handlers"
	"gitlab.com/mims-api-service/src/roadAsset/usecases"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type roadAssetRepository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewRoadAssetRepositoryHandler(conn *gorm.DB) *handlers.RoadAssetHandler {
	useCase := usecases.NewRoadAssetUseCase(&roadAssetRepository{conn})
	handler := handlers.NewRoadAssetHandler(useCase)
	return handler
}

type RoadAssetDetail struct {
	IconFilepath   string `json:"icon_filepath"`
	LineColorCode  string `json:"line_color_code"`
	TableName      string `json:"table_name"`
	ColumnName     string `json:"column_name"`
	TableNameRef   string `json:"table_name_ref"`
	ColumnDataType string `json:"column_data_type"`
	ComponentTitle string `json:"component_title"`
	ComponentType  string `json:"component_type"`
	Seq            string `json:"seq"`
}

// /////////////////////////
//
//	func (t *roadAssetRepository) GetRoadAssetDetailColumn(permissions []string, roadAssetID, refAssetTableID, deptID int) (models.RoadAssetTableColumnStaff, error) {
//		var roadAssetTableColumnStaff models.RoadAssetTableColumnStaff
//		query := t.conn
//		query = query.Preload("AssetTableColumn", func(db *gorm.DB) *gorm.DB {
//			return db.Where("is_visible_view = ?", true).Order("column_seq")
//		})
//		query = query.Preload("AssetTableStaff", func(db *gorm.DB) *gorm.DB {
//			andWheres1 := []string{}
//			andWheres2 := []string{}
//			orWheres := []string{}
//			if helpers.HasPermission([]string{"road_in_asset_manage_data", "road_out_asset_manage_data"}, permissions) {
//				andWheres1 = append(andWheres1, "is_approver = true")
//				andWheres1 = append(andWheres1, fmt.Sprintf("ref_department_id = %d", deptID))
//				whereString1 := "(" + strings.Join(andWheres1, " and ") + ")"
//				orWheres = append(orWheres, whereString1)
//			}
//			if helpers.HasPermission([]string{"road_in_asset_access", "road_out_asset_access"}, permissions) {
//				andWheres2 = append(andWheres2, "is_approver = false")
//				andWheres2 = append(andWheres2, fmt.Sprintf("ref_department_id = %d", deptID))
//				whereString2 := "(" + strings.Join(andWheres2, " and ") + ")"
//				orWheres = append(orWheres, whereString2)
//			}
//			queryWhereString2 := strings.Join(orWheres, " or ")
//			return db.Where(queryWhereString2)
//		})
//		query = query.Joins("JOIN road_asset on road_asset.ref_asset_table_id = ref_asset_table.id")
//		orWheres := []string{}
//		if helpers.HasPermission([]string{"road_in_asset_manage_data", "road_in_asset_access"}, permissions) {
//			orWheres = append(orWheres, "ref_asset_table.is_in_road = true")
//		}
//		if helpers.HasPermission([]string{"road_out_asset_manage_data", "road_out_asset_access"}, permissions) {
//			orWheres = append(orWheres, "ref_asset_table.is_in_road = false")
//		}
//		queryWhereString := strings.Join(orWheres, " or ")
//		query = query.Where(queryWhereString)
//		query = query.Where("road_asset.id = ?", roadAssetID)
//		query = query.Where("ref_asset_table.id = ?", refAssetTableID)
//		if err := query.Find(&roadAssetTableColumnStaff).Error; err != nil {
//			fmt.Println(err)
//		}
//		return roadAssetTableColumnStaff, nil
//	}
func (t *roadAssetRepository) GetRoadAssetDetailColumn(permissions []string, roadAssetID, refAssetTableID int) ([]responses.RoadAssetDetailColumn, error) {
	var roadAssetDetailColumn []responses.RoadAssetDetailColumn
	query := t.conn
	sqlSelects := []string{"DISTINCT ref_at.icon_filepath as icon_filepath", "ref_at.line_color_code as line_color_code", "ref_at.table_name as table_name", "ref_atc.column_name as column_name", "ref_atc.table_name_ref as table_name_ref", "ref_atc.column_data_type as column_data_type", "ref_atc.component_title as component_title", "ref_atc.component_type as component_type", "ref_atc.column_seq as seq"}
	query = query.Table("road_asset as ra")
	query = query.Select(helpers.Implode(",", sqlSelects))
	query = query.Joins("LEFT JOIN ref_asset_table as ref_at on ra.ref_asset_table_id = ref_at.id")
	query = query.Joins("LEFT JOIN ref_asset_table_columns as ref_atc on ref_at.id = ref_atc.ref_asset_table_id")
	// query = query.Joins("LEFT JOIN ref_asset_table_staff as ref_ats on ref_ats.ref_asset_table_id = ref_at.id")
	query = query.Where("ra.id = ?", roadAssetID)
	query = query.Where("ref_at.id = ?", refAssetTableID)
	query = query.Where("ref_atc.is_visible_view = ?", true)
	orWheres := []string{}
	// if helpers.HasPermission([]string{"view_all_road_in_asset", "view_department_road_in_asset", "road_in_asset_manage_data", "road_in_asset_department_manage_data"}, permissions) {
	// 	orWheres = append(orWheres, "ref_at.is_in_road = true")
	// }

	// if helpers.HasPermission([]string{"road_out_asset_manage_data", "road_out_asset_access"}, permissions) {
	// 	orWheres = append(orWheres, "ref_at.is_in_road = false")
	// }

	queryWhereString := strings.Join(orWheres, " or ")
	query = query.Where(queryWhereString)

	// andWheres1 := []string{}
	// andWheres2 := []string{}
	orWheres2 := []string{}

	// if helpers.HasPermission([]string{"road_in_asset_manage_data", "road_in_asset_department_manage_data"}, permissions) {
	// 	if !helpers.HasPermission([]string{"road_in_asset_manage_data"}, permissions) {
	// 		andWheres1 = append(andWheres1, fmt.Sprintf("ref_ats.ref_department_id = %d", deptID))
	// 		whereString1 := "(" + strings.Join(andWheres1, " and ") + ")"
	// 		orWheres2 = append(orWheres2, whereString1)
	// 	}
	// }

	// if helpers.HasPermission([]string{"view_all_road_in_asset", "view_department_road_in_asset", "road_out_asset_access"}, permissions) {
	// 	if !helpers.HasPermission([]string{"view_all_road_in_asset"}, permissions) {
	// 		andWheres2 = append(andWheres2, "ref_ats.is_approver = false")
	// 		andWheres2 = append(andWheres2, fmt.Sprintf("ref_ats.ref_department_id = %d", deptID))
	// 		whereString2 := "(" + strings.Join(andWheres2, " and ") + ")"
	// 		orWheres2 = append(orWheres2, whereString2)
	// 	}

	// }

	queryWhereString2 := strings.Join(orWheres2, " or ")
	query = query.Where(queryWhereString2)
	query = query.Order("seq")
	if err := query.Find(&roadAssetDetailColumn).Error; err != nil {
		return roadAssetDetailColumn, err
	}
	return roadAssetDetailColumn, nil
}

func (t *roadAssetRepository) GetRoadAssetDetailData(roadAssetID int, data []responses.RoadAssetDetailColumn, permissions []string, pageNumber, pageSize int, isPagination bool) ([]map[string]interface{}, error) {
	offset := (pageNumber - 1) * pageSize
	results := []map[string]interface{}{}

	if len(data) == 0 {
		return results, nil
	}
	sqlJoin := ""
	sqlSelect := "c.id as asset_object_id, ra.id as id, ra.id_parent AS id_parent, ra.status AS status, ref_ds.name AS status_text, ra.updated_by AS updated_by, ra.updated_date AS updated_date, ra.reject_reason AS reject_reason, ra.revision AS revision, ra.is_exclusive_lock AS is_exclusive_lock, c.id_parent AS id_parent_asset"
	refTableCount := 1
	sqlOrders := []string{}
	sqlOrders = append(sqlOrders, "id ASC")
	for _, item := range data {
		if item.TableNameRef != "" {
			tableAlias := fmt.Sprintf("r%d", refTableCount)
			sqlJoin += " " + fmt.Sprintf("LEFT JOIN %s as %s on %s.id = c.%s_id", item.TableNameRef, tableAlias, tableAlias, item.TableNameRef)
			// sqlJoin += " " + fmt.Sprintf("LEFT JOIN %s as %s on %s.id = c.%s", item.TableNameRef, tableAlias, tableAlias, item.ColumnName)
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
		//
		if item.ColumnName == "km" {
			sqlOrders = append(sqlOrders, "km ASC")
		} else if item.ColumnName == "km_start" {
			sqlOrders = append(sqlOrders, "km_start ASC")
		} else if item.ColumnName == "latitude" {
			sqlOrders = append(sqlOrders, "latitude ASC")
		}
	}
	sqlWhere := ""
	// if helpers.HasPermission([]string{"road_in_asset_manage_data", "road_in_asset_department_manage_data", "road_out_asset_manage_data"}, permissions) {
	sqlWhere += " where (ra.status = 'T' or ra.status = 'W' or ra.status = 'R' or ra.status = 'A' or ra.status = 'I')"
	// } else {
	// 	sqlWhere += " where ra.status = 'A'"
	// }

	sqlWhere += fmt.Sprintf("and c.is_deleted = false and ra.id = %d", roadAssetID) // input

	sqlOrders = append(sqlOrders, "revision DESC")
	order := " order by " + helpers.Implode(",", sqlOrders)
	sqlOffset := ""
	if isPagination {
		sqlOffset = fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
	}

	sql := "select " + sqlSelect + " from road_asset as ra JOIN " + data[0].TableName + " as c on ra.id = c.road_asset_id LEFT JOIN ref_data_status as ref_ds on ref_ds.status_code = ra.status " + sqlJoin + sqlWhere + order + sqlOffset
	rows, err := t.conn.Raw(sql).Rows()
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

func (t *roadAssetRepository) GetRoadAssetDetailInfo(roadAssetID int, permissions []string) (models.RoadAssetRefDataStatus, error) {
	var roadAssetRefDataStatus models.RoadAssetRefDataStatus
	query := t.conn
	//query = query.Preload("DataStatus")
	query = query.Preload("UserDepartment")
	// query = query.Preload("UserDepartment.Department")
	// if helpers.HasPermission([]string{"road_in_asset_manage_data", "road_in_asset_department_manage_data", "road_out_asset_manage_data"}, permissions) {
	var wheres []string
	wheres = append(wheres, "status = 'T'")
	wheres = append(wheres, "status = 'W'")
	wheres = append(wheres, "status = 'R'")
	wheres = append(wheres, "status = 'A'")
	wheres = append(wheres, "status = 'I'")
	queryWhereString := strings.Join(wheres, " or ")
	query = query.Where(queryWhereString)
	// } else {
	// 	query = query.Where("status = ?", "A")
	// }
	query = query.Where("id = ?", roadAssetID)
	if err := query.Order("revision DESC").First(&roadAssetRefDataStatus).Error; err != nil {
		return roadAssetRefDataStatus, err
	}
	return roadAssetRefDataStatus, nil
}

func (t *roadAssetRepository) IsApproverAssetTableStaff(refAssetTableID, departmentID int) (bool, error) {
	var assetTableStaff models.RefAssetTableStaff
	query := t.conn
	query = query.Where("is_approver = ?", true)
	query = query.Where("ref_department_id = ?", departmentID)
	query = query.Where("ref_asset_table_id = ?", refAssetTableID)
	if err := query.First(&assetTableStaff).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (t *roadAssetRepository) GetRoadAssetRevisions(params requests.AssetRevisionsQueryParams, permissions []string, roadID int) ([]responses.RoadAssetRevision, error) {
	var RoadAssetRevision []responses.RoadAssetRevision
	refAssetTableID := params.RefAssetTableID
	query := t.conn
	sqlSelects := []string{"DISTINCT ra.id as id", "ra.updated_date as updated_date", "ra.revision as revision_no", "ra.is_exclusive_lock as is_exclusive_lock", "ra.status as status", "ref_ds.name as status_text, ra.id_parent AS id_parent "}
	sqlSelect := helpers.Implode(", ", sqlSelects)
	query = query.Table("road_asset as ra")
	query = query.Select(sqlSelect)
	query = query.Joins("LEFT JOIN ref_asset_table as ref_at on ra.ref_asset_table_id = ref_at.id")
	query = query.Joins("LEFT JOIN ref_data_status as ref_ds on ref_ds.status_code = ra.status")
	// query = query.Joins("LEFT JOIN ref_asset_table_staff as ref_ats on ref_ats.ref_asset_table_id = ref_at.id")
	query = query.Where("ra.road_id = ?", roadID)
	query = query.Where("ref_at.id = ?", refAssetTableID)
	orWheres := []string{}
	// if helpers.HasPermission([]string{"road_in_asset_manage_data", "road_in_asset_department_manage_data", "road_in_asset_access", "view_all_road_in_asset", "view_department_road_in_asset"}, permissions) {
	// 	orWheres = append(orWheres, "ref_at.is_in_road = true")
	// }
	// if helpers.HasPermission([]string{"road_out_asset_manage_data", "road_out_asset_access"}, permissions) {
	// 	orWheres = append(orWheres, "ref_at.is_in_road = false")
	// }
	queryWhereString := strings.Join(orWheres, " or ")
	query = query.Where(queryWhereString)

	orWheres2 := []string{}
	// andWheres1 := []string{}
	// isUse := false
	// if helpers.HasPermission([]string{"road_in_asset_manage_data", "road_in_asset_department_manage_data"}, permissions) {
	// 	andWheres1 = append(andWheres1, "ref_ats.is_approver = true")
	// 	andWheres1 = append(andWheres1, fmt.Sprintf("ref_ats.ref_department_id = %d", deptID))
	// 	whereString1 := "(" + strings.Join(andWheres1, " and ") + ")"
	// 	orWheres2 = append(orWheres2, whereString1)
	// 	isUse = true
	// }
	// andWheres2 := []string{}
	// isUse := false
	// if helpers.HasPermission([]string{"view_all_road_in_asset"}, permissions) {
	// 	//
	// 	isUse = true
	// } else if helpers.HasPermission([]string{"view_all_road_in_asset"}, permissions) {
	// 	if !isUse {
	// 		andWheres2 = append(andWheres2, "ref_ats.is_approver = false")
	// 		andWheres2 = append(andWheres2, fmt.Sprintf("ref_ats.ref_department_id = %d", deptID))
	// 		whereString2 := "(" + strings.Join(andWheres2, " and ") + ")"
	// 		orWheres2 = append(orWheres2, whereString2)
	// 	}
	// } else {
	// 	return RoadAssetRevision, nil
	// }

	//if helpers.HasPermission([]string{"view_all_road_in_asset", "view_department_road_in_asset"}, permissions) {
	queryWhereString2 := strings.Join(orWheres2, " or ")
	query = query.Where(queryWhereString2)

	orWheres3 := []string{}
	// if helpers.HasPermission([]string{"road_in_asset_manage_data", "road_in_asset_department_manage_data"}, permissions) {
	orWheres3 = append(orWheres3, "ra.status = 'I'")
	orWheres3 = append(orWheres3, "ra.status = 'T'")
	orWheres3 = append(orWheres3, "ra.status = 'W'")
	orWheres3 = append(orWheres3, "ra.status = 'R'")
	orWheres3 = append(orWheres3, "ra.status = 'A'")
	query = query.Where(strings.Join(orWheres3, " or "))
	// } else {
	// query = query.Where("ra.status = ?", "A")
	// }
	query = query.Order("id_parent desc")
	query = query.Order("revision_no desc")
	query = query.Order("id desc")
	if err := query.Find(&RoadAssetRevision).Error; err != nil {
		return RoadAssetRevision, err
	}
	// if err := query.Find(&RoadAssetRevision).Error; err != nil {
	// 	return "", err
	// }
	// helpers.PrintlnJson(RoadAssetRevision)
	return RoadAssetRevision, nil
}

func (t *roadAssetRepository) GetRoadAssetTemplateColumn(params requests.AssetTemplateQueryParams, permissions []string) ([]responses.RoadAssetTemplateColumn, error) {
	var roadAssetTemplateColumn []responses.RoadAssetTemplateColumn
	refAssetTableID := params.RefAssetTableID
	query := t.conn
	sqlSelects := []string{"ref_at.table_name as table_name", "ref_atc.column_name as column_name", "ref_atc.table_name_ref as table_name_ref", "ref_atc.column_data_type as column_data_type", "ref_atc.component_title as component_title", "ref_atc.component_type as component_type", "ref_atc.is_required as is_required", "ref_atc.column_seq as seq"}
	query = query.Table("ref_asset_table as ref_at")
	query = query.Select(helpers.Implode(",", sqlSelects))
	query = query.Joins("LEFT JOIN ref_asset_table_columns as ref_atc on ref_at.id = ref_atc.ref_asset_table_id")
	query = query.Where("ref_at.id = ?", refAssetTableID)
	if params.Action == "edit" {
		query = query.Where("ref_atc.is_visible_edit = ?", true)
	}

	orWheres := []string{}
	// if helpers.HasPermission([]string{"road_in_asset_manage_data"}, permissions) {
	// 	orWheres = append(orWheres, "ref_at.is_in_road = true")
	// }
	// if helpers.HasPermission([]string{"road_out_asset_manage_data"}, permissions) {
	// 	orWheres = append(orWheres, "ref_at.is_in_road = false")
	// }
	queryWhereString := strings.Join(orWheres, " or ")
	query = query.Where(queryWhereString)
	if err := query.Order("seq").Find(&roadAssetTemplateColumn).Error; err != nil {
		return roadAssetTemplateColumn, err
	}
	return roadAssetTemplateColumn, nil
}

func (t *roadAssetRepository) GetRoadAssetTemplateData(params requests.AssetTemplateQueryParams, colList []responses.RoadAssetTemplateColumn) (interface{}, error) {
	sqlSelects := []string{}
	// refTableCount := 1
	for _, col := range colList {
		// if col.TableNameRef != "" && col.TableNameRef != "null" {
		// 	tableAlias := fmt.Sprintf("r%d", refTableCount)
		// 	sqlSelects = append(sqlSelects, fmt.Sprintf("%s.name as %s_name", tableAlias, col.TableNameRef))
		// 	refTableCount++
		// }

		if strings.Contains(col.ColumnName, "geom_camera") {
			sqlSelects = append(sqlSelects, fmt.Sprintf("ST_ASTEXT(%s) AS %s", col.ColumnName, col.ColumnName))
		} else if strings.Contains(col.ColumnName, "geom") {
			sqlSelects = append(sqlSelects, fmt.Sprintf("ST_ASTEXT(ST_FORCE2D(%s)) AS %s_2d", col.ColumnName, col.ColumnName))
		} else {
			sqlSelects = append(sqlSelects, fmt.Sprintf("%s as %s", col.ColumnName, col.ColumnName))
		}
	}
	results := []map[string]interface{}{}
	sql := fmt.Sprintf("select %s from %s where id = %d", helpers.Implode(",", sqlSelects), colList[0].TableName, params.AssetObjectID)
	helpers.PrintlnJson(sql)
	rows, err := t.conn.Raw(sql).Rows()
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

func (t *roadAssetRepository) GetMaxRevision(roadID, idParent int, refAssetTableID int) (models.RaData, error) {
	orWheres := []string{}
	query := t.conn

	var roadAsset models.RaData
	if idParent != 0 { // update
		query = query.Where("id_parent = ?", idParent)
	} else { //create
		query = query.Where("road_id = ?", roadID)
		query = query.Order("id_parent DESC")
	}

	query = query.Where("ref_asset_table_id = ?", refAssetTableID)
	orWheres = append(orWheres, "status = 'A'")
	orWheres = append(orWheres, "status = 'T'")
	orWheres = append(orWheres, "status = 'W'")
	orWheres = append(orWheres, "status = 'R'")
	orWheres = append(orWheres, "status = 'I'")
	query = query.Where(strings.Join(orWheres, " or "))
	query = query.Order("revision DESC")
	if err := query.First(&roadAsset).Error; err != nil {
		return roadAsset, err
	}
	return roadAsset, nil
}

func (t *roadAssetRepository) UpdateRoadAssetStauts(RoadAssetId int, status string) error {

	if err := t.conn.Table("road_asset").Where("id = ?", RoadAssetId).Updates(models.RoadAsset{Status: status}).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadAssetRepository) UpdateRoadAssetUpdatedDate(RoadAssetId int, updatedDate time.Time) error {

	if err := t.conn.Table("road_asset").Where("id = ?", RoadAssetId).Updates(models.RoadAsset{UpdatedDate: updatedDate}).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadAssetRepository) UpdateRoadAssetIDParent(RoadAssetId int, idParent int) error {

	if err := t.conn.Table("road_asset").Where("id = ?", RoadAssetId).Updates(models.RoadAsset{IdParent: idParent}).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadAssetRepository) UpdateRoadAssetData(RoadAssetId int, status string) error {

	if err := t.conn.Table("road_asset").Where("id = ?", RoadAssetId).Updates(models.RoadAsset{Status: status}).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadAssetRepository) GetRefAssetTableById(ID int) (models.RefAssetTable, error) {

	var refAssetTable models.RefAssetTable
	if err := t.conn.Where("id = ?", ID).First(&refAssetTable).Error; err != nil {
		return refAssetTable, err
	}
	return refAssetTable, nil
}

func (t *roadAssetRepository) InsertRoadAsset(roadAsset models.RoadAsset) (int, error) {
	query := t.conn
	if err := query.Create(&roadAsset).Error; err != nil {
		return 0, err
	}
	return roadAsset.Id, nil
}

func (t *roadAssetRepository) LoadData(tableName string, roadAssetId int) ([]map[string]interface{}, error) {
	results := []map[string]interface{}{}
	sql := fmt.Sprintf("SELECT * FROM %s t WHERE t.road_asset_id = %d AND t.is_deleted = false ORDER BY id", tableName, roadAssetId)
	rows, err := t.conn.Raw(sql).Rows()
	if err != nil {
		// return results, err
	}

	cols, err := rows.Columns()
	if err != nil {
		// return results, err
	}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			// return results, err
		}

		result := make(map[string]interface{})
		for i, colName := range cols {
			// if colName == "id" {
			// 	continue
			// }
			val := columnPointers[i].(*interface{})
			result[colName] = *val

		}

		results = append(results, result)
	}
	// helpers.PrintlnJson("results", results)
	return results, nil
}

func (t *roadAssetRepository) RawQuery(sql string) error {
	query := t.conn
	if err := query.Raw(sql).Error; err != nil {
		return err
	}
	return nil
}

type aID struct {
	ID int `json:"id"`
}

func (t *roadAssetRepository) RawQueryInsert(sql string) (int, error) {
	var aID aID
	query := t.conn
	if err := query.Raw(sql).Find(&aID).Error; err != nil {
		return aID.ID, err
	}
	return aID.ID, nil
}

func (t *roadAssetRepository) FileColumnFilepath(refAssetTableID int) ([]models.RefAssetTableColumns, error) {
	var assetTableColumns []models.RefAssetTableColumns
	query := t.conn
	query = query.Where("ref_asset_table_id = ?", refAssetTableID)
	query = query.Where("column_name like ?", "%filepath%")
	if err := query.Find(&assetTableColumns).Error; err != nil {
		return assetTableColumns, err
	}
	return assetTableColumns, nil
}

func (t *roadAssetRepository) GetOldAsset(idParentAsset, maxRevisionID int, tableName string, selects []string) ([]map[string]interface{}, error) {
	results := []map[string]interface{}{}
	// for test
	selectStr := helpers.Implode(", ", selects)
	sql := fmt.Sprintf("select %s from %s where road_asset_id = %d and id_parent = %d and is_deleted = %t order by id DESC", selectStr, tableName, maxRevisionID, idParentAsset, false)
	// fmt.Println(sql)
	rows, err := t.conn.Raw(sql).Rows()
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

func (t *roadAssetRepository) UpdateTableIsDeletedByRoadAssetID(roadAssetID int, tableName string) error {
	query := t.conn
	sql := fmt.Sprintf("UPDATE %s SET is_deleted = %t WHERE road_asset_id = %d", tableName, true, roadAssetID)
	// fmt.Println(sql)
	if err := query.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadAssetRepository) UpdateTableIsDeletedByRoadAssetIDAndIDParent(assetID int, IDParent int, tableName string) error {
	query := t.conn
	sql := fmt.Sprintf("UPDATE %s SET is_deleted = %t WHERE road_asset_id = %d and id_parent = %d", tableName, true, assetID, IDParent)
	// fmt.Println(sql)
	if err := query.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadAssetRepository) UpdateTableIsDeletedByID(ID int, tableName string) error {
	query := t.conn
	sql := fmt.Sprintf("UPDATE %s SET is_deleted = %t WHERE id = %d", tableName, true, ID)
	// fmt.Println(sql)
	if err := query.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadAssetRepository) UpdateIDParentByID(ID, IDParent int, tableName string) error {
	query := t.conn
	sql := fmt.Sprintf("UPDATE %s SET id_parent = %d WHERE id = %d", tableName, IDParent, ID)
	if err := query.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadAssetRepository) UpdateFilePathByID(ID int, column, filelPath, tableName string) error {
	query := t.conn
	sql := fmt.Sprintf("UPDATE %s SET %s = '%s' WHERE id = %d", tableName, column, filelPath, ID)
	if err := query.Exec(sql).Error; err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (t *roadAssetRepository) GetRoadKmByGeomLine(geomString string, roadID int) (float64, float64, int, error) {
	var resultStart struct {
		RoadID  int     `gorm:"column:road_id"`
		LaneNO  int     `gorm:"column:lane_no"`
		KmStart float64 `gorm:"column:km"`
	}
	var resultEnd struct {
		RoadID int     `gorm:"column:road_id"`
		LaneNO int     `gorm:"column:lane_no"`
		KmEnd  float64 `gorm:"column:km"`
	}
	query := t.conn

	queryStringStart := fmt.Sprintf("SELECT * FROM get_km_from_geom_road(ST_SetSRID(ST_StartPoint(ST_GeomFromText('%s')),4326), %d)", geomString, roadID)
	err := query.Raw(queryStringStart).Scan(&resultStart).Error
	if err != nil {
		return 0, 0, 0, err
	}

	queryStringEnd := fmt.Sprintf("SELECT * FROM get_km_from_geom_road(ST_SetSRID(ST_EndPoint(ST_GeomFromText('%s')),4326), %d)", geomString, roadID)
	err = query.Raw(queryStringEnd).Scan(&resultEnd).Error
	if err != nil {
		return 0, 0, 0, err
	}

	return resultStart.KmStart, resultEnd.KmEnd, resultEnd.LaneNO, nil

}

func (t *roadAssetRepository) GetRoadKmByGeomPoint(geomString string, roadID int) (float64, int, error) {
	var result struct {
		RoadID int     `gorm:"column:road_id"`
		LaneNO int     `gorm:"column:lane_no"`
		Km     float64 `gorm:"column:km"`
	}
	// fmt.Println(geomString)
	query := t.conn

	queryStringStart := fmt.Sprintf("SELECT * FROM get_km_from_geom_road(ST_SetSRID((ST_GeomFromText('%s')),4326), %d)", geomString, roadID)
	query = query.Raw(queryStringStart)
	err := query.Scan(&result).Error
	if err != nil {
		return 0, 0, err
	}
	return result.Km, result.LaneNO, nil
}

func (t *roadAssetRepository) UpdatetTriggerHashData(ID int, tableName, hash string) error {
	query := t.conn
	sql := fmt.Sprintf("UPDATE %s SET hash_data = '%s' WHERE id = %d;", tableName, hash, ID)
	if err := query.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadAssetRepository) GetRoadAssetStatusTByIdParent(idParent int) (models.RoadAsset, error) {
	var roadAsset models.RoadAsset

	if err := t.conn.Where("id_parent = ? AND status = ? ", idParent, "T").Order("revision DESC").First(&roadAsset).Error; err != nil {
		return roadAsset, err
	}
	return roadAsset, nil
}

func (t *roadAssetRepository) GetRoadAssetExclusiveLockByIdParent(idParent int, exclusiveLock bool, status string) (models.RoadAsset, error) {
	var roadAsset models.RoadAsset

	if err := t.conn.Where("id_parent = ? AND is_exclusive_lock = ? ", idParent, exclusiveLock).Where(status).Order("revision DESC").First(&roadAsset).Error; err != nil {
		return roadAsset, err
	}
	return roadAsset, nil
}

func (t *roadAssetRepository) UpdateConfirmRoadAsset(roadAssetID int, UserID int) (models.RoadAsset, error) {
	var roadAsset models.RoadAsset
	if err := t.conn.Raw("UPDATE road_asset SET status = ?, is_exclusive_lock = ?,updated_date = ?,updated_by = ? WHERE id = ? RETURNING id", "A", false, time.Now(), UserID, roadAssetID).Scan(&roadAsset).Error; err != nil {
		return roadAsset, err
	}
	return roadAsset, nil
}

func (t *roadAssetRepository) UpdateCancelRoadAsset(roadAssetID int, UserID int) (models.RoadAsset, error) {
	var roadAsset models.RoadAsset

	if err := t.conn.Raw("UPDATE road_asset SET is_exclusive_lock = ?,status = ?,updated_date = ?,updated_by = ? WHERE id = ? RETURNING id", false, "D", time.Now(), UserID, roadAssetID).Scan(&roadAsset).Error; err != nil {
		return roadAsset, err
	}
	return roadAsset, nil
}

func (t *roadAssetRepository) GetLastRoadAssetByIdParent(parentID int, status string) (models.RoadAsset, error) {
	var roadAsset models.RoadAsset

	if err := t.conn.Where("id_parent = ? ", parentID).Where(status).Last(&roadAsset).Error; err != nil {
		return roadAsset, err
	}

	return roadAsset, nil
}

func (t *roadAssetRepository) GetRoadAssetByID(assetID int, status string) (models.RoadAsset, error) {
	var roadAsset models.RoadAsset

	if err := t.conn.Where("id = ? ", assetID).Where(status).First(&roadAsset).Error; err != nil {
		return roadAsset, err
	}

	return roadAsset, nil
}

func (t *roadAssetRepository) UpdateRoadAssetStatus(roadAssetID int, status string) (bool, error) {

	if err := t.conn.Table("road_asset").Where("id = ?", roadAssetID).Updates(models.RoadAsset{Status: status}).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (t *roadAssetRepository) CreateRoadAssetRevision(roadAsset models.RoadAsset) (int, error) {

	if err := t.conn.Save(&roadAsset).Error; err != nil {
		return 0, err
	}
	return roadAsset.Id, nil
}

func (t *roadAssetRepository) GetAllRoadAssetTableByAssetId(roadAssetID int, tableName string) ([]interface{}, error) {
	var results []interface{}

	if err := t.conn.Table(tableName).Where("is_deleted = ? ", "false").Order("id").Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (t *roadAssetRepository) DeleteRoadAssetTableByRoadAssetID(roadAssetID int, tableName string) (bool, error) {

	if err := t.conn.Table(tableName).Where("road_asset_id = ? ", roadAssetID).Update("is_deleted", true).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (t *roadAssetRepository) DeleteRoadAssetTableByID(ID int, tableName string) (bool, error) {

	if err := t.conn.Table(tableName).Where("id = ? ", ID).Update("is_deleted", true).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (t *roadAssetRepository) DeleteRoadAssetTable(idParent int, roadAssetID int, tableName string) (bool, error) {

	if err := t.conn.Table(tableName).Where("id_parent = ? AND road_asset_id = ? ", idParent, roadAssetID).Update("is_deleted", 1).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (t *roadAssetRepository) UndeleteRoadAssetTableByID(ID int, tableName string) (bool, error) {

	if err := t.conn.Table(tableName).Where("id = ?", ID).Update("is_deleted", 0).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (t *roadAssetRepository) CreateAssetTableFromOldAsset(oldAssetData []map[string]interface{}, tableName string) error {

	for _, item := range oldAssetData {
		sql := ""
		value := ""
		i := 1
		for col, val := range item {
			if col != "id" {
				// if idx == 0 {
				sql += col
				if i != len(item)-1 {
					sql += ", "
				}
				dataType := fmt.Sprintf("%v", reflect.TypeOf(val))
				switch dataType {
				case "string":
					value += fmt.Sprintf("'%v'", val.(string))
				case "int64":
					value += fmt.Sprintf("%v", val.(int64))
				case "float64":
					value += fmt.Sprintf("%v", val.(float64))
				case "time.Time":
					value += fmt.Sprintf("'%v'", val.(time.Time).Format("2006-01-02 15:04:05"))
				case "[]uint8":
					var bytes []uint8
					bytes = append(bytes, val.([]uint8)...)
					str := string(bytes)
					value += fmt.Sprintf("'%v'", str)
				case "bool":
					value += fmt.Sprintf("%v", val.(bool))
				case "int":
					value += fmt.Sprintf("%v", val.(int))
				default:
					if val == nil {
						value += fmt.Sprintf("%v", "NULL")
					}
				}

				if i != len(item)-1 {
					value += ", "
				}
				i++
			}
		}
		smt := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s)", tableName, sql, value)
		err := t.conn.Raw(smt).Error
		if err != nil {
			logs.Error(err)
			return responses.NewAppErr(400, err.Error())
		}
	}
	return nil
}

func (t *roadAssetRepository) CreateAssetTableFromOldAssetByNewRoadAssetId(roadAssetID int, oldAssetData []map[string]interface{}, tableName string) error {

	for _, item := range oldAssetData {
		sql := "road_asset_id, "
		value := fmt.Sprintf("%v, ", roadAssetID)
		i := 1
		for col, val := range item {
			if col != "id" && col != "road_asset_id" {
				// if idx == 0 {
				sql += col
				if i != len(item)-2 {
					sql += ", "
				}
				dataType := fmt.Sprintf("%v", reflect.TypeOf(val))
				switch dataType {
				case "string":
					value += fmt.Sprintf("'%v'", val.(string))
				case "int64":
					value += fmt.Sprintf("%v", val.(int64))
				case "float64":
					value += fmt.Sprintf("%v", val.(float64))
				case "time.Time":
					value += fmt.Sprintf("'%v'", val.(time.Time).Format("2006-01-02 15:04:05"))
				case "[]uint8":
					var bytes []uint8
					bytes = append(bytes, val.([]uint8)...)
					str := string(bytes)
					value += fmt.Sprintf("'%v'", str)
				case "bool":
					value += fmt.Sprintf("%v", val.(bool))
				case "int":
					value += fmt.Sprintf("%v", val.(int))
				default:
					if val == nil {
						value += fmt.Sprintf("%v", "NULL")
					}
				}

				if i != len(item)-2 {
					value += ", "
				}
				i++
			}
		}
		smt := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s)", tableName, sql, value)

		err := t.conn.Exec(smt).Error
		if err != nil {
			logs.Error(err)
			return responses.NewAppErr(400, err.Error())
		}
	}
	return nil
}

func (t *roadAssetRepository) GetStatus(statusCode string) (string, error) {
	var status models.RefDataStatus
	if err := t.conn.Where("status_code = ?", statusCode).First(&status).Error; err != nil {
		return "", err
	}
	return status.Name, nil
}

type ClosestPoint struct {
	Point string `json:"point"`
}

func (t *roadAssetRepository) GetClosestRoad(roadID int, point string) (string, error) {
	var closestPoint ClosestPoint
	sql := fmt.Sprintf("SELECT st_astext(ST_ClosestPoint(the_geom, ST_GeomFromText('%v', 4326)), 4326) as point  FROM road_info WHERE road_id = %v", point, roadID)
	if err := t.conn.Raw(sql).First(&closestPoint).Error; err != nil {
		return "", err
	}
	return closestPoint.Point, nil
}

func (t *roadAssetRepository) GetAssetTableByID(assetTableID int) (models.RefAssetTable, error) {
	var refAssetTable models.RefAssetTable
	if err := t.conn.Where("id = ?", assetTableID).First(&refAssetTable).Error; err != nil {
		return refAssetTable, err
	}
	return refAssetTable, nil
}

func (t *roadAssetRepository) GetAssetTableByIDStaff(assetTableID int) ([]models.RefAssetTableStaff, error) {
	var refAssetTableStaff []models.RefAssetTableStaff
	if err := t.conn.Where("ref_asset_table_id = ?", assetTableID).Where("is_approver = ?", true).Find(&refAssetTableStaff).Error; err != nil {
		return refAssetTableStaff, err
	}
	return refAssetTableStaff, nil
}

func (t *roadAssetRepository) GetRoadInfoByRoadID(roadID int) (models.RoadInfo, error) {
	var roadInfo models.RoadInfo
	if err := t.conn.Where("road_id = ?", roadID).First(&roadInfo).Error; err != nil {
		return roadInfo, err
	}
	return roadInfo, nil
}

func (t *roadAssetRepository) UpdateRoadAssetStatuByIdParent(idParent int, status string) (models.RoadAsset, error) {
	var roadAsset models.RoadAsset
	if err := t.conn.Table("road_asset").Where("id_parent = ? AND status = ? ", idParent, "A").Update("status", "I").Error; err != nil {
		return roadAsset, err
	}
	return roadAsset, nil
}

func (t *roadAssetRepository) GetUserDepartmentById(userID int) (models.UserDepartment, error) {
	var user models.UserDepartment
	if err := t.conn.Where("id = ?", userID).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
