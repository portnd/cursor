package usecases

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/responses"
	"gitlab.com/mims-api-service/src/roadAsset/domains"
	"gorm.io/gorm"
)

type roadAssetUseCase struct {
	roadAssetRepo domains.RoadAssetRepository
}

// init usecase
func NewRoadAssetUseCase(repo domains.RoadAssetRepository) domains.RoadAssetUseCase {
	return &roadAssetUseCase{
		roadAssetRepo: repo,
	}
}

// =========================================================
type RoadAssetColData struct {
	Geom2d struct {
		ComponentTitle string `json:"component_title"`
		ComponentType  string `json:"component_type"`
		ColumnDataType string `json:"column_data_type"`
	} `json:"geom_2d"`
	Geom3d struct {
		ComponentTitle string `json:"component_title"`
		ComponentType  string `json:"component_type"`
		ColumnDataType string `json:"column_data_type"`
	} `json:"geom_3d"`
	GeomCamera struct {
		ComponentTitle string `json:"component_title"`
		ComponentType  string `json:"component_type"`
		ColumnDataType string `json:"column_data_type"`
	} `json:"geom_camera"`
}

type RoadAssetCol struct {
	ComponentTitle string `json:"component_title"`
	ComponentType  string `json:"component_type"`
	ColumnDataType string `json:"column_data_type"`
}

type RefAsset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RefAssetImg struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Filepath          string `json:"filepath"`
	ThumbnailFilepath string `json:"thumbnail_filepath"`
}

type RefAssetSingImg struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Filepath          string `json:"filepath"`
	ThumbnailFilepath string `json:"thumbnail_filepath"`
	Abbr              string `json:"abbr"`
}
type KeyValue struct {
	Key   string       `json:"key"`
	Value RoadAssetCol `json:"value"`
}

