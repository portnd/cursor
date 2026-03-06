package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/roadSurface/handlers"
	"gitlab.com/mims-api-service/src/roadSurface/usecases"
	"gorm.io/gorm"
)

type roadSurfaceRepositories struct {
	conn *gorm.DB
}

func NewRoadSurfaceRepositoryHandler(conn *gorm.DB) *handlers.RoadSurfaceHandler {
	usecase := usecases.NewRoadSurfaceUsecase(&roadSurfaceRepositories{conn})
	handler := handlers.NewRoadSurfaceHandler(usecase)
	return handler
}

func (r *roadSurfaceRepositories) GetRoadSurfaceData(roadID string, hasPermission bool) ([]models.RoadSurfacePreload, error) {
	var roadSurfaces []models.RoadSurfacePreload
	_, err := r.CheckDirection(roadID)
	if err != nil {
		return roadSurfaces, err
	}
	revision, err := r.latestRevision(roadID, hasPermission)
	if revision == 0 {
		return roadSurfaces, nil
	}
	if err != nil {
		return roadSurfaces, err
	}

	query := r.conn
	query = query.Preload("RoadSurfaceLane", func(db *gorm.DB) *gorm.DB {
		return db.Select("* , st_astext(the_geom) as the_geom").Order("id ASC")
	})
	query = query.Preload("RoadInfo")
	query = query.Preload("RoadSurfaceLane.RefSurface")
	query = query.Preload("RefSurfaceIdShoulderLeftData")
	query = query.Preload("RefSurfaceIdShoulderRightData")
	query = query.Preload("MaterialBase")
	query = query.Preload("MaterialSubbase")
	query = query.Preload("MaterialSubgrade")
	query = query.Where("revision = (?) AND road_id = (?) AND status != 'D'", strconv.Itoa(revision), roadID)

	err = query.Find(&roadSurfaces).Error
	if err != nil {
		return roadSurfaces, err
	}
	helpers.PrintlnJson(roadSurfaces[0].RoadSurfaceLane[0])

	return roadSurfaces, nil
}

func (r *roadSurfaceRepositories) CheckDirection(roadID string) (uint, error) {
	var direction uint
	if err := r.conn.Model(&models.Road{}).
		Select("CASE WHEN ref_direction_id = 2 THEN 2 ELSE 1 END AS direction_id").
		Where("id = (?) AND is_active = true", roadID).
		Row().
		Scan(&direction); err != nil {
		return direction, err
	}
	return direction, nil
}

func (r *roadSurfaceRepositories) latestRevision(roadID string, hasPermission bool) (int, error) {
	var latestRevision int
	var condition []string
	var where string

	condition = append(condition, "status = 'A'", "status = 'T'", "status = 'W'", "status = 'R'")
	where = strings.Join(condition, " OR ")

	direction, err := r.CheckDirection(roadID)
	if err != nil {
		return -1, err
	}
	order := " "
	if direction == 1 {
		order = ("revision DESC, km_start ASC, id DESC ")
	} else if direction == 2 {
		order = ("revision DESC, km_start DESC , id DESC ")
	}

	if err := r.conn.Model(&models.RoadSurface{}).
		Select("revision").
		Where(where).
		Where("road_id = (?)", roadID).
		Order(order).
		Limit(1).
		Row().
		Scan(&latestRevision); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No rows in result set, handle the case here
			latestRevision = 0 // Assign a default value
		} else {
			// Another error occurred, return it
			return latestRevision, err
		}
	}
	return latestRevision, nil
}

func (r *roadSurfaceRepositories) FindRefSurface(id int) (models.SurfaceShoulder, error) {
	var result models.SurfaceShoulder
	err := r.conn.Table("ref_surface").Where("id = ?", id).First(&result).Error
	if err != nil {
		return models.SurfaceShoulder{}, err
	}
	return result, nil
}

