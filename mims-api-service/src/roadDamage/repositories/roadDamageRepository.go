package repositories

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/helpers"
	models "gitlab.com/mims-api-service/models"
	requests "gitlab.com/mims-api-service/requests"
	"gitlab.com/mims-api-service/src/roadDamage/handlers"
	"gitlab.com/mims-api-service/src/roadDamage/usecases"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type roadDamageRepository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewRoadDamageRepositoryHandler(conn *gorm.DB) *handlers.RoadDamageHandler {
	useCase := usecases.NewRoadDamageUseCase(&roadDamageRepository{conn})
	handler := handlers.NewRoadDamageHandler(useCase)
	return handler
}

func (t *roadDamageRepository) GetRoadDamageByIDParent(parentId int) (models.RoadDamage, error) {
	var roadDamage models.RoadDamage
	if err := t.conn.Where("id_parent = ?", parentId).Where("status != ?", "D").Order("revision DESC").First(&roadDamage).Error; err != nil {
		return roadDamage, err
	}
	return roadDamage, nil
}
func (t *roadDamageRepository) GetRoadDamageList(roadId int, permissions []string) ([]models.RoadDamageList, error) {
	query := t.conn
	var roadDamageList []models.RoadDamageList
	query = query.Table("road")
	query = query.Select(`road_damage.*,
						ref_direction.id as direction_id,
						ref_direction.name as direction_name`)
	query = query.Joins("Left join road_damage ON road.id = road_damage.road_id")
	query = query.Joins("Left join road_info ON road.id = road_info.road_id")
	query = query.Joins("Left join ref_direction ON road.ref_direction_id = ref_direction.id")
	query = query.Where("road.id = ?", roadId)
	//-------- Permission Check --------
	var wheres []string
	wheres = append(wheres, "road_damage.status = 'A'")
	queryWhereString := strings.Join(wheres, " or ")
	query = query.Where(queryWhereString)

	query = query.Order("status DESC")
	query = query.Order("road_damage.year DESC")
	query = query.Order("lane_no")
	query = query.Order("surveyed_date DESC")
	query = query.Order("revision DESC")
	err := query.Find(&roadDamageList).Error
	if err != nil {
		return roadDamageList, err
	}
	return roadDamageList, nil
}

func (t *roadDamageRepository) GetRoadDamageForImport(roadID, IDParent int) (models.RoadDamage, error) {
	var wheres []string
	wheres = append(wheres, "status = 'A'")
	queryWhereString := strings.Join(wheres, " or ")

	var roadDamage models.RoadDamage
	query := t.conn
	query = query.Where(queryWhereString)
	if err := query.Where("id_parent = ?", IDParent).Where("road_id = ?", roadID).Order("revision DESC").First(&roadDamage).Error; err != nil {
		return roadDamage, err
	}
	return roadDamage, nil
}

func (t *roadDamageRepository) GetRoadById(roadId int) (models.Road, error) {
	var road models.Road
	if err := t.conn.Where("id = ?", roadId).Where("is_active = ?", true).First(&road).Error; err != nil {
		return road, err
	}
	return road, nil
}

func (t *roadDamageRepository) GetRoadGeom(roadId, laneNo int) (models.RoadGeom, error) {
	var roadGeom models.RoadGeom
	if err := t.conn.Where("road_id = ?", roadId).Where("lane_no = ?", laneNo).Where("status = ?", "A").Order("revision DESC").First(&roadGeom).Error; err != nil {
		return roadGeom, err
	}
	return roadGeom, nil
}