func (t *roadAssetUseCase) GetRoadAssetDetail(params requests.AssetDetailsQueryParams, permissions []string, roadAssetID, userID int) (interface{}, int, error) {

	limit := params.Limit
	page := params.Page
	refAssetTableID := params.RefAssetTableID

	colList, err := t.roadAssetRepo.GetRoadAssetDetailColumn(permissions, roadAssetID, refAssetTableID)
	if err != nil {
		logs.Error(err)
		return "", 0, responses.NewAppErr(400, err.Error())
	}
	// return colList, 0, nil
	total, err := t.roadAssetRepo.GetRoadAssetDetailData(roadAssetID, colList, permissions, page, limit, false)
	if err != nil {
		logs.Error(err)
		return "", 0, responses.NewAppErr(400, err.Error())
	}

	result, err := t.roadAssetRepo.GetRoadAssetDetailData(roadAssetID, colList, permissions, page, limit, true)
	if err != nil {
		logs.Error(err)
		return "", 0, responses.NewAppErr(400, err.Error())
	}

	roadAssetInfo, err := t.roadAssetRepo.GetRoadAssetDetailInfo(roadAssetID, permissions)
	if err != nil {
		logs.Error(err)
		return "", 0, responses.NewAppErr(400, err.Error())
	}
	roadAssetItems := []map[string]interface{}{}
	keySort := ""
	for _, item := range result {
		roadAssetItem := make(map[string]interface{})
		if int(item["revision"].(int64)) == roadAssetInfo.Revision {
			for _, col := range colList {
				roadAssetItem["id_parent"] = item["id_parent"]
				roadAssetItem["id_parent_asset"] = item["id_parent_asset"]
				if strings.Contains(col.ColumnName, "geom") {
					idxName := strings.Replace(col.ColumnName, "the_", "", 1)
					if idxName != "geom_camera" {
						roadAssetItem[idxName+"_2d"] = item[col.ColumnName+"_2d"]
						roadAssetItem[idxName+"_3d"] = item[col.ColumnName+"_3d"]
					} else {
						roadAssetItem[idxName] = item[col.ColumnName]
					}

				} else if strings.Contains(col.ColumnName, "ref_") {
					idxName := col.ColumnName
					idxNameTrim := strings.Replace(idxName, "_id", "", 1)
					if strings.Contains(col.TableNameRef, "_image") {
						if col.TableNameRef == "ref_asset_sign_image" {
							var refAsset RefAssetSingImg
							if item[idxNameTrim+"_id"] == nil {
								refAsset.ID = 0
							} else {
								refAsset.ID = int(item[idxNameTrim+"_id"].(int64))
							}

							if item[idxNameTrim+"_name"] == nil {
								refAsset.Name = ""
							} else {
								refAsset.Name = item[idxNameTrim+"_name"].(string)
							}

							if item[idxNameTrim+"_filepath"] == nil {
								refAsset.Filepath = ""
							} else {
								refAsset.Filepath = os.Getenv("STORAGE_IP") + "/" + item[idxNameTrim+"_filepath"].(string) // fileCreateUrl
							}

							if item[idxNameTrim+"_filepath"] == nil {
								refAsset.ThumbnailFilepath = ""
							} else {
								refAsset.ThumbnailFilepath = os.Getenv("STORAGE_IP") + "/" + item[idxNameTrim+"_filepath"].(string) // getIconDirUrl
							}

							if item[idxNameTrim+"_abbr"] == nil {
								refAsset.Abbr = ""
							} else {
								refAsset.Abbr = item[idxNameTrim+"_abbr"].(string)
							}
							roadAssetItem[idxNameTrim] = refAsset
						} else {
							var refAsset RefAssetImg
							if item[idxNameTrim+"_id"] == nil {
								refAsset.ID = 0
							} else {
								refAsset.ID = int(item[idxNameTrim+"_id"].(int64))
							}

							if item[idxNameTrim+"_name"] == nil {
								refAsset.Name = ""
							} else {
								refAsset.Name = item[idxNameTrim+"_name"].(string)
							}
							if item[idxNameTrim+"_filepath"] == nil {
								refAsset.Filepath = ""
							} else {
								refAsset.Filepath = os.Getenv("STORAGE_IP") + "/" + item[idxNameTrim+"_filepath"].(string) // fileCreateUrl
							}
							if item[idxNameTrim+"_filepath"] == nil {
								refAsset.ThumbnailFilepath = ""
							} else {
								refAsset.ThumbnailFilepath = os.Getenv("STORAGE_IP") + "/" + item[idxNameTrim+"_filepath"].(string) // getIconDirUrl
							}
							roadAssetItem[idxNameTrim] = refAsset
						}

					} else {
						var refAsset RefAsset
						if item[idxNameTrim+"_id"] == nil {
							refAsset.ID = 0
						} else {
							fmt.Println("item[idxNameTrim]", item[idxNameTrim+"_id"].(int64))
							refAsset.ID = int(item[idxNameTrim+"_id"].(int64))
						}

						if item[idxNameTrim+"_name"] == nil {
							refAsset.Name = ""
						} else {
							refAsset.Name = item[idxNameTrim+"_name"].(string)
						}
						roadAssetItem[idxNameTrim] = refAsset
					}
				} else if strings.Contains(col.ColumnName, "_filepath") {
					idxName := col.ColumnName
					if item[idxName] == nil {
						roadAssetItem[idxName] = ""              // fileCreateUrl
						roadAssetItem["thumbnail_"+idxName] = "" // getIconDirUrl
					} else {
						if item[idxName].(string) == "" {
							roadAssetItem[idxName] = ""
							roadAssetItem["thumbnail_"+idxName] = ""
						} else {
							roadAssetItem[idxName] = os.Getenv("STORAGE_IP") + "/" + item[idxName].(string)              // fileCreateUrl
							roadAssetItem["thumbnail_"+idxName] = os.Getenv("STORAGE_IP") + "/" + item[idxName].(string) // getIconDirUrl
						}

					}
				} else {
					switch col.ColumnDataType {
					case "integer":
						rValue := item[col.ColumnName] //intval($r->{$c->column_name});
						roadAssetItem[col.ColumnName] = rValue
					case "double precision":
						rValue := item[col.ColumnName] //doubleval($r->{$c->column_name});
						roadAssetItem[col.ColumnName] = rValue
					default:
						rValue := item[col.ColumnName] //$r->{$c->column_name};
						roadAssetItem[col.ColumnName] = rValue
					}
				}

				roadAssetItem["asset_object_id"] = item["asset_object_id"]
			}
		}
		roadAssetItems = append(roadAssetItems, roadAssetItem)
	}

	// column
	// colData := make(map[string]interface{})
	var colData []KeyValue
	for _, col := range colList {
		var column KeyValue
		idxName := ""
		if strings.Contains(col.ColumnName, "the_") {
			idxName = strings.Replace(col.ColumnName, "the_", "", 1)
			if idxName == "geom" {
				var roadAssetCol RoadAssetCol
				roadAssetCol.ComponentTitle = col.ComponentTitle
				roadAssetCol.ColumnDataType = col.ColumnDataType
				roadAssetCol.ComponentType = col.ComponentType
				// colData["geom_2d"] = roadAssetCol
				/////
				column.Key = "geom_2d"
				column.Value = roadAssetCol
				colData = append(colData, column)

				roadAssetCol.ComponentTitle = "พิกัด 3D Point"
				// colData["geom_3d"] = roadAssetCol
				/////
				column.Key = "geom_3d"
				column.Value = roadAssetCol
			}
		} else {
			idxName = col.ColumnName
		}

		if idxName != "geom" {
			var roadAssetCol RoadAssetCol
			roadAssetCol.ComponentTitle = col.ComponentTitle
			roadAssetCol.ColumnDataType = col.ColumnDataType
			roadAssetCol.ComponentType = col.ComponentType
			// colData[idxName] = roadAssetCol

			/////
			column.Key = idxName
			column.Value = roadAssetCol
		}
		if strings.Contains(column.Value.ComponentTitle, "ลำดับ") {
			keySort = column.Key
		}
		colData = append(colData, column)
	}
	// return colData, 1, nil
	//////////////////

	status, _ := t.roadAssetRepo.GetStatus(roadAssetInfo.Status)
	if roadAssetInfo.Status == "I" {
		if roadAssetInfo.Revision == 0 {
			status = "ข้อมูลเริ่มต้น"
		} else {
			status = fmt.Sprintf("ครั้งที่ %d", roadAssetInfo.Revision)
		}
	}

	iconFilepath := ""
	thumbnailIconFilepath := ""
	if colList[0].IconFilepath != "" {
		iconFilepath = os.Getenv("STORAGE_IP") + "/" + colList[0].IconFilepath
		thumbnailIconFilepath = os.Getenv("STORAGE_IP") + "/" + colList[0].IconFilepath
	}

	var canEdit bool
	if roadAssetInfo.Status == "A" || roadAssetInfo.Status == "T" {
		canEdit = true
	}

	userInfo, err := t.roadAssetRepo.GetUserDepartmentById(roadAssetInfo.UpdatedBy)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
			logs.Error(err)
			return responses.RoadConditionDetails{}, 0, err
		}
	}

	//sort data
	sort.Slice(roadAssetItems, func(i, j int) bool {
		valI, okI := roadAssetItems[i][keySort].(float64)
		valJ, okJ := roadAssetItems[j][keySort].(float64)
		if okI && okJ {
			return valI < valJ
		}

		if !okI && !okJ {
			return false
		}
		return okI
	})

	if len(colList) > 0 {
		data := responses.RoadAssetData{
			ID:          roadAssetInfo.Id,
			IDParent:    roadAssetInfo.IdParent,
			UpdatedDate: helpers.SetTimeToString(roadAssetInfo.UpdatedDate),
			Revision:    roadAssetInfo.Revision,

			Status:       status,
			StatusCode:   roadAssetInfo.Status,
			CanEdit:      canEdit,
			RejectReason: roadAssetInfo.RejectReason,
			// Permissions:           permission,
			IsExclusiveLock:       roadAssetInfo.IsExclusiveLock,
			RoadAssets:            roadAssetItems,
			DataColumns:           colData,
			IconFilepath:          iconFilepath,          // fileCreateUrl
			ThumbnailIconFilepath: thumbnailIconFilepath, // getIconDirUrl
			LineColorCode:         colList[0].LineColorCode,
		}
		data.UpdatedBy.UID = fmt.Sprintf("%d", userInfo.Id)
		data.UpdatedBy.UserName = userInfo.Username
		// data.UpdatedBy.FullName = userInfo.Firstname + " " + userInfo.Lastname == "" ? userInfo.Firstname + " " + userInfo.Lastname : ""
		if userInfo.Username != "" {
			data.UpdatedBy.FullName = userInfo.Firstname + " " + userInfo.Lastname
		} else {
			data.UpdatedBy.FullName = "System"
		}
		// data.UpdatedBy.Department.ID = userInfo.Department.ID
		// data.UpdatedBy.Department.Name = userInfo.Department.Name
		if userInfo.ProfileImgPath != "" {
			data.UpdatedBy.ProfilePicture = os.Getenv("STORAGE_IP") + "/" + userInfo.ProfileImgPath
		}
		return data, len(total), nil
	} else {
		return responses.RoadAssetData{}, 0, nil
	}
}

