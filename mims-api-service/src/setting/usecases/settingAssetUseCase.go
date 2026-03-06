package usecases

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
)

func (su *settingUseCase) GetAssetTables(params requests.AssetTableQueryParams) (helpers.ReturnValueOfGetAssetTables, error) {
	grpId, _ := helpers.ConvertStringToInt(params.GroupID)
	total, err := su.settingRepo.CountAssetTables(params.AssetType, params.Name, grpId)
	if err != nil {
		return helpers.ReturnValueOfGetAssetTables{}, responses.NewInternalServerError()
	}

	limit, offset, page := helpers.GetlimitOffsetPage(params.Limit, params.Page, total)

	assetTables, err := su.settingRepo.GetAssetTables(limit, offset, params.AssetType, params.Name, grpId)
	if err != nil {
		return helpers.ReturnValueOfGetAssetTables{}, responses.NewInternalServerError()
	}

	assetTableStaffs, err := su.settingRepo.GetAssetTableStaffs()
	if err != nil {
		return helpers.ReturnValueOfGetAssetTables{}, responses.NewInternalServerError()
	}

	return helpers.ReturnValueOfGetAssetTables{
		AssetTables:      assetTables,
		AssetTableStaffs: assetTableStaffs,
		Limit:            limit,
		Page:             page,
		Total:            total,
	}, nil
}

func (su *settingUseCase) GetAssetTableByID(c *gin.Context, id string) (helpers.ReturnValueOfGetAssetTableById, error) {
	assetTableID, err := helpers.ConvertStringToInt(id)
	if err != nil {
		return helpers.ReturnValueOfGetAssetTableById{}, err
	}

	assetTable, err := su.settingRepo.GetAssetTableByID(assetTableID)
	if err != nil {
		if err.Error() == "record not found" {
			return helpers.ReturnValueOfGetAssetTableById{}, responses.NewNotFoundError()
		}

		return helpers.ReturnValueOfGetAssetTableById{}, responses.NewInternalServerError()
	}

	assetTableStaffs, err := su.settingRepo.GetAssetTableStaffByID(assetTableID)
	if err != nil {
		return helpers.ReturnValueOfGetAssetTableById{}, responses.NewInternalServerError()
	}

	return helpers.ReturnValueOfGetAssetTableById{
		AssetTable:       assetTable,
		AssetTableStaffs: assetTableStaffs,
	}, nil
}

