package repositories

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/volumeAadt/handlers"
	"gitlab.com/mims-api-service/src/volumeAadt/usecases"

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

func (r *Repository) GetVolumeRevision(roadId int, permissions []string) ([]models.VolumeAadtRevision, error) {
	query := r.conn
	var volume []models.VolumeAadtRevision
	query = query.Where("road_id = ?", roadId)
	// query.roadDamageList
	//-------- Permission Check --------
	// if helpers.HasPermission([]string{"road_damage_manage_data"}, permissions) {
	query = query.Where("status = 'A'")
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
	// helpers.PrintlnJson(volumeAadt)
	return volume, nil
}

func (r *Repository) GetVolume(roadGrpID, aadtID int) (models.VolumeAadtList, error) {
	var wheres []string
	wheres = append(wheres, "status = 'A'")
	wheres = append(wheres, "status = 'I'")
	queryWhereString := strings.Join(wheres, " or ")

	var volume models.VolumeAadtList
	query := r.conn
	query = query.Where(queryWhereString)

	query = query.Preload("UserDepartment", func(db *gorm.DB) *gorm.DB {
		return db.Select("* , CONCAT('" + os.Getenv("STORAGE_IP") + "/" + "',profile_img_path) as profile_img_path")
	})
	// query = query.Preload("UserDepartment.Department")

	if err := query.Where("id = ?", aadtID).Where("road_id = ?", roadGrpID).Order("revision DESC").First(&volume).Error; err != nil {
		fmt.Println(err)
		return volume, err
	}
	return volume, nil
}

func (r *Repository) GetVolumeByID(roadGrpID int, aadtID int) (models.VolumeAadt, error) {
	var wheres []string
	wheres = append(wheres, "status = 'A'")
	wheres = append(wheres, "status = 'I'")
	queryWhereString := strings.Join(wheres, " or ")

	var volume models.VolumeAadt
	query := r.conn
	query = query.Where(queryWhereString)
	if err := query.Where("id = ?", aadtID).Where("road_id = ?", roadGrpID).Order("revision DESC").First(&volume).Error; err != nil {
		fmt.Println(err)
		return volume, err
	}
	return volume, nil
}

func (r *Repository) CreateVolume(data models.VolumeAadt) (int, int, error) {
	query := r.conn
	if err := query.Table("volume_aadt").Create(&data).Error; err != nil {
		return 0, 0, err
	}
	return data.ID, data.IDParent, nil
}

func (r *Repository) UpdateVolumeStatusT_To_DByGrpID(roadGrpID int) error {
	orWheres := []string{}
	orWheres = append(orWheres, "status = 'T'")
	orWheres = append(orWheres, "status = 'R'")
	queryWhereString := strings.Join(orWheres, " or ")
	if err := r.conn.Table("volume_aadt").Where("road_group_id = ?", roadGrpID).Where(queryWhereString).Updates(models.VolumeAadt{Status: "D", UpdatedDate: time.Now().UTC()}).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateStatusD(IDParent int) error {
	orWheres := []string{}
	orWheres = append(orWheres, "status = 'T'")
	orWheres = append(orWheres, "status = 'R'")
	queryWhereString := strings.Join(orWheres, " or ")
	if err := r.conn.Table("volume_aadt").Where("id_parent = ?", IDParent).Where(queryWhereString).Updates(models.VolumeAadt{Status: "D", UpdatedDate: time.Now().UTC()}).Error; err != nil {
		return err
	}
	return nil
	// return nil
}

func (r *Repository) UpdateStatusI(IDParent int) error {
	if err := r.conn.Table("volume_aadt").Where("id_parent = ?", IDParent).Where("status = ?", "A").Updates(models.VolumeAadt{Status: "I"}).Error; err != nil {
		return err
	}
	return nil
	// return nil
}
func (r *Repository) GetMaxRevision(roadId int, IDParent int) (models.VolumeAadt, error) {
	var volume models.VolumeAadt
	query := r.conn
	query = query.Where("road_id = ?", roadId)
	query = query.Where("id_parent = ?", IDParent)
	query = query.Where("status != ?", "D")
	if err := query.Order("revision DESC").First(&volume).Error; err != nil {
		return volume, err
	}
	return volume, nil
}

func (t *Repository) DeleteVolume(roadId, aadtID, userID int) error {
	var volume models.VolumeAadt
	if err := t.conn.Where("road_id = ?", roadId).Where("id = ?", aadtID).Where("status = 'A'").Order("revision DESC").First(&volume).Error; err != nil {
		return err
	}

	if volume.Status == "A" {
		var lastVolumeStatusI models.VolumeAadt
		if err := t.conn.Where("road_id = ?", roadId).Where("id_parent = ?", volume.IDParent).Where("status = ?", "I").Order("revision DESC").First(&lastVolumeStatusI).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		} else {
			if err := t.conn.Table("volume_aadt").Where("id = ?", lastVolumeStatusI.ID).Updates(&models.VolumeAadt{Status: "A"}).Error; err != nil {
				return err
			}
		}

	}

	ID := volume.ID
	volume.Status = "D"
	volume.UpdatedDate = (time.Now().UTC())
	volume.UpdatedBy = userID
	if err := t.conn.Table("volume_aadt").Where("id = ?", ID).Updates(&volume).Error; err != nil {
		return err
	}
	return nil
}

func (t *Repository) UpdateStatusIToA(roadId, idParent int) error {
	var volume models.VolumeAadt
	if err := t.conn.Where("road_id = ? and id_parent = ? and status = ?", roadId, idParent, "I").Order("revision desc").First(&volume).Error; err == nil {
		ID := volume.ID
		volume.Status = "A"
		if err := t.conn.Table("volume_aadt").Where("id = ?", ID).Updates(&volume).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) UpdateVolumUpdateIdParent(ID int) error {
	if err := r.conn.Table("volume_aadt").Where("id = ?", ID).Updates(models.VolumeAadt{IDParent: ID, UpdatedDate: time.Now().UTC()}).Error; err != nil {
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