func (t *roadAssetUseCase) GetRoadAssetPermission(params requests.AssetPermissionQueryParams, permissions []string) (interface{}, error) {

	return responses.Permission{}, nil
}

func (t *roadAssetUseCase) GetRoadAssetRevisions(params requests.AssetRevisionsQueryParams, permissions []string, roadID int) (interface{}, error) {
	data, err := t.roadAssetRepo.GetRoadAssetRevisions(params, permissions, roadID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	var roadAssetRevisions []responses.RoadAssetRevision
	for _, item := range data {
		var roadAssetRevision responses.RoadAssetRevision
		roadAssetRevision.ID = item.ID
		roadAssetRevision.IDParent = item.IDParent
		roadAssetRevision.RevisionNo = item.RevisionNo
		roadAssetRevision.IsExclusiveLock = item.IsExclusiveLock
		roadAssetRevision.UpdatedDate = item.UpdatedDate
		if item.Status != "I" {
			status, err := t.roadAssetRepo.GetStatus(item.Status)
			if err != nil {
				logs.Error(err)
				if err == gorm.ErrRecordNotFound {
					return responses.NoData{}, nil
				}
				return data, responses.NewAppErr(400, err.Error())
			}
			roadAssetRevision.Status = status
		} else {
			revision := item.RevisionNo
			status := ""
			if revision == 0 {
				status = "ข้อมูลเริ่มต้น"
			} else {
				status = fmt.Sprintf("ครั้งที่ %d", revision)
			}
			roadAssetRevision.Status = status
		}
		roadAssetRevisions = append(roadAssetRevisions, roadAssetRevision)
	}
	if len(roadAssetRevisions) == 0 {
		return []string{}, nil
	}
	return roadAssetRevisions, nil
}

func (t *roadAssetUseCase) GetRoadAssetTemplate(params requests.AssetTemplateQueryParams, permissions []string) (interface{}, error) {
	colList, err := t.roadAssetRepo.GetRoadAssetTemplateColumn(params, permissions)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	if params.Action == "edit" {
		template, err := t.roadAssetRepo.GetRoadAssetTemplateData(params, colList)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		templateData := template.([]map[string]interface{})
		res := []map[string]interface{}{}
		for _, item := range colList {

			data := make(map[string]interface{})
			data["table_name"] = item.TableName
			data["table_name_ref"] = item.TableNameRef
			data["column_name"] = item.ColumnName
			data["component_title"] = item.ComponentTitle
			data["component_type"] = item.ComponentType
			data["column_data_type"] = item.ColumnDataType
			data["is_required"] = item.IsRequired
			if item.ComponentType == "image" {
				STORAGE_IP := os.Getenv("STORAGE_IP") + "/"
				if templateData[0][item.ColumnName] != nil {
					data["value"] = STORAGE_IP + templateData[0][item.ColumnName].(string)
				} else {
					data["value"] = ""
				}
				//
			} else {
				data["value"] = templateData[0][item.ColumnName]
			}

			res = append(res, data)
		}
		return res, nil
	} else {
		return colList, nil
	}
}

type RoadAssetID struct {
	ID int `json:"id"`
}

func (t *roadAssetUseCase) GetRoadInfoByRoadID(roadID int) (models.RoadInfo, error) {
	roadInfo, err := t.roadAssetRepo.GetRoadInfoByRoadID(roadID)
	if err != nil {
		return roadInfo, err
	}
	return roadInfo, nil
}

func (t *roadAssetUseCase) CreateRoadAsset(reqs map[string]interface{}, roadID, IDParentAsset, userID int) (interface{}, error) {
	sameRevision := false
	newIdParent := false
	isExclusiveLock := false
	insertData := make(map[string]interface{})
	var maxRevisionData models.RaData
	var raData models.RoadAsset
	if reqs["id_parent"] != nil {
		idParent := int(reqs["id_parent"].(float64))
		refAssetTableID := int(reqs["ref_asset_table_id"].(float64))
		maxRevision, err := t.roadAssetRepo.GetMaxRevision(roadID, idParent, refAssetTableID)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		maxRevisionData = maxRevision

		if !maxRevision.IsExclusiveLock {
			isExclusiveLock = false
			raData.RoadId = roadID
			raData.CreatedDate = maxRevision.CreatedDate
			raData.CreatedBy = maxRevision.CreatedBy
			raData.UpdatedBy = maxRevision.UpdatedBy
			raData.UpdatedDate = maxRevision.UpdatedDate
			raData.IdParent = maxRevision.IDParent
			raData.Status = "T"
			raData.Revision = maxRevision.Revision + 1
			raData.IsExclusiveLock = true
			raData.RefAssetTableId = refAssetTableID
			if maxRevision.Status == "T" || maxRevision.Status == "R" {
				// t.roadAssetRepo.UpdateRoadAssetStauts(maxRevision.ID, "I")
				t.roadAssetRepo.UpdateRoadAssetStauts(maxRevision.ID, "D")
			} else if maxRevision.Status == "W" {
				return nil, responses.NewAppErr(http.StatusBadRequest, constants.DATA_WAITING_APPROVAL)
			}
		} else {
			isExclusiveLock = true
			raData.UpdatedDate = time.Now()
			sameRevision = true
		}
		// else if maxRevision.UpdatedBy != userID {
		// 	return nil, responses.NewAppErr(http.StatusBadRequest, constants.EDITED_BY_ANOTHER_USER)
		// }
	} else {
		idParent := 0
		refAssetTableID := int(reqs["ref_asset_table_id"].(float64))
		isMaxRevisionEmpty := false
		maxRevision, err := t.roadAssetRepo.GetMaxRevision(roadID, idParent, refAssetTableID)
		if err != nil {
			if err.Error() == "record not found" {
				isMaxRevisionEmpty = true
			}
		}
		fmt.Println("isMaxRevisionEmpty", isMaxRevisionEmpty)
		maxRevisionData = maxRevision
		// check maxRevision not empty
		if !isMaxRevisionEmpty {
			fmt.Println(isMaxRevisionEmpty)
			raData.RoadId = roadID
			raData.CreatedDate = maxRevision.CreatedDate
			raData.CreatedBy = maxRevision.CreatedBy
			raData.UpdatedBy = userID
			raData.UpdatedDate = maxRevision.UpdatedDate
			raData.IdParent = maxRevision.IDParent
			raData.Status = "T"
			raData.Revision = maxRevision.Revision + 1
			raData.IsExclusiveLock = true
			raData.RefAssetTableId = refAssetTableID
		} else {
			raData.RoadId = roadID
			raData.CreatedDate = time.Now()
			raData.CreatedBy = userID
			raData.UpdatedBy = userID
			raData.UpdatedDate = time.Now()
			raData.IdParent = 0
			raData.Status = "T"
			raData.Revision = 0
			raData.IsExclusiveLock = true
			raData.RefAssetTableId = refAssetTableID
			newIdParent = true
		}
	}

	for key, item := range reqs {
		if fmt.Sprintf("%v", item) == "\u003cnil\u003e" {
			continue
		}

		insertIdx := ""
		insertText := ""
		if !helpers.InArray(key, []string{"uuid", "id_parent", "id_parent_asset", "ref_asset_table_id", "road_id"}) {
			if strings.Contains(key, "geom") {
				if strings.Contains(key, "3d") {
					insertIdx = strings.TrimRight(key, "_3d")
				} else {
					insertIdx = key
				}
				insertText = fmt.Sprintf("ST_SetSRID(ST_GeomFromText('%v'), 4326)", item)
			} else {
				insertIdx = key
				insertText = fmt.Sprintf("%v", item)
			}
			insertData[insertIdx] = insertText
		}
	}
	// return insertData, nil

	//update road asset
	raID := 0
	if !sameRevision {
		// create
		roadAssetID, err := t.roadAssetRepo.InsertRoadAsset(raData)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		raID = roadAssetID
	} else {
		// update
		if isExclusiveLock {
			err := t.roadAssetRepo.UpdateRoadAssetUpdatedDate(maxRevisionData.ID, raData.UpdatedDate)
			if err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, err.Error())
			}
		}
		raID = maxRevisionData.ID
	}

	if reqs["id_parent"] == nil && newIdParent {
		err := t.roadAssetRepo.UpdateRoadAssetIDParent(raID, raID)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
	}

	//get asset table
	refAssetTableID := int(reqs["ref_asset_table_id"].(float64))
	refAssetTable, err := t.roadAssetRepo.GetRefAssetTableById(refAssetTableID)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	tableName := refAssetTable.TableNameColumn

	if IDParentAsset != 0 {
		insertData["id_parent"] = IDParentAsset
	} else {
		insertData["id_parent"] = 0
	}
	roadAssetID := 0
	if !sameRevision {
		roadAssetID = raID
	} else {
		roadAssetID = maxRevisionData.ID
	}
	insertData["road_asset_id"] = roadAssetID

	//clone old data if it's new revision
	if reqs["id_parent"] != nil && !maxRevisionData.IsExclusiveLock {
		//load data
		// oldAssetData, err := t.roadAssetRepo.LoadData(tableName, maxRevisionData.ID)
		oldAssetData, err := t.roadAssetRepo.LoadData(tableName, int(reqs["road_asset_id"].(float64)))
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}

		if len(oldAssetData) > 0 {
			for _, item := range oldAssetData {
				sql := ""
				value := ""
				i := 1
				for col, val := range item {
					if col != "id" {
						if col == "road_asset_id" {
							sql += col
							value += fmt.Sprintf("'%v'", roadAssetID)
						} else {
							sql += col
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

						}
						if i != len(item)-1 {
							sql += ", "
						}

						if i != len(item)-1 {
							value += ", "
						}
						i++
					}
				}

				smt := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s) RETURNING id", tableName, sql, value)
				_, err := t.roadAssetRepo.RawQueryInsert(smt)
				if err != nil {
					logs.Error(err)
					return "", responses.NewAppErr(400, err.Error())
				}
			}
		}
	}
	//remove previous revision asset
	if IDParentAsset != 0 {
		idParentAsset := IDParentAsset
		//get old image first
		fileCols, err := t.roadAssetRepo.FileColumnFilepath(refAssetTableID)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		selects := []string{"id"}
		for _, col := range fileCols {
			selects = append(selects, col.ColumnName)
		}
		oldAssets, err := t.roadAssetRepo.GetOldAsset(idParentAsset, maxRevisionData.ID, tableName, selects)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
		if len(oldAssets) > 0 {
			oldAsset := oldAssets[0]
			for idx, r := range oldAsset {
				// oldAssetInsert := make(map[string]interface{})
				if idx != "id" {
					if r != nil {
						// oldAssetInsert[idx] = r
						// insertData[idx] = r
						if r == nil {
							insertData[idx] = nil
						} else {
							insertData[idx] = r
						}
					}
				}
			}
			//delete old asset if it's T/R
			if maxRevisionData.Status == "T" || maxRevisionData.Status == "R" {
				assetID := int(oldAsset["id"].(int64))
				err = t.roadAssetRepo.UpdateTableIsDeletedByID(assetID, tableName)
				if err != nil {
					logs.Error(err)
					return "", responses.NewAppErr(400, err.Error())
				}
			}
		}
		fmt.Println("sameRevision", sameRevision)
		//if it's edited from confirmed data, also delete asset from cloned one
		if reqs["id_parent"] != nil && !sameRevision {
			idParentAsset := IDParentAsset
			err = t.roadAssetRepo.UpdateTableIsDeletedByRoadAssetIDAndIDParent(raID, idParentAsset, tableName)
			if err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, err.Error())
			}
		}
	}

	if reqs["id_parent"] != nil && !maxRevisionData.IsExclusiveLock {
		//delete old data if it's T,R
		if maxRevisionData.Status == "T" || maxRevisionData.Status == "R" {
			err = t.roadAssetRepo.UpdateTableIsDeletedByRoadAssetID(maxRevisionData.ID, tableName)
			if err != nil {
				logs.Error(err)
				return "", responses.NewAppErr(400, err.Error())
			}
		}
	}

	// insert to asset table
	column := ""
	value := ""
	i := 1
	for key, val := range insertData {
		dataType := fmt.Sprintf("%v", reflect.TypeOf(val))
		column += key
		if i <= len(insertData)-1 {
			column += ", "
		}
		switch dataType {
		case "string":
			if strings.Contains(val.(string), "ST_GeomFromText") {
				value += fmt.Sprintf("%v", val.(string))
			} else {
				value += fmt.Sprintf("'%v'", val.(string))
			}
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
			} else {
				value += fmt.Sprintf("%v", val)
			}
		}
		if i <= len(insertData)-1 {
			value += ", "
		}
		i++
	}
	mInsertString := fmt.Sprintf("INSERT INTO %s (%s) VALUES ", tableName, column)
	mInsertString = fmt.Sprintf("%s (%s) RETURNING id", mInsertString, value)
	aID, err := t.roadAssetRepo.RawQueryInsert(mInsertString)
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}

	if IDParentAsset == 0 {
		err = t.roadAssetRepo.UpdateIDParentByID(aID, aID, tableName)
		if err != nil {
			logs.Error(err)
			return "", responses.NewAppErr(400, err.Error())
		}
	}

	//insert images
	for key := range reqs {
		if strings.Contains(key, "_filepath") {
			imgBase64 := reqs[key].(string)
			base64 := ""
			fileType := ""
			if strings.Contains(imgBase64, "png") {
				fileType = "png"
			} else if strings.Contains(imgBase64, "jpge") {
				fileType = "jpge"
			} else {
				fileType = "jpg"
			}
			if strings.Contains(imgBase64, ",") {
				dataImgBase64 := helpers.Explode(",", imgBase64)
				base64 = dataImgBase64[1]
			} else {
				base64 = imgBase64
			}
			if base64 != "" {
				subPath := strings.TrimSuffix(key, "_filepath")
				pathOutput := fmt.Sprintf("storages/road/attachments/%d/asset/%s/", roadID, subPath)
				fileName := subPath + "_" + time.Now().Format("20060102150405")
				filePath, err := helpers.DecodeImgBase64(base64, pathOutput, fileName, fileType)
				fmt.Println(pathOutput)
				if err != nil {
					logs.Error(err)
					return "", responses.NewAppErr(400, err.Error())
				}
				err = t.roadAssetRepo.UpdateFilePathByID(aID, key, filePath, tableName)
				if err != nil {
					logs.Error(err)
					return "", responses.NewAppErr(400, err.Error())
				}
			} else {
				err = t.roadAssetRepo.UpdateFilePathByID(aID, key, "", tableName)
				if err != nil {
					logs.Error(err)
					return "", responses.NewAppErr(400, err.Error())
				}
			}

		}
	}
	//trigger hash data
	err = t.roadAssetRepo.UpdatetTriggerHashData(aID, tableName, "0")
	if err != nil {
		logs.Error(err)
		return "", responses.NewAppErr(400, err.Error())
	}
	// response
	var res RoadAssetID
	res.ID = aID
	return res, nil
}

