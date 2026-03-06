package repositories

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/volumeAccident/handlers"
	"gitlab.com/mims-api-service/src/volumeAccident/usecases"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type Repository struct {
	conn *gorm.DB
}

// init Repository Handler
func NewRepositoryHandler(conn *gorm.DB) *handlers.Handler {
	useCase := usecases.NewUseCase(&Repository{conn})
	handler := handlers.NewHandler(useCase)
	return handler
}

func (r *Repository) GetVolumeRevision(roadGrpID int, permissions []string) ([]models.VolumeAccidentRevision, error) {
	query := r.conn
	var volume []models.VolumeAccidentRevision
	query = query.Where("road_group_id = ?", roadGrpID)
	// query.roadDamageList
	//-------- Permission Check --------
	// if helpers.HasPermission([]string{"road_damage_manage_data"}, permissions) {
	var wheres []string
	wheres = append(wheres, "status = 'T'")
	wheres = append(wheres, "status = 'W'")
	wheres = append(wheres, "status = 'R'")
	wheres = append(wheres, "status = 'A'")
	queryWhereString := strings.Join(wheres, " or ")
	query = query.Where(queryWhereString)
	// } else {
	// 	query = query.Where("road_damage.status = ?", "A")
	// }

	query = query.Order("year DESC")
	query = query.Order("surveyed_date DESC")
	query = query.Order("updated_date DESC")
	query = query.Order("revision DESC")
	query = query.Order("status DESC")
	err := query.Find(&volume).Error
	if err != nil {
		return volume, err
	}
	return volume, nil
}

func (r *Repository) GetVolume(roadGrpID, ID int) (models.VolumeAccidentList, error) {
	var wheres []string
	wheres = append(wheres, "status = 'A'")
	wheres = append(wheres, "status = 'T'")
	wheres = append(wheres, "status = 'R'")
	wheres = append(wheres, "status = 'I'")
	wheres = append(wheres, "status = 'W'")
	queryWhereString := strings.Join(wheres, " or ")

	var volume models.VolumeAccidentList
	query := r.conn
	query = query.Preload("RoadGroup")
	query = query.Preload("UserDepartment", func(db *gorm.DB) *gorm.DB {
		return db.Select("email,username, firstname, lastname, tel, status, CONCAT('" + os.Getenv("STORAGE_IP") + "/" + "',profile_img_path) as profile_img_path")
	})
	// query = query.Preload("UserDepartment.Department")

	query = query.Where(queryWhereString)
	query = query.Where("id = ?", ID)
	query = query.Where("road_group_id = ?", roadGrpID)
	query = query.Order("revision DESC")
	if err := query.Preload("RoadGroup").First(&volume).Error; err != nil {
		fmt.Println(err)
		return volume, err
	}
	return volume, nil
}

func (r *Repository) GetVolumeByID(roadGrpID int, ID int) (models.VolumeAccident, error) {
	var wheres []string
	wheres = append(wheres, "status = 'A'")
	wheres = append(wheres, "status = 'T'")
	wheres = append(wheres, "status = 'R'")
	wheres = append(wheres, "status = 'I'")
	wheres = append(wheres, "status = 'W'")
	queryWhereString := strings.Join(wheres, " or ")

	var volume models.VolumeAccident
	query := r.conn
	query = query.Where(queryWhereString)
	if err := query.Where("id = ?", ID).Where("road_group_id = ?", roadGrpID).Order("revision DESC").First(&volume).Error; err != nil {
		fmt.Println(err)
		return volume, err
	}
	return volume, nil
}

func (r *Repository) CreateVolume(data models.VolumeAccident) (int, int, error) {
	query := r.conn
	if err := query.Table("volume_accident").Create(&data).Error; err != nil {
		return 0, 0, err
	}
	return data.ID, data.IDParent, nil
}