func (su *settingUseCase) CreateAssetTable(request requests.AssetTableData, icon *multipart.FileHeader, c *gin.Context) error {
	geomDataType, err := GetGeomTypeColumnDataType(request.GeomType)
	if err != nil {
		return responses.NewAppErr(http.StatusBadRequest, err.Error())
	}

	assetTalbes, err := su.settingRepo.GetAllAssetTables()
	if err != nil {
		// return responses.NewInternalServerError()
		return responses.NewAppErr(http.StatusBadRequest, err.Error())
	}

	for _, assetTable := range assetTalbes {
		if request.TableName == assetTable.TableNameColumn {
			return responses.NewDuplicatedNameError("TableName:duplicate")
		}
	}

	// build query for create table
	queryString := fmt.Sprintf(`
		CREATE TABLE %s (
		id serial PRIMARY KEY,
		road_asset_id integer NOT NULL,
		surveyed_date timestamp without time zone,
		the_geom geometry NOT NULL,
		id_parent integer,
		hash_data character varying,
		is_deleted boolean DEFAULT false,
		%s,
	`, request.TableName, geomDataType)
	// the_geom_camera geometry,

	for _, column := range request.Columns {
		queryString += column.ColumnName + " " + GetDataTypeFromComponentType(column.ComponentType) + ", "
	}

	for _, column := range request.Columns {
		if column.TableNameRef != "" {
			queryString += fmt.Sprintf(`
			CONSTRAINT %s_%s_fkey FOREIGN KEY (%s)
			REFERENCES %s (id) MATCH SIMPLE
			ON UPDATE NO ACTION
			ON DELETE RESTRICT,
			`, request.TableName, column.ColumnName, column.ColumnName, column.TableNameRef)
		}
	}

	queryString += fmt.Sprintf(`CONSTRAINT %s_road_asset_id_fkey FOREIGN KEY (road_asset_id)
	REFERENCES road_asset (id) MATCH SIMPLE
	ON UPDATE NO ACTION
	ON DELETE CASCADE);`, request.TableName)

	// build query for comment on table and columns
	queryString += fmt.Sprintf("COMMENT ON TABLE %s IS '%s';", request.TableName, request.TableLabel)

	for k, v := range GetDefaultComment() {
		queryString += fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s';", request.TableName, k, v)
	}

	for _, column := range request.Columns {
		queryString += fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s';", request.TableName, column.ColumnName, column.ComponentTitle)
	}

	// create trigger
	queryString += fmt.Sprintf("CREATE TRIGGER get_asset_hash BEFORE INSERT OR UPDATE ON public.%s FOR EACH ROW EXECUTE PROCEDURE get_asset_hash();", request.TableName)
	queryString += fmt.Sprintf("CREATE TRIGGER get_road_asset_updated_date AFTER INSERT OR UPDATE ON public.%s FOR EACH ROW EXECUTE PROCEDURE get_road_asset_updated_date();", request.TableName)

	//save icon file
	var iconFilePath string
	if icon != nil {
		if !helpers.IsImageTypeJPEGOrPNG(icon.Filename) {
			return responses.NewImageTypeError()
		}

		if helpers.IsFileSizeGreaterThanLimit(icon.Size, 5) {
			return responses.NewFileSizeExceedLimitError()
		}

		dstPath, err := SaveFile(c, icon, os.Getenv("ASSET_TABLE_ICON_DIR"))
		if err != nil {
			log.Println(err)
			return responses.NewInternalServerError()
		}
		iconFilePath = dstPath
	}

	//check correctness of line color code
	if request.LineColorCode != "" {
		if err := CheckLineColorCodeCorrectness(request.LineColorCode); err != nil {
			return responses.NewAppErr(http.StatusBadRequest, err.Error())
		}
	}

	// get geom type from input
	geomType, err := GetGeomType(request.GeomType)
	if err != nil {
		return responses.NewAppErr(http.StatusBadRequest, err.Error())
	}

	// insert data to ref_asset_table
	assetTable := models.RefAssetTable{
		RefAssetID:      request.AssetGroup,
		TableNameColumn: request.TableName,
		TableLabel:      request.TableLabel,
		IconFilepath:    iconFilePath,
		LineColorCode:   request.LineColorCode,
		Seq:             0,
		IsInRoad:        GetAssetType(request.AssetType),
		IsActive:        true,
		GeomType:        geomType,
		CanDelete:       true,
	}

	//insert data to ref_asset_table_columns
	insertColumns := GetInsertColumns(request)

	if len(insertColumns) > helpers.GetLimitNumberOfColumns() {
		return responses.NewLimitNumberOfColumnsError()
	}

	//check duplicated column when create asset table
	for i := 1; i < len(insertColumns); i++ {
		if insertColumns[i].ColumnName == insertColumns[i-1].ColumnName {
			return responses.NewColumnNameDuplicateError()
		}
	}

	//insert approver and viewer
	//staffs := GetStaffs(request)

	createData := helpers.CreateAssetTable{
		QueryString:   queryString,
		AssetTable:    assetTable,
		InsertColumns: insertColumns,
		//Staffs:        staffs,
	}

	err = su.settingRepo.CreateAssetTable(createData)
	if err != nil {
		// roll back state when internal server error
		if err := os.Remove(iconFilePath); err != nil {
			log.Println(err)
		}
		log.Println(err)
		// return responses.NewInternalServerError()
		return responses.NewAppErr(http.StatusBadRequest, err.Error())
	}

	return nil
}