func (r *roadSurfaceRepositories) FindRefMaterial(table string, id int) (interface{}, error) {
	switch table {
	case "ref_material_base":
		var result models.RefMaterialBase
		err := r.conn.Where("id = ?", id).First(&result).Error
		if err != nil {
			return nil, err
		}
		return result, nil
	case "ref_material_subbase":
		var result models.RefMaterialSubbase
		err := r.conn.Where("id = ?", id).First(&result).Error
		if err != nil {
			return nil, err
		}
		return result, nil

	case "ref_material_subgrade":
		var result models.RefMaterialSubgrade
		err := r.conn.Where("id = ?", id).First(&result).Error
		if err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, errors.New("Input table value invalid(From FindRefMaterial Function in Repo)")
	}
}

func (r *roadSurfaceRepositories) FindRefDirectionName(roadID int) (string, error) {
	var result models.RefDirection
	id, err := r.FindRefDirectionIDByRoadID(roadID)
	if err != nil {
		return "", err
	}
	err = r.conn.Where("id = ?", id).First(&result).Error
	if err != nil {
		return "", err
	}
	return result.Name, nil
}

func (r *roadSurfaceRepositories) FindRefSurfaceNameAndSurfaceGroup(id int) (string, string, error) {
	var result models.RefSurface
	err := r.conn.Where("id = ?", id).First(&result).Error
	if err != nil {
		return "", "", err
	}
	return result.Name, result.SurfaceGroup, nil
}

func (r *roadSurfaceRepositories) FindDataStatus(code string) (string, error) {
	var result models.RefDataStatus
	err := r.conn.Where("status_code = ?", code).First(&result).Error
	if err != nil {
		return "", err
	}
	return result.Name, nil
}

func (r *roadSurfaceRepositories) FindRefDirectionIDByRoadID(id int) (int, error) {
	var result models.RoadInfo
	err := r.conn.Where("road_id = ? and status = ?", id, "A").First(&result).Error
	if err != nil {
		return 0, err
	}
	return result.RefDirectionId, nil
}

func (r *roadSurfaceRepositories) GetRole(userId uint) ([]models.UserRole, error) {
	var userRole []models.UserRole
	if err := r.conn.Where("user_id = ?", userId).Find(&userRole).Error; err != nil {
		return userRole, err
	}
	return userRole, nil
}

// func (r *roadSurfaceRepositories) GetAccessControl(roles []int) ([]models.AccessControl, error) {
// 	var accessControl []models.AccessControl
// 	query := r.conn
// 	query = query.Joins("JOIN role_access_control on access_control.id = role_access_control.access_control_id")
// 	err := query.Find(&accessControl).Error
// 	if err != nil {
// 		return accessControl,err
// 	}
// 	return accessControl, nil
// }

func (r *roadSurfaceRepositories) GetUserByID(userId uint) (models.Users, error) {
	user := models.Users{}
	// userId = 1
	err := r.conn.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return models.Users{}, err
	}
	return user, nil
}

func (r *roadSurfaceRepositories) GetUserInfoByUpdatedBy(updatedBy int) models.UserDepartment {
	user := models.UserDepartment{}

	err := r.conn.Where("id = ?", updatedBy).First(&user).Error
	if err != nil {
		return user
	}

	return user
}

func (r *roadSurfaceRepositories) GetMaxRevisionByRoadID(roadID int) (int, error) {
	var maxRevision sql.NullInt64
	row := r.conn.Model(&models.RoadSurface{}).Select("MAX(revision)").Where("road_id = ?", roadID).Row()
	err := row.Scan(&maxRevision)

	var intValue int
	if maxRevision.Valid {
		intValue = int(maxRevision.Int64)
	} else {
		intValue = 0
	}
	return intValue, err
}