func (r *Repository) UpdateStatusD(idParent int) error {
	query := r.conn
	orWheres := []string{}
	orWheres = append(orWheres, "status = 'T'")
	orWheres = append(orWheres, "status = 'R'")
	queryWhereString := strings.Join(orWheres, " or ")
	if err := query.Table("volume_accident").Where("id_parent = ?", idParent).Where(queryWhereString).Updates(models.VolumeAccident{Status: "D", UpdatedDate: time.Now().UTC()}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateStatusI(IDParent int) error {
	if err := r.conn.Table("volume_accident").Where("id_parent = ?", IDParent).Where("status = ?", "A").Updates(models.VolumeAadt{Status: "I"}).Error; err != nil {
		return err
	}
	return nil
	// return nil
}

func (r *Repository) UpdateVolumeStatusT_To_DByGrpID(roadGrpID int) error {
	orWheres := []string{}
	orWheres = append(orWheres, "status = 'T'")
	orWheres = append(orWheres, "status = 'R'")
	queryWhereString := strings.Join(orWheres, " or ")
	if err := r.conn.Table("volume_accident").Where("road_group_id = ?", roadGrpID).Where(queryWhereString).Updates(models.VolumeAccident{Status: "D", UpdatedDate: time.Now().UTC()}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetMaxRevision(roadGrpID int, IDParent int) (models.VolumeAccident, error) {
	var volume models.VolumeAccident
	query := r.conn
	query = query.Where("road_group_id = ?", roadGrpID)
	query = query.Where("id_parent = ?", IDParent)
	query = query.Where("status != ?", "D")
	if err := query.Order("revision DESC").First(&volume).Error; err != nil {
		return volume, err
	}
	return volume, nil
}

func (t *Repository) DeleteVolume(roadGrpID, accidentID, userID int) error {
	var volume models.VolumeAccident
	var wheres []string
	wheres = append(wheres, "status = 'A'")
	wheres = append(wheres, "status = 'T'")
	wheres = append(wheres, "status = 'R'")
	wheres = append(wheres, "status = 'W'")
	queryWhereString := strings.Join(wheres, " or ")
	if err := t.conn.Where("road_group_id = ?", roadGrpID).Where("id = ?", accidentID).Where(queryWhereString).Order("revision DESC").First(&volume).Error; err != nil {
		return err
	}

	if volume.Status == "W" {
		return errors.New(constants.DATA_WAITING_APPROVAL)
	}

	if volume.Status == "A" {
		var lastVolumeStatusI models.VolumeAccident
		if err := t.conn.Where("road_group_id = ?", roadGrpID).Where("id_parent = ?", volume.IDParent).Where("status = ?", "I").Order("revision DESC").First(&lastVolumeStatusI).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		} else {
			if err := t.conn.Table("volume_accident").Where("id = ?", lastVolumeStatusI.ID).Updates(&models.VolumeAccident{Status: "A"}).Error; err != nil {
				return err
			}
		}

	}

	ID := volume.ID
	volume.Status = "D"
	volume.UpdatedDate = (time.Now().UTC())
	volume.UpdatedBy = userID
	if err := t.conn.Table("volume_accident").Where("id = ?", ID).Updates(&volume).Error; err != nil {
		return err
	}

	return nil
}

func (t *Repository) UpdateStatusIToA(roadGrpID, idParent int) error {
	var volume models.VolumeAccident
	if err := t.conn.Where("road_group_id = ? and id_parent = ? and status = ?", roadGrpID, idParent, "I").Order("revision desc").First(&volume).Error; err == nil {
		ID := volume.ID
		volume.Status = "A"
		if err := t.conn.Table("volume_accident").Where("id = ?", ID).Updates(&volume).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) UpdateVolumUpdateIdParent(ID int) error {
	if err := r.conn.Table("volume_accident").Where("id = ?", ID).Updates(models.VolumeAccident{IDParent: ID, UpdatedDate: time.Now().UTC()}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetTheGeomByRoadGrpID(roadGrpID int) ([]models.RoadInfo, error) {
	var roadInfo []models.RoadInfo
	if err := r.conn.Raw("SELECT road_info.id as id, road_id, year, road_info.ref_direction_id as ref_direction_id, name,km_start, km_end,st_astext(the_geom) as the_geom, revision, road_info.status as status, ramp_id, road_color_code from road JOIN road_info on road.id = road_info.road_id where road_group_id = ? ", roadGrpID).Find(&roadInfo).Error; err != nil {
		return roadInfo, err
	}
	return roadInfo, nil
}

func (r *Repository) GetStatus(statusCode string) (string, error) {
	var status models.RefDataStatus
	if err := r.conn.Where("status_code = ?", statusCode).First(&status).Error; err != nil {
		return "", err
	}
	return status.Name, nil
}