func (t *roadAssetUseCase) GetRoadKmByGeom(geomStr string, roadID int) (responses.RoadKmByGeom, error) {
	var result responses.RoadKmByGeom
	geom, err := t.roadAssetRepo.GetClosestRoad(roadID, geomStr)
	if err != nil {
		logs.Error(err)
		return result, responses.NewAppErr(400, err.Error())
	}
	if strings.Contains(geom, "LINESTRING") {
		start, end, laneNo, err := t.roadAssetRepo.GetRoadKmByGeomLine(geomStr, roadID)
		if err != nil {
			logs.Error(err)
			return result, err
		}
		result.RoadID = roadID
		result.TheGeom = geom
		result.KmStart = start
		result.KmEnd = end
		result.Km = 0
		result.LaneNo = laneNo
		result.Type = "LINESTRING"
		return result, nil
	} else {
		point, laneNo, err := t.roadAssetRepo.GetRoadKmByGeomPoint(geomStr, roadID)
		if err != nil {
			logs.Error(err)
			return result, err
		}
		result.RoadID = roadID
		result.LaneNo = laneNo
		result.TheGeom = geom
		result.KmStart = 0
		result.KmEnd = 0
		result.Km = point
		result.Type = "POINT"
	}

	return result, nil
}

func (t *roadAssetUseCase) ConfirmRoadAsset(idParent int, UID uint) (int, error) {
	var roadAssetId int
	userID := int(UID)
	roadAsset, err := t.roadAssetRepo.GetRoadAssetStatusTByIdParent(idParent)
	if err != nil {
		logs.Error(err)
		return 0, responses.NewAppErr(500, err.Error())
	}
	helpers.PrintlnJson(roadAsset)
	switch roadAsset.IsExclusiveLock {
	case true:
		_, err := t.roadAssetRepo.UpdateRoadAssetStatuByIdParent(idParent, "A")
		if err != nil {
			logs.Error(err)
			return 0, responses.NewAppErr(500, err.Error())
		}
		rs, err := t.roadAssetRepo.UpdateConfirmRoadAsset(roadAsset.Id, userID)
		if err != nil {
			logs.Error(err)
			return 0, responses.NewAppErr(500, err.Error())
		}
		roadAssetId = rs.Id
	case false:
		return 0, responses.NewAppErr(400, constants.DATA_ALREADY_CONFIRMED)
	}

	return roadAssetId, nil
}