func (r *roadSurfaceRepositories) GetNewIDs(roadID int) ([]int, error) {
	var result []int
	maxRevision, err := r.GetMaxRevisionByRoadID(roadID)
	if err != nil {
		return result, err
	}
	query := r.conn.Model(&models.RoadSurface{}).Select("id").Where("road_id = ? AND revision = ?", roadID, maxRevision)
	query = query.Pluck("id", &result)
	return result, err
}

func (r *roadSurfaceRepositories) GetGeomByRoadID(roadID int) ([]models.RoadGeom, error) {
	var result []models.RoadGeom
	err := r.conn.Select("the_geom,km_start,km_end,lane_no").Where("road_id = ? AND status = 'A'", roadID).Order("revision DESC,  lane_no ASC").Find(&result).Error

	return result, err
}

func (r *roadSurfaceRepositories) ClearPreviousData(roadID int, uid int, kmStart float64, kmEnd float64, numberLane int, tx *gorm.DB) error {
	var result []models.RoadSurface
	laneNo := make([]int, numberLane)
	for i := 0; i < numberLane; i++ {
		laneNo[i] = i + 1
	}
	query := r.conn
	query = query.Table("road_surface").Select("road_surface.id")
	query = query.Joins("JOIN road_surface_lane ON road_surface.id = road_surface_lane.road_surface_id")
	query = query.Where("road_surface.road_id = ?", roadID)
	direction, err := r.FindRefDirectionIDByRoadID(roadID)
	if err != nil {
		return err
	}
	if direction == 1 {
		query = query.Where("road_surface.km_start <= ? AND road_surface.km_end > ?", kmStart, kmStart)
	} else if direction == 2 {
		query = query.Where("road_surface.km_start >= ? AND road_surface.km_end < ?", kmStart, kmStart)
	}
	query = query.Where("road_surface_lane.lane_no IN (?)", laneNo)
	query = query.Where("status = 'T' OR status = 'R'")
	err = query.Find(&result).Error
	if err != nil {
		return err
	}

	err = r.UpdateStatusToDelete(result, uid, tx)
	if err != nil {
		return err
	}
	return nil
}

func (r *roadSurfaceRepositories) InsertRoadSurface(data models.RoadSurfacePointer, tx *gorm.DB) (int, error) {
	var roadSurface models.RoadSurface

	// for _,data := range datas {
	formattedCreatedDate := data.CreatedDate.Format("2006-01-02 15:04:05")
	thicknessBase := "NULL"
	if data.ThicknessBase != nil {
		thicknessBase = fmt.Sprintf("%f", *data.ThicknessBase)
	}

	refMaterialBaseId := "NULL"
	if data.RefMaterialBaseId != nil {
		refMaterialBaseId = fmt.Sprintf("%d", *data.RefMaterialBaseId)
	}

	thicknessSubbase := "NULL"
	if data.ThicknessSubbase != nil {
		thicknessSubbase = fmt.Sprintf("%f", *data.ThicknessSubbase)
	}

	refMaterialSubbaseId := "NULL"
	if data.RefMaterialSubbaseId != nil {
		refMaterialSubbaseId = fmt.Sprintf("%d", *data.RefMaterialSubbaseId)
	}

	thicknessSubgrade := "NULL"
	if data.ThicknessSubgrade != nil {
		thicknessSubgrade = fmt.Sprintf("%f", *data.ThicknessSubgrade)
	}

	refMaterialSubgradeId := "NULL"
	if data.RefMaterialSubgradeId != nil {
		refMaterialSubgradeId = fmt.Sprintf("%d", *data.RefMaterialSubgradeId)
	}

	thicknessConcreteSlab := "NULL"
	if data.ThicknessConcreteSlab != nil {
		thicknessConcreteSlab = fmt.Sprintf("%f", *data.ThicknessConcreteSlab)
	}

	thicknessSurfaceConcrete := "NULL"
	if data.ThicknessSurfaceConcrete != nil {
		thicknessSurfaceConcrete = fmt.Sprintf("%f", *data.ThicknessSurfaceConcrete)
	}

	smt := fmt.Sprintf(`INSERT INTO road_surface (road_id, "year", km_start, km_end, width_surface, width_shoulder_left, width_shoulder_right, created_by, revision, status, id_parent, reject_reason, thickness_surface, thickness_surface_concrete, ref_surface_id_shoulder_left, ref_surface_id_shoulder_right, thickness_base, ref_material_base_id, thickness_subbase, ref_material_subbase_id, thickness_subgrade, ref_material_subgrade_id, hash_data, surface_crosssection_code, ref_structure_id, area_surface,created_date,updated_by,updated_date,thickness_concrete_slab, surface_group) VALUES(%d, %d, %f, %f, %f, %f, %f,%d, %d, '%s', 0, '',  %f, %s, %d, %d, %s, %s, %s, %s, %s, %s, '', %d, %d, 0,'%s',%d,'%s',%s, %d);`,
		data.RoadId, data.Year, data.KmStart, data.KmEnd, data.WidthSurface, data.WidthShoulderLeft, data.WidthShoulderRight, data.CreatedBy, data.Revision, data.Status, data.ThicknessSurface, thicknessSurfaceConcrete, data.RefSurfaceIdShoulderLeft, data.RefSurfaceIdShoulderRight, thicknessBase, refMaterialBaseId, thicknessSubbase, refMaterialSubbaseId, thicknessSubgrade, refMaterialSubgradeId, data.SurfaceCrosssectionCode, data.SurfaceCrosssectionCode, formattedCreatedDate, data.CreatedBy, formattedCreatedDate, thicknessConcreteSlab, data.SurfaceGroup)
	if err := tx.Exec(smt).Last(&roadSurface).Error; err != nil {
		return 0, err
		// }

	}
	fmt.Println(smt)
	return roadSurface.Id, nil
}