func (t *roadDamageRepository) UpdateRoadDamageIDParent(roadDamageId int) error {
	var roadDamage models.RoadDamage
	if err := t.conn.Model(roadDamage).Where("id = ?", roadDamageId).Updates(models.RoadDamage{IdParent: roadDamageId}).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadDamageRepository) CreateRoadDamage(item requests.SumValues, rcData requests.RcData) (models.RoadDamage, error) {
	DIGIT, _ := strconv.Atoi(os.Getenv("DIGIT"))
	var roadDamage models.RoadDamage
	roadDamage.RoadId = rcData.RoadId
	roadDamage.LaneNo = rcData.LaneNo
	roadDamage.Year = rcData.Year
	roadDamage.KmStart = rcData.KmStart
	roadDamage.KmEnd = rcData.KmEnd
	roadDamage.SurveyedDate = rcData.SurveyedDate
	roadDamage.Revision = rcData.Revision
	roadDamage.IdParent = rcData.IDParent
	roadDamage.Status = rcData.Status
	roadDamage.SurveyType = item.SurveyType
	roadDamage.AcIcrack = helpers.RoundFloat(item.AcIcrack, DIGIT)
	roadDamage.AcUcrack = helpers.RoundFloat(item.AcUcrack, DIGIT)
	roadDamage.AcRavelling = helpers.RoundFloat(item.AcRavelling, DIGIT)
	roadDamage.AcPatching = helpers.RoundFloat(item.AcPatching, DIGIT)
	roadDamage.AcPotholeArea = helpers.RoundFloat(item.AcPotholeArea, DIGIT)
	roadDamage.AcBleeding = helpers.RoundFloat(item.AcBleeding, DIGIT)
	roadDamage.AcPotholeCount = helpers.RoundFloat(item.AcPotholeCount, DIGIT)
	roadDamage.CcTransverseCrack = helpers.RoundFloat(item.CcTransverseCrack, DIGIT)
	roadDamage.CcNonTransverseCrack = helpers.RoundFloat(item.CcNonTransverseCrack, DIGIT)
	roadDamage.CcCornerBreak = helpers.RoundFloat(item.CcCornerBreak, DIGIT)
	roadDamage.CcJointSealDamage = helpers.RoundFloat(item.CcJointSealDamage, DIGIT)
	roadDamage.CcPatching = helpers.RoundFloat(item.CcPatching, DIGIT)
	roadDamage.CcSpalling = helpers.RoundFloat(item.CcSpalling, DIGIT)
	roadDamage.CcScaling = helpers.RoundFloat(item.CcScaling, DIGIT)
	roadDamage.CreatedDate = time.Now()
	roadDamage.CreatedBy = rcData.CreatedBy
	roadDamage.UpdatedDate = time.Now()
	roadDamage.UpdatedBy = rcData.UpdatedBy

	if err := t.conn.Create(&roadDamage).Error; err != nil {
		return roadDamage, err
	}
	return roadDamage, nil

}

func (t *roadDamageRepository) CreateRoadDamageRange(item requests.RaodRangeItem, rcData requests.RcData, roadDamageId int) (int, string, error) {
	DIGIT, _ := strconv.Atoi(os.Getenv("DIGIT"))
	var roadDamageRange models.RoadDamageRange
	acIcrack := helpers.RoundFloat(item.AcIcrack, DIGIT)
	acUcrack := helpers.RoundFloat(item.AcUcrack, DIGIT)
	acRavelling := helpers.RoundFloat(item.AcRavelling, DIGIT)
	acPatching := helpers.RoundFloat(item.AcPatching, DIGIT)
	acPotholeArea := helpers.RoundFloat(item.AcPotholeArea, DIGIT)
	acBleeding := helpers.RoundFloat(item.AcBleeding, DIGIT)
	acPotholeCount := helpers.RoundFloat(item.AcPotholeCount, DIGIT)
	ccTransverseCrack := helpers.RoundFloat(item.CcTransverseCrack, DIGIT)
	ccNonTransverseCrack := helpers.RoundFloat(item.CcNonTransverseCrack, DIGIT)
	ccCornerBreak := helpers.RoundFloat(item.CcCornerBreak, DIGIT)
	ccJointSealDamage := helpers.RoundFloat(item.CcJointSealDamage, DIGIT)
	ccPatching := helpers.RoundFloat(item.CcPatching, DIGIT)
	ccSpalling := helpers.RoundFloat(item.CcSpalling, DIGIT)
	ccScaling := helpers.RoundFloat(item.CcScaling, DIGIT)
	// helpers.PrintlnJson("item.TheGeom", item.TheGeom)
	smt := fmt.Sprintf(`INSERT INTO road_damage_range 
	(road_damage_id, km_start, km_end, the_geom, survey_type, ac_icrack, ac_ucrack, ac_ravelling, ac_patching, 
	ac_pothole_area, ac_bleeding, ac_pothole_count, cc_transverse_crack, cc_non_transverse_crack, cc_corner_break, 
	cc_joint_seal_damage, cc_patching, cc_spalling, cc_scaling) 
	VALUES (%d, %f, %f, %s, '%s', %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f)`,
		roadDamageId, item.KMStart, item.KMEnd, item.TheGeom, item.SurveyType, acIcrack, acUcrack, acRavelling,
		acPatching, acPotholeArea, acBleeding, acPotholeCount, ccTransverseCrack, ccNonTransverseCrack, ccCornerBreak,
		ccJointSealDamage, ccPatching, ccSpalling, ccScaling)
	// helpers.PrintlnJson(item)
	if err := t.conn.Exec(smt).Last(&roadDamageRange).Error; err != nil {
		return 0, "", err
	}
	return roadDamageRange.Id, item.SurveyType, nil
}