func CheckLineColorCodeCorrectness(lineColorCode string) error {
	re := regexp.MustCompile("#[0-9A-Fa-f]{6}")
	if len(lineColorCode) > 7 || !re.MatchString(lineColorCode) {
		return errors.New("line color code's format is invalid")
	}
	return nil
}

func GetInsertColumns(request requests.AssetTableData) []models.RefAssetTableColumns {
	result := []models.RefAssetTableColumns{
		{
			ColumnName:      "id_parent",
			ColumnSeq:       0,
			ColumnDataType:  "integer",
			ComponentType:   "hidden",
			ComponentTitle:  "Unique ID",
			IsRequired:      true,
			IsVisibleView:   false,
			IsVisibleEdit:   true,
			IsMandatory:     true,
			IsVisibleReport: false,
		},
		{
			ColumnName:      "id",
			ColumnSeq:       1,
			ColumnDataType:  "integer",
			ComponentType:   "hidden",
			ComponentTitle:  "คีย์หลัก",
			IsRequired:      true,
			IsVisibleView:   false,
			IsVisibleEdit:   true,
			IsMandatory:     true,
			IsVisibleReport: false,
		},
		{
			ColumnName:      "road_asset_id",
			ColumnSeq:       2,
			ColumnDataType:  "integer",
			ComponentType:   "hidden",
			ComponentTitle:  "รหัสข้อมูลสินทรัพย์",
			IsRequired:      true,
			IsVisibleView:   false,
			IsVisibleEdit:   true,
			IsMandatory:     true,
			IsVisibleReport: false,
		},
		{
			ColumnName:      "the_geom",
			ColumnSeq:       3,
			ColumnDataType:  "geometry",
			ComponentType:   GetComponentTypeForColumnTheGeom(request.GeomType),
			ComponentTitle:  "พิกัด",
			IsRequired:      true,
			IsVisibleView:   true,
			IsVisibleEdit:   true,
			IsMandatory:     true,
			IsVisibleReport: false,
		},
		// {
		// 	// ColumnName:     "the_geom_camera",
		// 	ColumnSeq:      4,
		// 	ColumnDataType: "geometry",
		// 	ComponentType:  "geom-point",
		// 	ComponentTitle: "พิกัดกล้อง",
		// 	IsRequired:     true,
		// 	IsVisibleView:  true,
		// 	IsVisibleEdit:  true,
		// 	IsMandatory:    true,
		// },
	}

	var nextSeq int
	switch request.GeomType {
	case "km":
		column := models.RefAssetTableColumns{
			ColumnName:      "km",
			ColumnSeq:       5,
			ColumnDataType:  "duoble precision",
			ComponentType:   "text-km",
			ComponentTitle:  "กม. ที่ตั้ง",
			IsRequired:      true,
			IsVisibleView:   true,
			IsVisibleEdit:   true,
			IsMandatory:     true,
			IsVisibleReport: true,
		}
		result = append(result, column)
		nextSeq = 6
	case "km_range":
		columns := []models.RefAssetTableColumns{
			{
				ColumnName:      "km_start",
				ColumnSeq:       5,
				ColumnDataType:  "duoble precision",
				ComponentType:   "text-km",
				ComponentTitle:  "กม. เริ่มต้น",
				IsRequired:      true,
				IsVisibleView:   true,
				IsVisibleEdit:   true,
				IsMandatory:     true,
				IsVisibleReport: true,
			},
			{
				ColumnName:      "km_end",
				ColumnSeq:       6,
				ColumnDataType:  "duoble precision",
				ComponentType:   "text-km",
				ComponentTitle:  "กม. สิ้นสุด",
				IsRequired:      true,
				IsVisibleView:   true,
				IsVisibleEdit:   true,
				IsMandatory:     true,
				IsVisibleReport: true,
			},
		}
		result = append(result, columns...)
		nextSeq = 7
	case "point":
		columns := []models.RefAssetTableColumns{
			{
				ColumnName:      "latitude",
				ColumnSeq:       5,
				ColumnDataType:  "duoble precision",
				ComponentType:   "text-number",
				ComponentTitle:  "Latitude",
				IsRequired:      true,
				IsVisibleView:   true,
				IsVisibleEdit:   true,
				IsMandatory:     true,
				IsVisibleReport: true,
			},
			{
				ColumnName:      "longitude",
				ColumnSeq:       6,
				ColumnDataType:  "duoble precision",
				ComponentType:   "text-number",
				ComponentTitle:  "Longitude",
				IsRequired:      true,
				IsVisibleView:   true,
				IsVisibleEdit:   true,
				IsMandatory:     true,
				IsVisibleReport: true,
			},
			{
				ColumnName:      "altitude",
				ColumnSeq:       7,
				ColumnDataType:  "duoble precision",
				ComponentType:   "text-number",
				ComponentTitle:  "Altitude",
				IsRequired:      true,
				IsVisibleView:   true,
				IsVisibleEdit:   true,
				IsMandatory:     true,
				IsVisibleReport: true,
			},
		}
		result = append(result, columns...)
		nextSeq = 8
	}

	columnsAdded := GetColumnsAdded(nextSeq, request)

	result = append(result, columnsAdded...)

	return result
}