func (t *roadAssetUseCase) CancelRoadAsset(idParent int, UID uint) (interface{}, error) {
	userID := int(UID)
	roadAsset, err := t.roadAssetRepo.GetRoadAssetStatusTByIdParent(idParent)
	if err != nil {
		return false, responses.NewAppErr(500, err.Error())
	}
	switch roadAsset.IsExclusiveLock {
	case true:
		// if roadAsset.UpdatedBy == userID {
		_, err := t.roadAssetRepo.UpdateCancelRoadAsset(roadAsset.Id, userID)
		if err != nil {
			return false, nil
		}
		// } else {
		// 	return false, responses.NewAppErr(400, constants.DATA_NOT_MATCH_USER)
		// }

	case false:
		return false, responses.NewAppErr(400, constants.DATA_ALREADY_CONFIRMED)
	}

	return true, nil
}

func (t *roadAssetUseCase) DeleteRoadAsset(c *gin.Context, idParent int, UID uint) (interface{}, error) {

	userID := int(UID)

	orWheres := []string{}
	orWheres = append(orWheres, "status = 'A'")
	orWheres = append(orWheres, "status = 'T'")
	orWheres = append(orWheres, "status = 'R'")
	orWheres = append(orWheres, "status = 'W'")
	status := strings.Join(orWheres, " or ")

	exclusiveLock := false
	roadAsset, err := t.roadAssetRepo.GetRoadAssetExclusiveLockByIdParent(idParent, exclusiveLock, status)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			exclusiveLock = true
			roadAsset, err = t.roadAssetRepo.GetRoadAssetExclusiveLockByIdParent(idParent, exclusiveLock, status)
			if err != nil {
				logs.Error(err)
				return false, responses.NewAppErr(400, err.Error())
			}
			// if roadAsset.Id != 0 {
			// 	return false, responses.NewAppErr(400, constants.EDITED_BY_ANOTHER_USER)
			// }

		}
		// return false, responses.NewAppErr(400, err.Error())
	}

	if roadAsset.Id != 0 {

		_, err := t.roadAssetRepo.UpdateCancelRoadAsset(roadAsset.Id, userID)
		if err != nil {
			logs.Error(err)
			return false, responses.NewAppErr(400, err.Error())
		}
		//ถ้ามีการลบสถานะ A และมี I ก่อนหน้ากำหนดให้ I เป็น A
		if roadAsset.Status == "A" {

			exclusiveLock = false
			status = "status = 'I'"
			roadAsset, err = t.roadAssetRepo.GetRoadAssetExclusiveLockByIdParent(idParent, exclusiveLock, status)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return true, nil
				}
				logs.Error(err)
				return false, responses.NewAppErr(500, err.Error())
			}

			status = "A"
			_, err = t.roadAssetRepo.UpdateRoadAssetStatus(roadAsset.Id, status)
			if err != nil {
				logs.Error(err)
				return false, responses.NewAppErr(500, err.Error())
			}
		}
	}

	return true, nil
}

