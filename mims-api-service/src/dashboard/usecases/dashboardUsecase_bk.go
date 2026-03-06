package usecases

// import (
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"

// 	"gitlab.com/mims-api-service/helpers"
// 	"gitlab.com/mims-api-service/models"
// 	"gitlab.com/mims-api-service/requests"
// 	"gitlab.com/mims-api-service/responses"
// 	"gitlab.com/mims-api-service/src/dashboard/domains"
// )

// type useCase struct {
// 	repo domains.Repository
// }

// func NewUsecase(repo domains.Repository) domains.UseCase {
// 	return &useCase{repo: repo}
// }

// func (u *useCase) GetAssetMap(roadID []int, assetId []int, filter requests.AssetMap) (interface{}, error) {
// 	var resps []responses.AssetMap

// 	listTable, err := u.repo.GetTableResult(roadID, assetId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	zoom, err := strconv.Atoi(filter.Zoom)
// 	if err != nil {
// 		return nil, err
// 	}

// 	tables := []string{}
// 	for _, t := range listTable {
// 		tables = append(tables, t.TableName)
// 	}

// 	if zoom <= 0 {
// 		return "", responses.NewAppErr(400, "Zoom level must be greater than 0")
// 	}

// 	if zoom >= 12 {

// 		unionQuery := buildUnionQuery(tables)

// 		fmt.Println("buildUnionQuery:", unionQuery)

// 		theGeomCusters, err := u.repo.GetRoadAssetTheGeomCuster(unionQuery, roadID, filter)
// 		if err != nil {
// 			return nil, err
// 		}
// 		fmt.Println("theGeomCuster", theGeomCusters)
// 		for _, theGeoms := range theGeomCusters {
// 			var resp responses.AssetMap

// 			if len(theGeoms.TheGeomCluster) == 0 {
// 				continue
// 			}

// 			theGeomCuster, err := helpers.ConvertThegeomToGeomJSON(theGeoms.TheGeomCluster)
// 			if err != nil {
// 				return nil, err
// 			}

// 			fmt.Println("TableName:", t.TableName)
// 			resp.Name = t.TableLabel
// 			resp.IconFilepath = t.IconFilepath
// 			resp.ThumbnailIconFilepath = ""
// 			resp.LineColorCode = t.LineColorCode
// 			resp.IsCluster = true
// 			resp.ID = theGeoms.AssetID
// 			resp.Cluster = theGeoms.TotalTheGeomCluster
// 			resp.TheGeom = theGeomCuster

// 			resps = append(resps, resp)

// 		}

// 		return resps, nil

// 	} else {

// 		roadSigns, err := u.repo.GetRoadAssetSignDataByRoadIDs(t, roadID, filter)
// 		if err != nil {
// 			return nil, err
// 		}

// 		fmt.Println("roadSigns", roadSigns)

// 		createAssetLocation := CreateAssetLocation(t, roadSigns)
// 		fmt.Println("createAssetLocation", createAssetLocation)
// 		for _, v := range createAssetLocation {
// 			var resp responses.AssetMap

// 			resp.ID = v.ID
// 			resp.Name = v.Name
// 			resp.IconFilepath = v.IconFilepath
// 			resp.ThumbnailIconFilepath = v.ThumbnailIconFilepath
// 			resp.LineColorCode = v.LineColorCode

// 			geomJson, err := helpers.ConvertThegeomToGeomJSON(v.Wkt)
// 			if err != nil {
// 				return nil, err
// 			}
// 			resp.TheGeom = geomJson

// 			resps = append(resps, resp)

// 		}
// 	}

// 	return resps, nil

// }

// func buildUnionQuery(tables []string) string {
// 	var parts []string

// 	for _, table := range tables {
// 		parts = append(parts, generateSQLForTable(table))
// 	}
// 	return strings.Join(parts, " UNION ALL ")
// }

// func generateSQLForTable(tableName string) string {
// 	// Assume tableName is validated or fetched from a controlled source to prevent SQL Injection
// 	return fmt.Sprintf(`
//         SELECT t.the_geom FROM %s
//         JOIN road_asset ra ON ra.id = t.road_asset_id
//         WHERE ra.status = 'A' AND t.is_deleted = false AND t.the_geom && ST_MakeEnvelope('100.408635', '13.643633137698634', '100.7165502011776', '13.731250697431998', 4326)
//     `, tableName)
// }

// func CreateAssetLocation(t models.TableResult, roadSign []models.RoadAssetSign) []models.AssetLocation {
// 	var results []models.AssetLocation
// 	STORAGE_IP := os.Getenv("STORAGE_IP") + "/"
// 	for _, g := range roadSign {
// 		var assLo models.AssetLocation
// 		ThumbnailIconFilepath := t.IconFilepath
// 		if g.SignImageFilePath != "" {
// 			ThumbnailIconFilepath = g.SignImageFilePath
// 		} else if g.ImgFilePath != "" {
// 			ThumbnailIconFilepath = g.ImgFilePath
// 		}
// 		//type_var := strings.Split(g.GeomCl, "(")[0]
// 		assLo.ID = t.AssetId
// 		assLo.Name = t.TableLabel
// 		//assLo.Type = type_var

// 		if t.IconFilepath != "" {
// 			assLo.IconFilepath = STORAGE_IP + t.IconFilepath
// 		}

// 		if ThumbnailIconFilepath != "" {
// 			assLo.ThumbnailIconFilepath = STORAGE_IP + ThumbnailIconFilepath
// 		}

// 		assLo.LineColorCode = t.LineColorCode

// 		assLo.Wkt = g.TheGeom
// 		results = append(results, assLo)
// 	}
// 	return results
// }