func GetColumnsAdded(seq int, request requests.AssetTableData) []models.RefAssetTableColumns {
	result := []models.RefAssetTableColumns{}
	for _, reqColumn := range request.Columns {
		if reqColumn.ColumnID == 0 {

			var componentType string

			switch reqColumn.ComponentType {
			case "geometry":
				componentType = GetComponentTypeForColumnTheGeom(request.GeomType)
			default:
				componentType = reqColumn.ComponentType
			}

			column := models.RefAssetTableColumns{
				ColumnName:      reqColumn.ColumnName,
				ColumnSeq:       seq,
				ColumnDataType:  GetComponentTypeForColumnTheGeom(reqColumn.ComponentType),
				ComponentType:   componentType,
				ComponentTitle:  reqColumn.ComponentTitle,
				IsRequired:      reqColumn.IsRequired,
				IsVisibleView:   reqColumn.IsVisibleView,
				IsVisibleEdit:   reqColumn.IsVisibleEdit,
				IsMandatory:     false,
				IsVisibleReport: reqColumn.IsVisibleReport,
			}

			if reqColumn.TableNameRef != "" {
				column.TableNameRef = reqColumn.TableNameRef
			}
			result = append(result, column)
			seq++
		}
	}

	return result
}

func GetComponentTypeForColumnTheGeom(geomType string) string {
	if geomType == "km" || geomType == "point" {
		return "geom-point"
	}

	return "geom-line"
}

func GetAssetType(assetType string) bool {
	return assetType == "in"
}

func GetGeomType(geomType string) (int, error) {
	switch geomType {
	case "km":
		return 1, nil
	case "km_range":
		return 2, nil
	case "point":
		return 3, nil
	default:
		return 0, errors.New("geom type is incorrect")
	}
}

func GetDefaultComment() map[string]string {
	return map[string]string{
		"id":            "คีย์หลัก",
		"road_asset_id": "รหัสข้อมูลสินทรัพย์ (อ้างอิง road_asset.id)",
		"surveyed_date": "วันที่สำรวจ",
		"the_geom":      "พิกัด",
		// "the_geom_camera": "พิกัดกล้อง",
		"id_parent":  "Unique ID ของสินทรัพย์ (ใช้สำหรับเปรียบเทียบการเปลี่ยนแปลงข้อมูล)",
		"hash_data":  "MD5 ของข้อมูล",
		"is_deleted": "ลบสินทรัพย์ออกจากรายการเมื่ออนุมัติข้อมูลหรือไม่",
	}
}

func GetGeomTypeColumnDataType(geomType string) (string, error) {
	switch geomType {
	case "km":
		return "km double precision", nil
	case "km_range":
		return "km_start double precision,km_end double precision", nil
	case "point":
		return "latitude double precision, longitude double precision, altitude double precision", nil
	default:
		return "", errors.New("geom_type is incorrect")
	}
}