func (t *roadAssetUseCase) DeleteRoadAssetObject(assetID int, refAssetTableID int, assetObjectID int, UID uint) (interface{}, error) {
	userID := int(UID)

	refAssetTable, err := t.roadAssetRepo.GetRefAssetTableById(refAssetTableID)
	if err != nil {
		return false, responses.NewAppErr(500, err.Error())
	}
	tableName := strings.ToLower(refAssetTable.TableNameColumn)

	orWheres := []string{}
	orWheres = append(orWheres, "status = 'A'")
	orWheres = append(orWheres, "status = 'T'")
	orWheres = append(orWheres, "status = 'R'")
	orWheres = append(orWheres, "status = 'W'")
	orWheres = append(orWheres, "status = 'I'")
	status := strings.Join(orWheres, " or ")

	roadAsset, err := t.roadAssetRepo.GetRoadAssetByID(assetID, status)
	if err != nil {
		return false, responses.NewAppErr(500, err.Error())
	}

	orWheres = append(orWheres, "status = 'A'")
	orWheres = append(orWheres, "status = 'T'")
	orWheres = append(orWheres, "status = 'R'")
	orWheres = append(orWheres, "status = 'W'")
	orWheres = append(orWheres, "status = 'I'")
	status = strings.Join(orWheres, " or ")
	lastRoadAsset, err := t.roadAssetRepo.GetLastRoadAssetByIdParent(roadAsset.IdParent, status)
	if err != nil {
		return false, responses.NewAppErr(500, err.Error())
	}

	// roadAsset, err := t.roadAssetRepo.GetRoadAssetByID(assetID, status)
	// if err != nil {
	// 	return false, responses.NewAppErr(500, err.Error())
	// }

	// A  IsExclusiveLock = false

	if lastRoadAsset.IsExclusiveLock {
		// if lastRoadAsset.CreatedBy == userID {
		status = "T"
		isUpdated, err := t.roadAssetRepo.UpdateRoadAssetStatus(lastRoadAsset.Id, status)
		if err != nil {
			return false, responses.NewAppErr(500, err.Error())
		}
		if isUpdated {
			_, err := t.roadAssetRepo.DeleteRoadAssetTableByID(assetObjectID, tableName)
			if err != nil {
				return false, responses.NewAppErr(500, err.Error())
			}
		}

		// } else {
		// 	return false, responses.NewAppErr(400, constants.EDITED_BY_ANOTHER_USER)
		// }
	} else {

		switch lastRoadAsset.Status {
		case "W":
			return false, responses.NewAppErr(400, constants.DATA_WAITING_APPROVAL)
		case "T", "R":
			status = "D"
			_, err := t.roadAssetRepo.UpdateRoadAssetStatus(assetID, status)
			if err != nil {
				return false, responses.NewAppErr(500, err.Error())
			}
			_, err = t.roadAssetRepo.DeleteRoadAssetTableByID(assetObjectID, tableName)
			if err != nil {
				return false, responses.NewAppErr(500, err.Error())
			}
		}

		roadAssetID, err := t.roadAssetRepo.CreateRoadAssetRevision(
			models.RoadAsset{
				RoadId:          lastRoadAsset.RoadId,
				RefAssetTableId: refAssetTableID,
				UpdatedDate:     time.Now(),
				CreatedDate:     time.Now(),
				CreatedBy:       userID,
				UpdatedBy:       userID,
				IdParent:        lastRoadAsset.IdParent,
				Status:          "T",
				Revision:        lastRoadAsset.Revision + 1,
				IsExclusiveLock: true,
			})
		if err != nil {
			return false, responses.NewAppErr(500, err.Error())
		}

		_, err = t.roadAssetRepo.DeleteRoadAssetTableByID(assetObjectID, tableName)
		if err != nil {
			return false, responses.NewAppErr(500, err.Error())
		}

		oldAssetData, err := t.roadAssetRepo.LoadData(tableName, assetID)
		if err != nil {
			logs.Error(err)
			return false, responses.NewAppErr(400, err.Error())
		}

		if len(oldAssetData) > 0 {
			err := t.roadAssetRepo.CreateAssetTableFromOldAssetByNewRoadAssetId(roadAssetID, oldAssetData, tableName)
			if err != nil {
				logs.Error(err)
				return false, responses.NewAppErr(400, err.Error())
			}
		}

		switch roadAsset.Status {
		case "T", "R":
			_, err = t.roadAssetRepo.DeleteRoadAssetTableByRoadAssetID(assetID, tableName)
			if err != nil {
				return false, responses.NewAppErr(500, err.Error())
			}

			status = "D"
			_, err = t.roadAssetRepo.UpdateRoadAssetStatus(lastRoadAsset.Id, status)
			if err != nil {
				return false, responses.NewAppErr(500, err.Error())
			}

		default:

			// status = "D"
			// _, err = t.roadAssetRepo.UpdateRoadAssetStatus(lastRoadAsset.Id, status)
			// if err != nil {
			// 	return false, responses.NewAppErr(500, err.Error())
			// }

			_, err = t.roadAssetRepo.UndeleteRoadAssetTableByID(assetObjectID, tableName)
			if err != nil {
				return false, responses.NewAppErr(500, err.Error())
			}

		}

	}
	return true, nil
}