func (t *roadDamageRepository) CreateRoadDamageM(item requests.RoadDamageMItem, rcData requests.RcData) error {
	var roadDamageRange models.RoadDamageRange
	if err := t.conn.Where("id = ?", item.RoadDamageRangeID).First(&roadDamageRange).Error; err != nil {
		return err
	}
	surveyType := roadDamageRange.SurveyType
	if surveyType == item.SurveyType {
		DIGIT, _ := strconv.Atoi(os.Getenv("DIGIT"))
		acIcrack := helpers.RoundFloat(item.AcIcrack, DIGIT)
		acUcrack := helpers.RoundFloat(item.AcUcrack, DIGIT)
		acRavelling := helpers.RoundFloat(item.AcRavelling, DIGIT)
		acPatching := helpers.RoundFloat(item.AcPatching, DIGIT)
		acPotholeArea := helpers.RoundFloat(item.AcPotholeArea, DIGIT)
		acBleeding := helpers.RoundFloat(item.AcBleeding, DIGIT)
		acPotholeCount := helpers.RoundFloat(item.AcPotholeCount, DIGIT)
		ccTransverseCrack := helpers.RoundFloat(item.CcTransverseCrack, DIGIT)
		ccNonTransverseCrack := helpers.RoundFloat(item.CcNonTransverseCrack, DIGIT)
		ccCornerBreak := helpers.RoundFloat(item.CcCornerBreak, DIGIT)
		ccJointSealDamage := helpers.RoundFloat(item.CcJointSealDamage, DIGIT)
		ccPatching := helpers.RoundFloat(item.CcPatching, DIGIT)
		ccSpalling := helpers.RoundFloat(item.CcSpalling, DIGIT)
		ccScaling := helpers.RoundFloat(item.CcScaling, DIGIT)
		smt := fmt.Sprintf(`INSERT INTO road_damage_m 
	(road_damage_range_id, km, the_geom, img_filepath, survey_type, ac_icrack, ac_ucrack, ac_ravelling, ac_patching, 
	ac_pothole_area, ac_bleeding, ac_pothole_count, cc_transverse_crack, cc_non_transverse_crack, cc_corner_break, 
	cc_joint_seal_damage, cc_patching, cc_spalling, cc_scaling) 
	VALUES (%d, %f, %s, '%s', '%s', %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f)`,
			item.RoadDamageRangeID, item.KM, item.TheGeomPoint, item.ImgFilepath, item.SurveyType, acIcrack, acUcrack,
			acRavelling, acPatching, acPotholeArea, acBleeding, acPotholeCount, ccTransverseCrack, ccNonTransverseCrack,
			ccCornerBreak, ccJointSealDamage, ccPatching, ccSpalling, ccScaling)
		if err := t.conn.Exec(smt).Error; err != nil {
			return err
		}
	}

	return nil
}