func GetDataTypeFromComponentType(componentType string) string {
	switch componentType {
	case "text-number":
		return "double precision"
	case "text-integer":
		return "integer"
	case "text-km":
		return "double precision"
	case "geom-point":
		return "geometry"
	case "geom-line":
		return "geometry"
	case "datepicker":
		return "timestamp without time zone"
	case "select":
		return "integer"
	case "hidden":
		return "integer"
	case "text":
		return "character varying"
	default:
		return "character varying"
	}
}

func (su *settingUseCase) UpdateAssetTableByID(request requests.AssetTableData, icon *multipart.FileHeader, iconFilePathStatus string, c *gin.Context) error {
	assetTableId, err := helpers.ConvertStringToInt(c.Param("id"))
	if err != nil {
		return err
	}

	assetTable, err := su.settingRepo.GetAssetTableByID(assetTableId)
	if err != nil {
		if err.Error() == "record not found" {
			return responses.NewNotFoundError()
		}
		return responses.NewInternalServerError()
	}

	// insert update values to existed asset table
	assetTable.RefAssetID = request.AssetGroup
	assetTable.TableLabel = request.TableLabel

	var dstPath string
	oldDstPath := assetTable.IconFilepath
	helpers.PrintlnJson("icon", icon)
	if icon != nil {
		if !helpers.IsImageTypeJPEGOrPNG(icon.Filename) {
			return responses.NewImageTypeError()
		}

		if helpers.IsFileSizeGreaterThanLimit(icon.Size, 5) {
			return responses.NewFileSizeExceedLimitError()
		}

		if iconFilePathStatus != "not_edit" {
			if assetTable.IconFilepath != "" {
				if err := RemoveFileIfExist(assetTable.IconFilepath); err != nil {
					logs.Error(err)
					return responses.NewInternalServerError()
				}
			}
		}

		if iconFilePathStatus != "delete" && iconFilePathStatus != "not_edit" && iconFilePathStatus != "no_file" {
			dstPath, err = SaveFile(c, icon, os.Getenv("ASSET_TABLE_ICON_DIR"))
			if err != nil {
				logs.Error(err)
				return responses.NewInternalServerError()
			}
		}

		switch iconFilePathStatus {
		case "no_file":
			assetTable.IconFilepath = ""
		case "delete":
			assetTable.IconFilepath = ""
		case "not_edit":
			assetTable.IconFilepath = oldDstPath
		case "upload":
			assetTable.IconFilepath = dstPath
		default:
			assetTable.IconFilepath = dstPath
		}
	} else {
		switch iconFilePathStatus {
		case "no_file":
			assetTable.IconFilepath = ""
		case "delete":
			assetTable.IconFilepath = ""
		case "not_edit":
			assetTable.IconFilepath = oldDstPath
		}
	}

	if request.LineColorCode != "" {
		if err := CheckLineColorCodeCorrectness(request.LineColorCode); err != nil {
			return responses.NewAppErr(http.StatusBadRequest, err.Error())
		}

		assetTable.LineColorCode = request.LineColorCode
	}

	//build query to change comment table name
	commentQuery := fmt.Sprintf("COMMENT ON TABLE %s IS '%s';\n", assetTable.TableNameColumn, request.TableLabel)

	//TODO del columns
	// get del columns from request
	delColumns := []models.RefAssetTableColumns{}
	for _, columnId := range request.DeleteColumns {
		column, err := su.settingRepo.GetColumnsByID(columnId)
		if err != nil {
			if err.Error() == "record not found" {
				return responses.NewNotFoundError()
			}

			log.Println(err)
			return responses.NewInternalServerError()
		}

		delColumns = append(delColumns, column)
	}

	// build query for rename col before delete
	var renameQuery string
	for _, column := range delColumns {
		renameQuery += fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO id_%d_%s_deleted;\n", assetTable.TableNameColumn, column.ColumnName, column.ID, column.ColumnName)
	}

	//get maxSeq in asset table columns by id
	maxSeq, err := su.settingRepo.GetColumnMaxSeqByAssetTableID(assetTableId)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}
	nextSeq := maxSeq + 1

	insertNewColums := GetColumnsAdded(nextSeq, request)

	numberOfColumns, err := su.settingRepo.CountAssetTableColumns(assetTable.TableNameColumn)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}

	if len(insertNewColums) > 0 && numberOfColumns >= helpers.GetLimitNumberOfColumns() {
		return responses.NewLimitNumberOfColumnsError()
	}

	if IsColumnNameDuplicate(insertNewColums, assetTable.AssetTableColumns) {
		return responses.NewColumnNameDuplicateError()
	}

	//build comment and add query for insert new columns
	var addNewColumnsQuery string
	for _, column := range insertNewColums {
		addNewColumnsQuery += fmt.Sprintf("ALTER TABLE %s ADD %s %s;\n", assetTable.TableNameColumn, column.ColumnName, GetDataTypeFromComponentType(column.ComponentType))
		commentQuery += fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s';\n", assetTable.TableNameColumn, column.ColumnName, column.ComponentTitle)
	}

	updateColumns := GetUpdateColumns(request.Columns, assetTable.AssetTableColumns)
	// build comment query for update columns
	for _, column := range updateColumns {
		commentQuery += fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s';\n", assetTable.TableNameColumn, column.ColumnName, column.ComponentTitle)
	}
	helpers.PrintlnJson(commentQuery)
	//staffs := GetStaffs(request)

	updateAssetTable := helpers.UpdateAssetTable{
		CommentQuery:       commentQuery,
		RenameColumnsQuery: renameQuery,
		AddColumnsQuery:    addNewColumnsQuery,
		AssetTable:         assetTable.RefAssetTable,
		InsertNewColumns:   insertNewColums,
		UpdateColumns:      updateColumns,
		//Staffs:             staffs,
		DeleteColumns: request.DeleteColumns,
	}

	if err := su.settingRepo.UpdateAssetTable(updateAssetTable); err != nil {
		//roll back state when internal server error
		if err := os.Remove(dstPath); err != nil {
			log.Println(err)
		}
		_, err := os.Create(oldDstPath)
		if err != nil {
			log.Println(err)
		}
		log.Println(err)
		return responses.NewInternalServerError()
	}

	return nil
}