func (r *roadSurfaceRepositories) InsertRoadSurfaceLane(data models.RoadSurfaceLane, tx *gorm.DB) error {
	// var roadSurfaceLane models.RoadSurfaceLane
	smt := fmt.Sprintf(`INSERT INTO road_surface_lane (road_surface_id, lane_no, ref_surface_id, the_geom, ref_surface_params_id, road_id) VALUES(%d, %d, %d, %s, %d, %d);`, data.RoadSurfaceId, data.LaneNo, data.RefSurfaceId, data.TheGeom, data.RefSurfaceParamsID, data.RoadId)
	if err := tx.Exec(smt).Error; err != nil {
		return err
	}
	return nil
}

func (r *roadSurfaceRepositories) UpdateIDParent(id, idParent int, tx *gorm.DB) error {
	err := tx.Model(&models.RoadSurface{}).Where("id = ?", id).Updates(map[string]interface{}{"id_parent": idParent}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *roadSurfaceRepositories) UpdateStatusToDelete(data []models.RoadSurface, uid int, tx *gorm.DB) error {
	for _, v := range data {
		v.Status = "D"
		v.UpdatedBy = uid
		v.UpdatedDate = time.Now()
		if err := tx.Model(&models.RoadSurface{}).Where("id = ?", v.Id).Updates(v).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *roadSurfaceRepositories) UpdateStatusTToDeleteAll(roadID int, tx *gorm.DB) error {
	if err := tx.Model(&models.RoadSurface{}).Where("status = ? OR status = ? AND road_id = ?", "T", "R", roadID).Update("Status", "D").Error; err != nil {
		return err
	}

	return nil
}

func (r *roadSurfaceRepositories) GetTotalKm(roadID int) (float64, error) {
	type TotalKm struct {
		Start float64 `gorm:"column:km_start"`
		End   float64 `gorm:"column:km_end"`
	}

	result := TotalKm{}
	err := r.conn.Table("road_info").Select("km_start,km_end").Where("road_id = ?", roadID).Scan(&result).Error
	if err != nil {
		return 0, err
	}
	totalKm := math.Abs(result.End - result.Start)
	return totalKm, nil
}

func (r *roadSurfaceRepositories) StartTransSection() *gorm.DB {
	tx := r.conn.Begin()
	return tx
}

func (r *roadSurfaceRepositories) RollBack(tx *gorm.DB) error {
	tx.Rollback()
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func (r *roadSurfaceRepositories) Commit(tx *gorm.DB) error {
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *roadSurfaceRepositories) GetRoadSurfaceGroupByRoadID(roadID int) (models.RoadSurface, error) {
	var roadSurface models.RoadSurface
	err := r.conn.Select("max(surface_group) as surface_group").Where("road_id = ?", roadID).Find(&roadSurface).Error
	if err != nil {
		return roadSurface, err
	}
	return roadSurface, nil
}

func (r *roadSurfaceRepositories) GetRefSurfaceParam() ([]models.RefSurfaceParam, error) {
	var refSurfaceParam []models.RefSurfaceParam
	err := r.conn.Where("is_latest = ?", true).Find(&refSurfaceParam).Error
	if err != nil {
		return refSurfaceParam, err
	}
	return refSurfaceParam, nil
}

func (r *roadSurfaceRepositories) GetDataList(model interface{}, where string) error {
	query := r.conn
	if where != "" {
		query = query.Where(where)
	}
	err := query.Find(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *roadSurfaceRepositories) GetLastRoadInfoByID(roadId int) (*models.RoadInfo, error) {
	var roadInfo models.RoadInfo
	query := t.conn
	query = query.Where("status = 'A' and road_id = ?", roadId)
	err := query.Find(&roadInfo).Error
	if err != nil {
		return nil, err
	}
	return &roadInfo, nil
}

func (r *roadSurfaceRepositories) ClearPreviousDataStatus(roadID int, uid int, tx *gorm.DB) error {
	var result []models.RoadSurface

	query := r.conn
	query = query.Table("road_surface").Select("road_surface.id")
	query = query.Joins("JOIN road_surface_lane ON road_surface.id = road_surface_lane.road_surface_id")
	query = query.Where("road_surface.road_id = ?", roadID)
	// query = query.Where("road_surface.km_start <= ? AND road_surface.km_end > ?", kmStart, kmStart)
	// query = query.Where("road_surface_lane.lane_no IN (?)", laneNumbers)
	// query = query.Where("road_surface_lane.direction_id IN (?)", directions)
	query = query.Where("status = 'A'")
	err := query.Find(&result).Error
	if err != nil {
		return err
	}

	uniqueResults := make([]models.RoadSurface, 0)
	seenIDs := make(map[int]bool)
	for _, v := range result {
		if _, found := seenIDs[v.Id]; !found {
			seenIDs[v.Id] = true
			uniqueResults = append(uniqueResults, v)
		}
	}
	for _, v := range uniqueResults {
		dataUpdate := models.RoadSurface{
			Status:      "I",
			UpdatedBy:   uid,
			UpdatedDate: time.Now(),
		}
		// Update the record in the database
		if err := r.conn.Model(&models.RoadSurface{}).Where("id = ?", v.Id).Updates(dataUpdate).Error; err != nil {
			// Log the error or handle it as needed
			return err
		}

	}

	return nil
}

func (r *roadSurfaceRepositories) GetRoadSurfaceIconById(roadId int) ([]models.RoadSurfaceIcon, error) {
	var result []models.RoadSurfaceIcon
	query := r.conn
	query = query.
		Select("DISTINCT ON (ref_surface.id) ref_surface.id, ref_surface.name, ref_surface.color as color_code").
		Joins("RIGHT JOIN road_surface_lane ON road_surface.id = road_surface_lane.road_surface_id").
		Joins("LEFT JOIN ref_surface ON ref_surface.id = road_surface_lane.ref_surface_id").
		Where("road_surface.road_id = ? AND road_surface.status = ?", roadId, "A").
		Order("ref_surface.id")
	err := query.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