func (t *roadDamageRepository) UpdateRoadDamage(ID, userId int, req requests.RoadDamageImport) error {
	query := t.conn

	// new road damage
	var roadDamage models.RoadDamage
	if err := query.Where("id = ?", ID).First(&roadDamage).Error; err != nil {
		return err
	}
	surveyedDate, err := time.Parse("2006-01-02 15:04:05", req.SurveyedDate)
	if err != nil {
		return err
	}
	// roadDamage.Id = 0
	// roadDamage.Status = "A"
	roadDamage.SurveyedDate = surveyedDate
	// roadDamage.CreatedBy = userId
	roadDamage.UpdatedBy = userId
	// roadDamage.CreatedDate = time.Now()
	roadDamage.UpdatedDate = time.Now()
	// if err := query.Create(&roadDamage).Error; err != nil {
	// 	return err
	// }
	if err := query.Where("id = ?", ID).Updates(&roadDamage).Error; err != nil {
		return err
	}

	// // new road damage range
	// var roadDamageRange []models.RoadDamageRange
	// if err := query.Where("road_damage_id = ?", ID).Find(&roadDamageRange).Error; err != nil {
	// 	return err
	// }
	// for _, damageRange := range roadDamageRange {
	// 	oldDamageRange := damageRange.Id
	// 	damageRange.Id = 0
	// 	damageRange.RoadDamageId = roadDamage.Id
	// 	if err := query.Create(&damageRange).Error; err != nil {
	// 		return err
	// 	}

	// 	///
	// 	var roadDamageM []models.RoadDamageM
	// 	if err := query.Where("road_damage_range_id = ?", oldDamageRange).Find(&roadDamageM).Error; err != nil {
	// 		return err
	// 	}
	// 	for _, damageM := range roadDamageM {

	// 		damageM.Id = 0
	// 		damageM.RoadDamageRangeId = damageRange.Id
	// 		helpers.PrintlnJson(damageM)
	// 		if err := query.Create(&damageM).Error; err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}

// var roadDamage models.RoadDamage
func (t *roadDamageRepository) UpdateRoadDamageStatusT(ID int, req requests.RoadDamageImport) error {
	surveyedDate, err := time.Parse("2006-01-02 15:04:05", req.SurveyedDate)
	if err != nil {
		return err
	}

	// var roadDamage models.RoadDamage
	// if err := t.conn.LogMode(true).Where("id = ?", ID).First(&roadDamage).Error; err != nil {
	// 	return err
	// }
	// roadDamage.Id = 0
	// roadDamage.SurveyedDate = surveyedDate
	// roadDamage.Year = surveyedDate.Year()
	// roadDamage.LaneNo = req.LaneNo
	// roadDamage.Status = "T"
	// if err := t.conn.Create(&roadDamage).Error; err != nil {
	// 	return err
	// }

	// var roadDamageRange models.RoadDamage
	// if err := t.conn.LogMode(true).Where("road_damage_id = ?", ID).First(&roadDamage).Error; err != nil {
	// 	return err
	// }
	if err := t.conn.Table("road_damage").Where("id = ?", ID).Updates(models.RoadDamage{SurveyedDate: surveyedDate, Year: surveyedDate.Year(), LaneNo: req.LaneNo, Status: "T"}).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadDamageRepository) UpdateRoadDamageStatusI(ID int, req requests.RoadDamageImport) error {
	if err := t.conn.Table("road_damage").Where("id = ?", ID).Updates(models.RoadDamage{Status: "I"}).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadDamageRepository) UpdateRoadDamageCsvFilepath(ID int, csvFilepath string) error {
	var roadDamage models.RoadDamage
	if err := t.conn.Model(roadDamage).Where("id = ?", ID).Updates(models.RoadDamage{DamageInputFilepath: csvFilepath}).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadDamageRepository) UpdateRoadDamageImgFilepath(ID int, imgFilepath string) error {
	// var roadDamage models.RoadDamage
	// helpers.PrintlnJson("imgFilepath", imgFilepath)
	query := t.conn
	sql := fmt.Sprintf("UPDATE road_damage SET img_filepath = '%s' WHERE id = %d ", imgFilepath, ID)
	if err := query.Exec(sql).Error; err != nil {
		query.Rollback()
		return err
	}
	// if err := query.Model(roadDamage).Where("id = ?", ID).Updates(models.RoadDamage{ImgFilepath: imgFilepath}).Error; err != nil {
	// 	return err
	// }
	return nil
}

func (t *roadDamageRepository) GetRoadDamageById(ID int) (models.RoadDamage, error) {
	var roadDamage models.RoadDamage
	if err := t.conn.Where("id = ?", ID).First(&roadDamage).Error; err != nil {
		return roadDamage, err
	}
	return roadDamage, nil
}

func (t *roadDamageRepository) GetDirectionById(roadID int) (models.RefDirection, error) {
	var direction models.RefDirection
	var roadInfo models.RoadInfo
	if err := t.conn.Where("road_id = ?", roadID).First(&roadInfo).Error; err != nil {
		return direction, err
	}
	directionID := roadInfo.RefDirectionId
	if err := t.conn.Where("id = ?", directionID).First(&direction).Error; err != nil {
		return direction, err
	}
	return direction, nil
}

func (t *roadDamageRepository) GetRoadDamageDirection(roadId, idParent int, permissions []string) (int, error) {
	direction := struct {
		ID          int `json:"id"`
		DirectionId int `json:"direction_id"`
	}{}
	query := t.conn
	query = query.Table("road_damage")
	query = query.Select("road_info.ref_direction_id as direction_id")
	query = query.Joins("lEFT JOIN road on road.id = road_damage.road_id")
	query = query.Joins("lEFT JOIN road_info on road_info.road_id = road.id")
	query = query.Where("road_damage.id_parent =? ", idParent)
	// if helpers.HasPermission([]string{"road_damage_manage_data"}, permissions) {
	var wheres []string
	wheres = append(wheres, "road_damage.status = 'I'")
	wheres = append(wheres, "road_damage.status = 'A'")
	queryWhereString := strings.Join(wheres, " or ")
	query = query.Where(queryWhereString)
	// } else {
	// query = query.Where("road_damage.status = ?", "A")
	// }
	query = query.Where("road_damage.road_id = ?", roadId)
	query = query.Where("road.is_active = ?", true)
	if err := query.Order("road_damage.revision DESC").First(&direction).Error; err != nil {
		return 0, err
	}

	return direction.DirectionId, nil
}

func (t *roadDamageRepository) GetRoadDamageDetail(roadId, idParent, directionId int, permissions []string) (models.RoadDamageDetail, error) {
	var road models.RoadDamageDetail
	query := t.conn
	query = query.Preload("RoadDamageStatus")
	query = query.Preload("RoadDamageRange", func(db *gorm.DB) *gorm.DB {
		order := "km_start DESC"
		if directionId == 1 {
			order = "km_start"
		}

		return db.Select("id, road_damage_id, km_start, km_end, ST_ASTEXT(ST_FORCE2D(the_geom)) as the_geom, survey_type, ac_icrack, ac_ucrack, ac_ravelling, ac_patching, ac_pothole_area, ac_bleeding, ac_pothole_count, cc_transverse_crack, cc_non_transverse_crack, cc_corner_break, cc_joint_seal_damage, cc_patching, cc_spalling, cc_scaling").Order(order)
	})
	query = query.Preload("RoadDamageRange.RoadDamageM", func(db *gorm.DB) *gorm.DB {
		order := "km DESC"
		if directionId == 1 {
			order = "km"
		}
		return db.Select("id, road_damage_range_id, km, ST_ASTEXT(ST_FORCE2D(the_geom)) as the_geom, img_filepath, hash_data, survey_type, ac_icrack, ac_ucrack, ac_ravelling, ac_patching, ac_pothole_area, ac_bleeding, ac_pothole_count, cc_transverse_crack, cc_non_transverse_crack, cc_corner_break, cc_joint_seal_damage, cc_patching, cc_spalling, cc_scaling").Order(order)
	})
	// if helpers.HasPermission([]string{"road_damage_manage_data"}, permissions) {
	var wheres []string
	wheres = append(wheres, "road_damage.status = 'I'")
	wheres = append(wheres, "road_damage.status = 'A'")
	queryWhereString := strings.Join(wheres, " or ")
	query = query.Where(queryWhereString)
	// } else {
	// query = query.Where("road_damage.status = ?", "A")
	// }
	query = query.Select("id, road_id, lane_no, year, km_start, km_end, surveyed_date, created_by, created_date, updated_by, updated_date, revision, status, id_parent, CONCAT('" + os.Getenv("STORAGE_IP") + "/" + "',img_filepath) as img_filepath, CONCAT('" + os.Getenv("STORAGE_IP") + "/" + "',damage_input_filepath) as damage_input_filepath, reject_reason, survey_type, ac_icrack, ac_ucrack, ac_ravelling, ac_patching, ac_pothole_area, ac_bleeding, ac_pothole_count, cc_transverse_crack, cc_non_transverse_crack, cc_corner_break, cc_joint_seal_damage, cc_patching, cc_spalling, cc_scaling")
	query = query.Where("id_parent = ?", idParent).Where("road_id = ?", roadId)
	query = query.Order("revision DESC")
	query = query.Order("road_id")
	query = query.Order("updated_date DESC")
	if err := query.Find(&road).Error; err != nil {
		// fmt.Println(err)
		return road, err

	}
	return road, nil
}

func (t *roadDamageRepository) GetUserById(UserID int) (models.UserDepartment, error) {
	user := models.UserDepartment{}
	err := t.conn.Where("id = ?", UserID).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (t *roadDamageRepository) GetRoadDamageTemplate(roadId int) (models.RoadInfoGeom, error) {
	var roadInfoGeom models.RoadInfoGeom
	query := t.conn
	query = query.Preload("RoadInfo")
	query = query.Preload("RoadGeom", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", "A").Order("revision DESC")
	})

	query = query.Preload("RoadSection", func(db *gorm.DB) *gorm.DB {
		db = db.Select("road_section.id,road_section.road_group_id,road_section.number,road_section.name_origin_th,road_section.name_destination_th,road_section.name_origin_en,road_section.name_destination_en,road_section.km_start, road_section.km_end, road_section.distance,road_section.ref_division_code,road_section.ref_district_code,road_section.ref_depot_code, ARRAY_AGG(ref_province.name) AS province").
			Joins("JOIN ref_province ON ref_province.province_code = ANY(road_section.province_code)").Group("road_section.id,road_section.road_group_id,road_section.number,road_section.name_origin_th,road_section.name_destination_th,road_section.name_origin_en,road_section.name_destination_en,road_section.km_start, road_section.km_end, road_section.distance,road_section.ref_division_code,road_section.ref_district_code,road_section.ref_depot_code")

		return db.Order("road_section.id ASC")
	})
	query = query.Preload("RoadSection.RoadGroup")
	if err := query.Where("id = ?", roadId).Where("is_active = ?", true).Order("road_code").First(&roadInfoGeom).Error; err != nil {
		return roadInfoGeom, err
	}
	return roadInfoGeom, nil
}

func (t *roadDamageRepository) DeleteRoadDamageForImport(id, userID int) error {
	var roadDamage models.RoadDamage
	var wheres []string
	wheres = append(wheres, "status = 'A'")
	wheres = append(wheres, "status = 'T'")
	wheres = append(wheres, "status = 'R'")
	wheres = append(wheres, "status = 'W'")
	queryWhereString := strings.Join(wheres, " or ")
	if err := t.conn.Where("id_parent = ?", id).Where(queryWhereString).Order("revision DESC").First(&roadDamage).Error; err != nil {
		return err
	}

	if roadDamage.Status == "W" {
		return errors.New(constants.DATA_WAITING_APPROVAL)
	}

	if roadDamage.Status == "A" {
		var lastRoadDamageStatusI models.RoadDamage
		if err := t.conn.Where("id_parent = ?", roadDamage.IdParent).Where("status = ?", "I").Order("revision DESC").First(&lastRoadDamageStatusI).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		} else {
			if err := t.conn.Table("road_damage").Where("id = ?", lastRoadDamageStatusI.Id).Updates(&models.RoadDamage{Status: "A"}).Error; err != nil {
				return err
			}
		}

	}

	roadDamageID := roadDamage.Id
	roadDamage.Status = "D"
	roadDamage.UpdatedDate = time.Now()
	roadDamage.UpdatedBy = userID
	if err := t.conn.Table("road_damage").Where("id = ?", roadDamageID).Updates(&roadDamage).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadDamageRepository) GetRoadInfo(roadID int) (models.RoadInfo, error) {
	var roadInfo models.RoadInfo
	err := t.conn.Where("road_id = ?", roadID).First(&roadInfo).Error
	if err != nil {
		return roadInfo, err
	}
	return roadInfo, nil
}