func (t *roadAssetUseCase) GetAssetTableByID(assetTableID int) (interface{}, error) {
	var assetTableType responses.AssetTableType
	assetTable, err := t.roadAssetRepo.GetAssetTableByID(assetTableID)
	if err != nil {
		return false, responses.NewAppErr(400, err.Error())
	}
	staff, _ := t.roadAssetRepo.GetAssetTableByIDStaff(assetTableID)
	departmentIDs := []int{}
	for _, item := range staff {
		departmentIDs = append(departmentIDs, item.RefDepartmentID)
	}
	switch assetTable.GeomType {
	case 1:
		assetTableType.GeomTypeID = 1
		assetTableType.GeomTypeName = "กม."
		assetTableType.Color = assetTable.LineColorCode
		assetTableType.DepartmentManage = departmentIDs
	case 2:
		assetTableType.GeomTypeID = 2
		assetTableType.GeomTypeName = "ช่วงกม."
		assetTableType.Color = assetTable.LineColorCode
		assetTableType.DepartmentManage = departmentIDs
	case 3:
		assetTableType.GeomTypeID = 3
		assetTableType.GeomTypeName = "latitude. longitude"
		assetTableType.Color = assetTable.LineColorCode
		assetTableType.DepartmentManage = departmentIDs
	}
	return assetTableType, nil
}