func GetUpdateColumns(updateColumns []requests.Columns, existedColumns []models.RefAssetTableColumns) []models.RefAssetTableColumns {
	results := []models.RefAssetTableColumns{}
	oldColumns := map[int]models.RefAssetTableColumns{}
	for _, column := range existedColumns {
		oldColumns[column.ID] = column
	}

	for _, column := range updateColumns {
		if column.ColumnID != 0 && oldColumns[column.ColumnID].ID != 0 {
			updateColumn := oldColumns[column.ColumnID]
			updateColumn.ComponentTitle = column.ComponentTitle
			updateColumn.IsRequired = column.IsRequired
			updateColumn.IsVisibleEdit = column.IsVisibleEdit
			updateColumn.IsVisibleView = column.IsVisibleView
			updateColumn.IsVisibleReport = column.IsVisibleReport
			results = append(results, updateColumn)
		}
	}

	return results
}

func IsColumnNameDuplicate(newColumns, existedColumns []models.RefAssetTableColumns) bool {
	insertColumnName := map[string]bool{}
	for _, column := range newColumns {
		insertColumnName[column.ColumnName] = true
	}

	for _, column := range existedColumns {
		if insertColumnName[column.ColumnName] {
			return true
		}
	}

	return false
}

func (su *settingUseCase) DeleteAssetTableByID(id string) error {
	assetTableId, err := helpers.ConvertStringToInt(id)
	if err != nil {
		return err
	}

	err = su.settingRepo.DeleteAssetTable(assetTableId)
	if err != nil {
		log.Println(err)
		return responses.NewInternalServerError()
	}

	return nil
}