func (t *roadDamageRepository) GetRoadGeomByRoadIDByLaneNo(roadID int, laneNo int) (models.RoadGeom, error) {
	var roadGeom models.RoadGeom
	err := t.conn.Select("ST_ASTEXT(the_geom) as the_geom").Where("road_id = ?", roadID).Where("lane_no = ?", laneNo).Where("status = ?", "A").First(&roadGeom).Error
	if err != nil {
		return roadGeom, err
	}
	return roadGeom, nil
}

func (t *roadDamageRepository) UpdateRoadDamageImgPath(rdID int, imgPath string) error {
	var roadDamage models.RoadDamage
	if err := t.conn.Model(roadDamage).Where("id = ?", rdID).Update("img_filepath", imgPath).Error; err != nil {
		return err
	}
	return nil
}

func (t *roadDamageRepository) GetRoadSurfaceLane(roadID, laneNo int, km float64) (string, error) {
	reSurface := struct {
		ID         int    `json:"id"`
		RefSurface string `json:"ref_surface"`
	}{}
	// var roadSurface models.RoadSurfaceData // RoadSurfaceGroupType
	query := t.conn
	query = query.Table("road_surface")
	query = query.Select("CASE WHEN ref_surface_type.surface_group = 'Asphalt' THEN 'AC' ELSE 'CC' END as ref_surface")
	query = query.Where("road_surface.road_id = ?", roadID)
	query = query.Where("status = ?", "A")
	query = query.Where("lane_no = ?", laneNo)
	query = query.Where("km_start >= ? and km_start <= ?", km, km)
	query = query.Joins("INNER JOIN road_surface_lane on road_surface_lane.road_surface_id = road_surface.id")
	query = query.Joins("INNER JOIN ref_surface_type on ref_surface_type.id = road_surface_lane.ref_surface_id")
	if err := query.First(&reSurface).Error; err != nil {
		if err := query.First(&reSurface).Error; err != nil {
			if err.Error() != gorm.ErrRecordNotFound.Error() {
				return "", err
			}
		}
	}
	return reSurface.RefSurface, nil
}
func (t *roadDamageRepository) GetLastRoadInfoByID(roadID int) (models.RoadInfoGeomData, error) {
	var roadInfo models.RoadInfoGeomData
	query := t.conn
	query = query.Select("* , ST_ASTEXT(the_geom) as line_string")
	query = query.Table("road_info")
	query = query.Preload("RoadGeom", func(db *gorm.DB) *gorm.DB {
		return db.Select("*, ST_ASTEXT(the_geom) as line_string").Where("status = ?", "A")
	})
	if err := query.Where("road_id = ?", roadID).Where("status = ?", "A").Order("road_id asc").Find(&roadInfo).Error; err != nil {
		return roadInfo, err
	}

	return roadInfo, nil
}
